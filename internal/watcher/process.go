package watcher

import (
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// Process manages the lifecycle of the running application binary.
// It is safe for concurrent use — a sync.Mutex serialises all Restart calls
// so that rapid file-change events never spawn duplicate processes.
type Process struct {
	mu     sync.Mutex
	cmd    *exec.Cmd
	waitCh chan struct{} // closed by the reaper goroutine when the process exits
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
func (p *Process) Restart(bin string) error {
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

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Start(); err != nil {
		return err
	}

	// Allocate a fresh channel before handing cmd off to the reaper so there
	// is never a window where p.waitCh is nil while p.cmd is non-nil.
	waitCh := make(chan struct{})
	p.cmd = cmd
	p.waitCh = waitCh

	// ── 3. Reap the process in the background ────────────────────────────────
	// cmd.Wait() is called exactly once per process — only here — so there is
	// no risk of a "wait: process already finished" error.
	go func() {
		_ = cmd.Wait()
		close(waitCh)
	}()

	return nil
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
