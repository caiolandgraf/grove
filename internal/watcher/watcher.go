package watcher

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
)

// ── ANSI helpers (self-contained so this package has no dep on package main) ──

const (
	ansiReset  = "\033[0m"
	ansiBold   = "\033[1m"
	ansiDim    = "\033[2m"
	ansiRed    = "\033[38;2;220;60;60m"
	ansiGreen  = "\033[38;2;40;210;90m"
	ansiYellow = "\033[38;2;230;200;40m"
	ansiCyan   = "\033[38;2;80;220;220m"
	ansiGray   = "\033[38;2;130;130;145m"

	ansiBgGreen = "\033[48;2;40;180;80m\033[38;2;255;255;255m"
	ansiBgRed   = "\033[48;2;195;55;55m\033[38;2;255;255;255m"
	ansiBgBlue  = "\033[48;2;60;120;220m\033[38;2;255;255;255m"
	ansiBgGrove = "\033[48;2;40;180;80m\033[38;2;255;255;255m"
)

func badge(
	bg, label string,
) string {
	return bg + " " + label + " " + ansiReset
}
func logDev(msg string) { fmt.Println("  " + msg) }

// ──────────────────────────────────────────────────────────────────────────────

// Watcher ties together the fsnotify subscription, the debounce timer, and the
// child-process manager.  Create one with New and call Start to begin the loop.
type Watcher struct {
	cfg  Config
	proc *Process

	fsw *fsnotify.Watcher

	// mu guards the debounce timer so that concurrent fsnotify callbacks never
	// schedule two simultaneous rebuilds.
	mu       sync.Mutex
	debounce *time.Timer

	// rebuildCh decouples the fsnotify goroutine from the rebuild goroutine so
	// that a slow build never blocks the watcher event loop.
	rebuildCh chan struct{}
}

// New returns a ready-to-use Watcher.  Call Start to begin watching.
func New(cfg Config) *Watcher {
	return &Watcher{
		cfg:       cfg,
		proc:      &Process{},
		rebuildCh: make(chan struct{}, 1),
	}
}

// Start performs an initial build+run, then enters the fsnotify event loop.
// It blocks until the user sends SIGINT / SIGTERM.
func (w *Watcher) Start() error {
	// ── Ensure tmp directory exists ──────────────────────────────────────────
	if err := os.MkdirAll(w.cfg.TmpDir, 0o755); err != nil {
		return fmt.Errorf("cannot create tmp_dir %q: %w", w.cfg.TmpDir, err)
	}

	// ── Set up the fsnotify watcher ──────────────────────────────────────────
	fsw, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("cannot create fsnotify watcher: %w", err)
	}
	w.fsw = fsw
	defer fsw.Close() //nolint:errcheck

	// Recursively watch every configured directory.
	for _, dir := range w.cfg.WatchDirs {
		if err := w.addRecursive(dir); err != nil {
			logDev(
				ansiYellow + "⚠  Cannot watch " + dir + ": " + err.Error() + ansiReset,
			)
		}
	}

	// ── Initial build + launch ───────────────────────────────────────────────
	w.printHeader()
	w.runRebuild()

	// ── Signal handling ───────────────────────────────────────────────────────
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	// ── Rebuild worker ────────────────────────────────────────────────────────
	// A dedicated goroutine drains rebuildCh so the fsnotify loop is never
	// blocked by a long compilation.
	go func() {
		for range w.rebuildCh {
			w.runRebuild()
		}
	}()

	// ── Event loop ───────────────────────────────────────────────────────────
	for {
		select {

		case event, ok := <-fsw.Events:
			if !ok {
				return nil
			}

			// Auto-watch newly created subdirectories so files added after
			// startup are still detected.
			if event.Has(fsnotify.Create) {
				if info, err := os.Stat(
					event.Name,
				); err == nil &&
					info.IsDir() {
					_ = w.addRecursive(event.Name)
				}
			}

			if w.shouldHandle(event) {
				w.scheduleRebuild()
			}

		case err, ok := <-fsw.Errors:
			if !ok {
				return nil
			}
			logDev(ansiYellow + "⚠  Watcher error: " + err.Error() + ansiReset)

		case <-sigCh:
			fmt.Println()
			logDev(ansiGray + "Stopping application…" + ansiReset)
			w.proc.Stop()
			fmt.Println()
			logDev(
				badge(
					ansiBgGrove,
					"GROVE DEV",
				) + "  " + ansiGray + "stopped." + ansiReset,
			)
			fmt.Println()
			return nil
		}
	}
}

