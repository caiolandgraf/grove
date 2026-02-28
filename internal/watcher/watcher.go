package watcher

import (
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
			logDev(ansiBold + "🌿 Grove dev stopped." + ansiReset)
			fmt.Println()
			return nil
		}
	}
}

// ── Filtering ─────────────────────────────────────────────────────────────────

// shouldHandle returns true when event should trigger a rebuild:
//   - Op must be Write or Create (Rename/Remove/Chmod are ignored).
//   - The path must not be inside an excluded directory.
//   - The file extension must be in the configured allow-list.
func (w *Watcher) shouldHandle(event fsnotify.Event) bool {
	if !event.Has(fsnotify.Write) && !event.Has(fsnotify.Create) {
		return false
	}

	if w.isExcluded(event.Name) {
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
	logDev(ansiBold + ansiCyan + "🔨 Recompilando…" + ansiReset)

	start := time.Now()

	if err := w.build(); err != nil {
		fmt.Println()
		logDev(
			badge(
				ansiBgRed,
				"ERRO",
			) + "  " + ansiRed + "❌ Build error:" + ansiReset,
		)
		// The build command already wrote the compiler output to Stderr, so
		// we only print a short summary here.
		logDev(ansiDim + "    " + err.Error() + ansiReset)
		fmt.Println()
		return
	}

	elapsed := time.Since(start)

	if err := w.proc.Restart(w.cfg.Bin); err != nil {
		logDev(
			badge(
				ansiBgRed,
				"ERRO",
			) + "  " + ansiRed + "❌ Failed to start binary: " + err.Error() + ansiReset,
		)
		return
	}

	fmt.Println()
	logDev(
		badge(ansiBgGreen, "OK") + "  " +
			ansiGreen + "✅ App reiniciado" + ansiReset +
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
	// Pipe compiler output indented for visual alignment.
	cmd.Stdout = newPrefixWriter(os.Stdout, "    ")
	cmd.Stderr = newPrefixWriter(os.Stderr, "    ")

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
		ansiBold + "🌿 Grove dev" + ansiReset + "  " + ansiGray + "watching for changes — Ctrl+C to stop" + ansiReset,
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

// ── prefixWriter ─────────────────────────────────────────────────────────────

// newPrefixWriter wraps w and prepends prefix to every line of output.
// This keeps build errors visually aligned with grove's own log lines.
func newPrefixWriter(w *os.File, prefix string) *prefixWriter {
	return &prefixWriter{w: w, pfx: []byte(prefix), atSOL: true}
}

type prefixWriter struct {
	w     *os.File
	pfx   []byte
	atSOL bool
}

func (pw *prefixWriter) Write(p []byte) (n int, err error) {
	for len(p) > 0 {
		if pw.atSOL {
			if _, err = pw.w.Write(pw.pfx); err != nil {
				return
			}
			pw.atSOL = false
		}
		idx := strings.IndexByte(string(p), '\n')
		if idx < 0 {
			var nn int
			nn, err = pw.w.Write(p)
			n += nn
			return
		}
		var nn int
		nn, err = pw.w.Write(p[:idx+1])
		n += nn
		if err != nil {
			return
		}
		p = p[idx+1:]
		pw.atSOL = true
	}
	return
}
