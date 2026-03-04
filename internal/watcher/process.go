package watcher

import (
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// appOut is the shared writer used by all child processes to format their
// stdout/stderr through appOutputWriter before printing to the terminal.
var appOut = newAppOutputWriter(os.Stdout)

// Process manages the lifecycle of the running application binary.
// It is safe for concurrent use — a sync.Mutex serialises all Restart calls
// so that rapid file-change events never spawn duplicate processes.
type Process struct {
	mu       sync.Mutex
	cmd      *exec.Cmd
	waitCh   chan struct{}   // closed by the reaper goroutine when the process exits
	lastDone <-chan struct{} // DoneCh of the most recently launched process
}

// RestartResult is returned by Restart and lets the caller observe whether the
// newly launched process is still alive after a short stabilisation window.
type RestartResult struct {
	// ReadyCh is closed once the stabilisation window has elapsed AND the
	// process is still running. If the process crashed before the window
	// expired, ReadyCh is never closed — CrashCh is closed instead.
	ReadyCh <-chan struct{}

	// CrashCh is closed if the process exits within the stabilisation window
	// (i.e. a fast crash / startup panic). The caller should treat this as a
	// failed start rather than a clean launch.
	CrashCh <-chan struct{}

	// DoneCh mirrors the internal waitCh: it is closed when the process has
	// fully exited and all pipe output has been flushed. The caller can wait
	// on this channel to ensure all output (including panic dumps) has been
	// printed before starting a new rebuild.
	DoneCh <-chan struct{}
}

// Restart stops the currently running process (if any) and starts a new one
// from the binary at bin.
//
// Shutdown sequence:
//  1. Send os.Interrupt (SIGINT on Unix) so the app can clean up.
//  2. Wait up to 5 s for the process to exit on its own.
//  3. If it is still alive after the grace period, send os.Kill (SIGKILL).
//
// A new process is then launched unconditionally regardless of how the old
// one exited, so a build that produces a working binary always gets a chance
// to run.
//
// The returned RestartResult lets the caller distinguish between a healthy
// start and an immediate crash (e.g. a startup panic).
func (p *Process) Restart(bin string) (RestartResult, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// ── 1. Gracefully stop the old process ───────────────────────────────────
	if p.cmd != nil && p.cmd.Process != nil {
		// Best-effort interrupt — ignore the error (the process may have
		// already exited on its own between the check and the signal).
		_ = p.cmd.Process.Signal(os.Interrupt)

		select {
		case <-p.waitCh:
			// Exited cleanly within the grace period.
		case <-time.After(5 * time.Second):
			// Timed out — escalate to a hard kill.
			_ = p.cmd.Process.Kill()
			// Still wait so the OS can reclaim the PID before we start the
			// next process; the reaper goroutine will close waitCh shortly.
			<-p.waitCh
		}
	}

	// ── 2. Launch the new binary ─────────────────────────────────────────────
	// Split the binary path from any embedded arguments (unusual but safe to
	// support, e.g. ".grove/tmp/app --port 8080").
	parts := strings.Fields(bin)
	var cmd *exec.Cmd
	if len(parts) == 1 {
		cmd = exec.Command(parts[0])
	} else {
		cmd = exec.Command(parts[0], parts[1:]...)
	}

	// Use explicit pipes for stdout and stderr so we can drain them fully
	// before calling Flush(). This guarantees that panic output written by
	// the Go runtime to stderr is never lost even when the process exits
	// immediately after writing (before the OS delivers the last bytes via
	// the default cmd.Stdout / cmd.Stderr assignment).
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return RestartResult{}, err
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return RestartResult{}, err
	}
	cmd.Stdin = os.Stdin

	if err := cmd.Start(); err != nil {
		return RestartResult{}, err
	}

	// Allocate a fresh channel before handing cmd off to the reaper so there
	// is never a window where p.waitCh is nil while p.cmd is non-nil.
	waitCh := make(chan struct{})
	p.cmd = cmd
	p.waitCh = waitCh

	// ── 3. Drain pipes + reap ─────────────────────────────────────────────────
	// We drain stdout and stderr in dedicated goroutines and wait for both to
	// finish before calling Flush(). This ensures every byte written by the
	// child — including a panic dumped by the Go runtime right before exit —
	// has been processed by appOut before we declare the process done.
	go func() {
		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			_, _ = io.Copy(appOut, stdoutPipe)
		}()
		go func() {
			defer wg.Done()
			_, _ = io.Copy(appOut, stderrPipe)
		}()

		// Wait for both pipes to be fully drained (EOF), then wait for the
		// process to exit, then flush any partial panic buffer.
		wg.Wait()
		_ = cmd.Wait()
		appOut.Flush()
		close(waitCh)
	}()

	// ── 4. Build result channels ──────────────────────────────────────────────
	// stabilisationWindow is how long we wait before declaring the process
	// healthy. Startup panics in Go programs typically print and exit within
	// a few milliseconds; 500 ms is a comfortable margin that still feels
	// instant to the developer.
	const stabilisationWindow = 500 * time.Millisecond

	readyCh := make(chan struct{})
	crashCh := make(chan struct{})

	go func() {
		select {
		case <-waitCh:
			// Process exited within the stabilisation window — it's a crash.
			close(crashCh)
		case <-time.After(stabilisationWindow):
			// Process is still running after the window — declare it ready.
			close(readyCh)
		}
	}()

	result := RestartResult{
		ReadyCh: readyCh,
		CrashCh: crashCh,
		DoneCh:  waitCh,
	}
	p.lastDone = waitCh
	return result, nil
}

// WaitDone waits for the most recently launched process's output to be fully
// flushed, but ONLY if the process has already exited on its own.
//
// If the process is still running, this is a no-op — the caller (runRebuild)
// will subsequently call Restart, which sends SIGINT and drains waitCh itself
// before launching the next binary. Blocking here in that case would deadlock
// because Restart has not yet been called to stop the process.
func (p *Process) WaitDone() {
	p.mu.Lock()
	done := p.lastDone
	p.mu.Unlock()

	if done == nil {
		return
	}

	// Non-blocking check: only wait if the process has already exited.
	select {
	case <-done:
		// Already done — output has been flushed.
	default:
		// Still running — Restart will handle teardown.
	}
}

// Stop sends an interrupt signal to the running process and waits for it to
// exit. It is intended for clean shutdown (e.g. when the user hits Ctrl-C).
// If the process has already exited Stop is a no-op.
func (p *Process) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.cmd == nil || p.cmd.Process == nil {
		return
	}

	_ = p.cmd.Process.Signal(os.Interrupt)

	select {
	case <-p.waitCh:
	case <-time.After(5 * time.Second):
		_ = p.cmd.Process.Kill()
		<-p.waitCh
	}

	p.cmd = nil
}