// ── Filtering ─────────────────────────────────────────────────────────────────

// shouldHandle returns true when event should trigger a rebuild:
//   - Op must be Write or Create (Rename/Remove/Chmod are ignored).
//   - The path must not be inside an excluded directory.
//   - The filename must not end in _spec.go (test specs never trigger a rebuild).
//   - The file extension must be in the configured allow-list.
func (w *Watcher) shouldHandle(event fsnotify.Event) bool {
	if !event.Has(fsnotify.Write) && !event.Has(fsnotify.Create) {
		return false
	}

	if w.isExcluded(event.Name) {
		return false
	}

	// Spec files must never cause a rebuild regardless of their location.
	if strings.HasSuffix(filepath.Base(event.Name), "_spec.go") {
		return false
	}

	ext := filepath.Ext(event.Name)
	for _, allowed := range w.cfg.Extensions {
		if ext == allowed {
			return true
		}
	}

	return false
}

// isExcluded returns true when any path component of p matches an entry in
// cfg.Exclude exactly (e.g. ".grove" excludes ".grove/tmp/app").
func (w *Watcher) isExcluded(p string) bool {
	// Normalise to forward slashes for consistent splitting on all platforms.
	parts := strings.Split(filepath.ToSlash(p), "/")
	for _, part := range parts {
		for _, excl := range w.cfg.Exclude {
			if part == excl {
				return true
			}
		}
	}
	return false
}

// ── Debounce ──────────────────────────────────────────────────────────────────

// scheduleRebuild arms (or resets) the debounce timer.  When the timer fires
// it sends a single token on rebuildCh, which the rebuild worker drains.
func (w *Watcher) scheduleRebuild() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.debounce != nil {
		w.debounce.Stop()
	}

	delay := time.Duration(w.cfg.DebounceMs) * time.Millisecond

	w.debounce = time.AfterFunc(delay, func() {
		// Non-blocking send: if a rebuild is already queued the worker will
		// pick it up; we don't need to queue another.
		select {
		case w.rebuildCh <- struct{}{}:
		default:
		}
	})
}

// ── Build + restart ───────────────────────────────────────────────────────────

// runRebuild compiles the project and, on success, restarts the binary.
// Build errors are printed but do not stop the watcher.
func (w *Watcher) runRebuild() {
	fmt.Println()
	logDev(
		badge(ansiBgBlue, "RE-BUILDING"),
	)
	fmt.Println()

	start := time.Now()

	if err := w.build(); err != nil {
		fmt.Println()
		logDev(badge(ansiBgRed, "BUILD FAILED"))
		fmt.Println()
		return
	}

	elapsed := time.Since(start)

	if err := w.proc.Restart(w.cfg.Bin); err != nil {
		fmt.Println()
		logDev(
			badge(
				ansiBgRed,
				"ERRO",
			) + "  " + ansiRed + ansiBold + "❌ Failed to start binary: " + ansiReset + ansiRed + err.Error() + ansiReset,
		)
		fmt.Println()
		return
	}

	fmt.Println()
	logDev(
		badge(ansiBgGreen, "APP RESTARTED") +
			"  " + ansiGray + "(" + fmtElapsed(elapsed) + ")" + ansiReset,
	)
	fmt.Println()
}

// build runs cfg.BuildCmd in cfg.Root, piping stderr to the terminal.
func (w *Watcher) build() error {
	parts := strings.Fields(w.cfg.BuildCmd)
	if len(parts) == 0 {
		return fmt.Errorf("build_cmd is empty")
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Dir = w.cfg.Root
	// Pipe compiler output through the build writer which colourises each line.
	bw := newBuildOutputWriter(os.Stderr)
	cmd.Stdout = bw
	cmd.Stderr = bw

	return cmd.Run()
}

// ── Recursive watch ───────────────────────────────────────────────────────────

// addRecursive walks dir and registers every non-excluded subdirectory with
// the fsnotify watcher.  Errors for individual directories are silently
// skipped so that a missing watch_dir entry does not abort startup.
func (w *Watcher) addRecursive(dir string) error {
	return filepath.WalkDir(
		dir,
		func(path string, d os.DirEntry, err error) error {
			if err != nil {
				// Unreadable entry — skip rather than abort the whole walk.
				return nil
			}
			if !d.IsDir() {
				return nil
			}
			if w.isExcluded(path) {
				return filepath.SkipDir
			}

			if watchErr := w.fsw.Add(path); watchErr != nil {
				logDev(
					ansiYellow + "⚠  Cannot watch " + path + ": " + watchErr.Error() + ansiReset,
				)
			}
			return nil
		},
	)
}

// ── UI helpers ────────────────────────────────────────────────────────────────

func (w *Watcher) printHeader() {
	sep := "  " + ansiDim + strings.Repeat("─", 54) + ansiReset

	exts := strings.Join(w.cfg.Extensions, " ")
	dirs := strings.Join(w.cfg.WatchDirs, ", ")

	fmt.Println()
	fmt.Println(sep)
	fmt.Println()
	logDev(
		badge(
			ansiBgGrove,
			"GROVE DEV",
		) + "  " + ansiGray + "watching for changes — Ctrl+C to stop" + ansiReset,
	)
	fmt.Println()
	logDev(
		ansiGray + "  extensions  " + ansiReset + ansiBold + exts + ansiReset,
	)
	logDev(
		ansiGray + "  watch dirs  " + ansiReset + ansiBold + dirs + ansiReset,
	)
	logDev(
		ansiGray + "  debounce    " + ansiReset + ansiBold + fmt.Sprintf(
			"%d ms",
			w.cfg.DebounceMs,
		) + ansiReset,
	)
	logDev(
		ansiGray + "  binary      " + ansiReset + ansiBold + w.cfg.Bin + ansiReset,
	)
	fmt.Println()
	fmt.Println(sep)
	fmt.Println()
}

func fmtElapsed(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	return fmt.Sprintf("%.1fs", d.Seconds())
}

// ── buildOutputWriter ─────────────────────────────────────────────────────────

// buildOutputWriter processes go compiler output line by line and applies
// context-aware colouring:
//
//   - Package header lines ("# module/pkg") → gray + dim, no error marker.
//   - All other non-empty lines            → fully red, prefixed with ×.
//
// This matches the visual language of gest and makes it easy to scan which
// package failed and which specific symbols are undefined.
type buildOutputWriter struct {
	w   *os.File
	buf []byte
}

func newBuildOutputWriter(w *os.File) *buildOutputWriter {
	return &buildOutputWriter{w: w}
}

func (bw *buildOutputWriter) Write(p []byte) (n int, err error) {
	bw.buf = append(bw.buf, p...)

	for {
		nl := bytes.IndexByte(bw.buf, '\n')
		if nl < 0 {
			break
		}
		bw.writeLine(string(bw.buf[:nl]))
		bw.buf = bw.buf[nl+1:]
	}

	return len(p), nil
}

// writeLine emits a single compiler output line with appropriate styling.
func (bw *buildOutputWriter) writeLine(line string) {
	if strings.TrimSpace(line) == "" {
		fmt.Fprintln(bw.w)
		return
	}

	if strings.HasPrefix(line, "# ") {
		// Package header — subdued so it doesn't compete with the errors.
		fmt.Fprintf(bw.w, "  %s%s%s\n", ansiGray+ansiDim, line, ansiReset)
		return
	}

	// Error / warning line — fully red with × marker.
	fmt.Fprintf(bw.w, "  %s× %s%s\n", ansiRed, line, ansiReset)
}
