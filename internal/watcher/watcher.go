package watcher

import (
	"bytes"
	"encoding/json"
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

	ansiBgGreen  = "\033[48;2;40;180;80m\033[38;2;255;255;255m"
	ansiBgRed    = "\033[48;2;195;55;55m\033[38;2;255;255;255m"
	ansiBgBlue   = "\033[48;2;60;120;220m\033[38;2;255;255;255m"
	ansiBgYellow = "\033[48;2;185;140;20m\033[38;2;255;255;255m"
	ansiBgGrove  = "\033[48;2;40;180;80m\033[38;2;255;255;255m"
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
	// Wait for the previous process to fully drain its output (including any
	// panic dump) before resetting state and printing the RE-BUILDING banner.
	// This prevents the banner from interleaving with the crash output of a
	// process that panicked right after the last restart.
	if w.proc != nil {
		w.proc.WaitDone()
	}

	// Reset per-session state (hints, panic buffer) so every rebuild starts
	// clean and hints are shown again if the error persists.
	appOut.resetSession()

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

	result, err := w.proc.Restart(w.cfg.Bin)
	if err != nil {
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

	// Wait for either the stabilisation window (process is healthy) or an
	// immediate crash (startup panic). In the crash case, DoneCh will be
	// closed shortly after — we wait on it so the full panic output is
	// printed before we return and the caller can schedule the next rebuild.
	select {
	case <-result.ReadyCh:
		// Process survived the stabilisation window — it looks healthy.
		fmt.Println()
		logDev(
			badge(ansiBgGreen, "APP RESTARTED") +
				"  " + ansiGray + "(" + fmtElapsed(elapsed) + ")" + ansiReset,
		)
		fmt.Println()

	case <-result.CrashCh:
		// Process exited immediately — wait for all pipe output (panic dump,
		// error messages) to be fully flushed before returning.
		<-result.DoneCh
		fmt.Println()
		logDev(
			badge(ansiBgRed, "APP CRASHED") +
				"  " + ansiGray + "process exited immediately after start" + ansiReset,
		)
		fmt.Println()
	}
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

// ── appOutputWriter ───────────────────────────────────────────────────────────

// appOutputWriter processes the running application's stdout/stderr line by
// line and applies context-aware formatting:
//
//   - JSON log lines (slog / zap / zerolog structured output) are parsed and
//     rendered as a human-readable coloured line.
//   - "panic:" lines and the goroutine stack trace that follows are rendered
//     in a red block with a clear PANIC badge.
//   - All other lines are indented and passed through as-is.
type appOutputWriter struct {
	w        *os.File
	buf      []byte
	inPanic  bool
	panicBuf []string
	hintSeen map[string]bool
}

func newAppOutputWriter(w *os.File) *appOutputWriter {
	return &appOutputWriter{w: w, hintSeen: map[string]bool{}}
}

// resetSession clears per-run state so hints are shown again on every rebuild.
func (aw *appOutputWriter) resetSession() {
	aw.inPanic = false
	aw.panicBuf = nil
	aw.hintSeen = map[string]bool{}
}

// Flush flushes any buffered output that has not yet been written.
// It must be called when the child process exits so that a panic dump
// that ends at EOF (without a trailing non-stack line) is not silently
// discarded.
func (aw *appOutputWriter) Flush() {
	// Flush any partial line left in the byte buffer (no trailing newline).
	if len(aw.buf) > 0 {
		aw.writeLine(string(aw.buf))
		aw.buf = nil
	}
	// If we were accumulating a panic dump, print it now.
	if aw.inPanic {
		aw.flushPanic()
	}
}

func (aw *appOutputWriter) Write(p []byte) (n int, err error) {
	aw.buf = append(aw.buf, p...)

	for {
		nl := bytes.IndexByte(aw.buf, '\n')
		if nl < 0 {
			break
		}
		aw.writeLine(string(aw.buf[:nl]))
		aw.buf = aw.buf[nl+1:]
	}

	return len(p), nil
}

func (aw *appOutputWriter) writeLine(line string) {
	trimmed := strings.TrimSpace(line)

	// ── blank line ────────────────────────────────────────────────────────────
	if trimmed == "" {
		if aw.inPanic {
			aw.panicBuf = append(aw.panicBuf, "")
		} else {
			fmt.Fprintln(aw.w)
		}
		return
	}

	// ── panic accumulation ────────────────────────────────────────────────────
	// Once we see "panic:" we collect lines until the goroutine dump ends
	// (i.e. we hit a line that is not a stack frame, continuation or blank).
	if strings.HasPrefix(trimmed, "panic:") {
		aw.inPanic = true
		aw.panicBuf = []string{trimmed}
		return
	}

	if aw.inPanic {
		// Goroutine header, stack frames, file:line references all belong to
		// the panic dump.  We keep collecting until the line doesn't look like
		// part of a stack trace anymore.
		isStackLine := strings.HasPrefix(trimmed, "goroutine ") ||
			strings.HasPrefix(trimmed, "main.") ||
			strings.HasPrefix(trimmed, "runtime.") ||
			strings.HasPrefix(line, "\t") ||
			strings.Contains(trimmed, ".go:")

		if isStackLine {
			aw.panicBuf = append(aw.panicBuf, line)
			return
		}

		// End of panic dump — flush everything at once.
		aw.flushPanic()
	}

	// ── JSON structured log line ──────────────────────────────────────────────
	if len(trimmed) > 0 && trimmed[0] == '{' {
		if rendered, allText, ok := renderJSONLog(trimmed); ok {
			fmt.Fprintln(aw.w, rendered)
			aw.detectHints(allText)
			return
		}
	}

	// ── slog text format (key=value) ──────────────────────────────────────────
	// Matches the default output of log/slog's TextHandler:
	//   time=2006-01-02T15:04:05.000Z level=INFO msg="hello" key=val
	// Also matches logfmt-style lines that lack a time= prefix as long as
	// they contain at least level= and msg= (or message=).
	if rendered, allText, ok := renderSlogText(trimmed); ok {
		fmt.Fprintln(aw.w, rendered)
		aw.detectHints(allText)
		return
	}

	// ── plain line ────────────────────────────────────────────────────────────
	fmt.Fprintf(aw.w, "  %s\n", line)
}

// flushPanic prints the accumulated panic dump as a styled red block.
// Known panic messages get an extra actionable hint printed below the block.
func (aw *appOutputWriter) flushPanic() {
	aw.inPanic = false
	if len(aw.panicBuf) == 0 {
		return
	}

	// Extract the raw panic message for hint detection.
	panicMsg := ""
	for _, l := range aw.panicBuf {
		trimmed := strings.TrimSpace(l)
		if strings.HasPrefix(trimmed, "panic:") {
			panicMsg = strings.ToLower(
				strings.TrimSpace(strings.TrimPrefix(trimmed, "panic:")),
			)
			break
		}
	}

	fmt.Fprintln(aw.w)
	fmt.Fprintf(aw.w, "  %s\n", badge(ansiBgRed, "PANIC"))
	fmt.Fprintln(aw.w)

	for _, l := range aw.panicBuf {
		trimmed := strings.TrimSpace(l)
		if trimmed == "" {
			fmt.Fprintln(aw.w)
			continue
		}
		// First line is the panic message itself — bold red.
		if strings.HasPrefix(trimmed, "panic:") {
			msg := strings.TrimSpace(strings.TrimPrefix(trimmed, "panic:"))
			fmt.Fprintf(aw.w, "  %s%s%s\n", ansiRed+ansiBold, msg, ansiReset)
			continue
		}
		// Goroutine header.
		if strings.HasPrefix(trimmed, "goroutine ") {
			fmt.Fprintf(aw.w, "\n  %s%s%s\n", ansiGray, trimmed, ansiReset)
			continue
		}
		// File/line references (indented with tab in original).
		if strings.HasPrefix(l, "\t") {
			fmt.Fprintf(
				aw.w,
				"    %s%s%s\n",
				ansiDim+ansiGray,
				trimmed,
				ansiReset,
			)
			continue
		}
		// Function names.
		fmt.Fprintf(aw.w, "  %s%s%s\n", ansiGray, trimmed, ansiReset)
	}

	fmt.Fprintln(aw.w)

	// ── Actionable hints for known panic messages ─────────────────────────────
	if strings.Contains(panicMsg, ".env") &&
		strings.Contains(panicMsg, "not found") {
		printHint(aw.w,
			"Environment file not found.",
			[]string{
				"cp .env.example .env",
				"# then edit .env with your database credentials",
			},
		)
	}

	aw.panicBuf = nil
}

// renderSlogText attempts to parse a log/slog TextHandler (or logfmt) line of
// the form:
//
//	time=2006-01-02T15:04:05.000Z level=INFO msg="hello world" key=val key2="v 2"
//
// It returns (rendered, allText, true) on success, ("", "", false) if the line
// does not look like a slog text entry (must have at least level= and msg= /
// message= keys).
func renderSlogText(line string) (string, string, bool) {
	kv := parseSlogKV(line)
	if len(kv) == 0 {
		return "", "", false
	}

	level := firstOf(kv, "level", "lvl", "severity")
	msg := firstOf(kv, "msg", "message", "Message")
	ts := firstOf(kv, "time", "ts", "timestamp", "Time")

	// Require at least level + msg to treat this as a structured log line.
	if level == "" || msg == "" {
		return "", "", false
	}

	// ── Level badge ───────────────────────────────────────────────────────────
	var levelBadge string
	switch strings.ToUpper(level) {
	case "DEBUG":
		levelBadge = ansiGray + ansiBold + "DBG" + ansiReset
	case "INFO":
		levelBadge = ansiGreen + ansiBold + "INF" + ansiReset
	case "WARN", "WARNING":
		levelBadge = ansiYellow + ansiBold + "WRN" + ansiReset
	case "ERROR", "ERR":
		levelBadge = ansiRed + ansiBold + "ERR" + ansiReset
	case "FATAL", "PANIC":
		levelBadge = ansiBgRed + " " + level + " " + ansiReset
	default:
		levelBadge = ansiGray + ansiBold + "LOG" + ansiReset
	}

	// ── Timestamp — keep only HH:MM:SS ───────────────────────────────────────
	timeStr := ""
	if ts != "" {
		if idx := strings.IndexByte(ts, 'T'); idx >= 0 && len(ts) > idx+9 {
			timeStr = ansiDim + ansiGray + ts[idx+1:idx+9] + ansiReset + "  "
		}
	}

	// ── Message colour ────────────────────────────────────────────────────────
	var msgColour string
	switch strings.ToUpper(level) {
	case "ERROR", "ERR", "FATAL", "PANIC":
		msgColour = ansiRed
	case "WARN", "WARNING":
		msgColour = ansiYellow
	}
	msgPart := msg
	if msgColour != "" {
		msgPart = msgColour + msg + ansiReset
	}

	// ── Extra fields (skip standard ones already rendered) ────────────────────
	skip := map[string]bool{
		"level": true, "lvl": true, "severity": true,
		"msg": true, "message": true, "Message": true,
		"time": true, "ts": true, "timestamp": true, "Time": true,
	}

	// Build allText for hint detection while collecting extras.
	allText := strings.ToLower(msg)
	var extras []string
	for _, k := range kvKeys(kv) {
		v := kv[k]
		allText += " " + strings.ToLower(v)
		if skip[k] {
			continue
		}
		val := v
		if len(val) > 120 {
			val = val[:117] + "…"
		}
		extras = append(extras, ansiGray+k+"="+ansiReset+ansiDim+val+ansiReset)
	}

	extra := ""
	if len(extras) > 0 {
		extra = "  " + strings.Join(extras, "  ")
	}

	return fmt.Sprintf(
		"  %s%s%s  %s%s",
		timeStr, levelBadge, ansiReset, msgPart, extra,
	), allText, true
}

// parseSlogKV parses a slog TextHandler / logfmt line into an ordered map.
// Values may be bare tokens or double-quoted strings (with \" escaping).
// The insertion order of keys is preserved via kvKeys.
func parseSlogKV(line string) map[string]string {
	out := make(map[string]string)
	rest := line
	for len(rest) > 0 {
		rest = strings.TrimLeft(rest, " \t")
		if rest == "" {
			break
		}

		// Find the '=' that separates key from value.
		eqIdx := strings.IndexByte(rest, '=')
		if eqIdx <= 0 {
			break
		}
		key := rest[:eqIdx]
		rest = rest[eqIdx+1:]

		// Keys must be plain identifiers (no spaces).
		if strings.ContainsAny(key, " \t\"") {
			break
		}

		var val string
		if len(rest) > 0 && rest[0] == '"' {
			// Quoted value — scan to the closing unescaped '"'.
			i := 1
			var buf strings.Builder
			for i < len(rest) {
				if rest[i] == '\\' && i+1 < len(rest) {
					buf.WriteByte(rest[i+1])
					i += 2
					continue
				}
				if rest[i] == '"' {
					i++
					break
				}
				buf.WriteByte(rest[i])
				i++
			}
			val = buf.String()
			rest = rest[i:]
		} else {
			// Bare value — ends at the next space.
			end := strings.IndexAny(rest, " \t")
			if end < 0 {
				val = rest
				rest = ""
			} else {
				val = rest[:end]
				rest = rest[end:]
			}
		}

		if _, exists := out[key]; !exists {
			out[key] = val
		}
	}
	return out
}

// kvKeys returns the keys of a slog KV map in a stable order: standard fields
// first (time, level, msg), then the rest in the order they were encountered.
// Because Go maps are unordered we reconstruct order from the original line.
func kvKeys(kv map[string]string) []string {
	priority := []string{
		"time",
		"ts",
		"timestamp",
		"Time",
		"level",
		"lvl",
		"severity",
		"msg",
		"message",
		"Message",
	}
	seen := make(map[string]bool)
	var out []string
	for _, k := range priority {
		if _, ok := kv[k]; ok {
			out = append(out, k)
			seen[k] = true
		}
	}
	for k := range kv {
		if !seen[k] {
			out = append(out, k)
		}
	}
	return out
}

// firstOf returns the first non-empty value found for any of the given keys.
func firstOf(m map[string]string, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k]; ok && v != "" {
			return v
		}
	}
	return ""
}

// renderJSONLog attempts to parse line as a structured log entry (slog/zap/
// zerolog) and renders it as a human-readable coloured string.
// Returns (rendered, true) on success, ("", false) if line is not valid JSON
// or doesn't look like a log entry.
//
// For known error patterns (database connection refused, etc.) an actionable
// hint is printed to os.Stdout immediately after the log line.
func renderJSONLog(line string) (string, string, bool) {
	var entry map[string]interface{}
	if err := json.Unmarshal([]byte(line), &entry); err != nil {
		return "", "", false
	}

	// ── Extract standard fields ───────────────────────────────────────────────
	level := extractString(entry, "level", "lvl", "severity")
	msg := extractString(entry, "msg", "message", "Message")
	ts := extractString(entry, "time", "ts", "timestamp", "Time")

	if msg == "" {
		return "", "", false
	}

	// ── Detect known error patterns and queue a hint ──────────────────────────
	// Collect all string values in the entry to search across msg + extra fields.
	allText := strings.ToLower(msg)
	for _, v := range entry {
		if s, ok := v.(string); ok {
			allText += " " + strings.ToLower(s)
		}
	}
	// Build allText — the combined lowercased text of all fields — returned to
	// the caller so it can call detectHints without needing to re-parse.

	// ── Level badge ───────────────────────────────────────────────────────────
	var levelBadge string
	switch strings.ToUpper(level) {
	case "DEBUG":
		levelBadge = ansiGray + ansiBold + "DBG" + ansiReset
	case "INFO":
		levelBadge = ansiGreen + ansiBold + "INF" + ansiReset
	case "WARN", "WARNING":
		levelBadge = ansiYellow + ansiBold + "WRN" + ansiReset
	case "ERROR", "ERR":
		levelBadge = ansiRed + ansiBold + "ERR" + ansiReset
	case "FATAL", "PANIC":
		levelBadge = ansiBgRed + " " + level + " " + ansiReset
	default:
		levelBadge = ansiGray + ansiBold + "LOG" + ansiReset
	}

	// ── Timestamp — keep only HH:MM:SS ───────────────────────────────────────
	timeStr := ""
	if ts != "" {
		// ISO-8601: take only the time portion HH:MM:SS
		if idx := strings.IndexByte(ts, 'T'); idx >= 0 && len(ts) > idx+9 {
			timeStr = ansiDim + ansiGray + ts[idx+1:idx+9] + ansiReset + "  "
		}
	}

	// ── Extra fields (skip standard ones already rendered) ───────────────────
	skip := map[string]bool{
		"level": true, "lvl": true, "severity": true,
		"msg": true, "message": true, "Message": true,
		"time": true, "ts": true, "timestamp": true, "Time": true,
	}

	var extras []string
	for k, v := range entry {
		if skip[k] {
			continue
		}
		var valStr string
		switch vv := v.(type) {
		case string:
			valStr = vv
		default:
			b, _ := json.Marshal(vv)
			valStr = string(b)
		}
		// Truncate very long values.
		if len(valStr) > 120 {
			valStr = valStr[:117] + "…"
		}
		extras = append(
			extras,
			ansiGray+k+"="+ansiReset+ansiDim+valStr+ansiReset,
		)
	}

	// ── Colour the message based on level ────────────────────────────────────
	var msgColour string
	switch strings.ToUpper(level) {
	case "ERROR", "ERR", "FATAL", "PANIC":
		msgColour = ansiRed
	case "WARN", "WARNING":
		msgColour = ansiYellow
	default:
		msgColour = ""
	}

	msgPart := msgColour + msg + ansiReset
	if msgColour == "" {
		msgPart = msg
	}

	extra := ""
	if len(extras) > 0 {
		extra = "  " + strings.Join(extras, "  ")
	}

	return fmt.Sprintf(
		"  %s%s%s  %s%s",
		timeStr,
		levelBadge,
		ansiReset,
		msgPart,
		extra,
	), allText, true
}

// extractString returns the first non-empty string value found for any of the
// given keys in the map.
func extractString(m map[string]interface{}, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k]; ok {
			if s, ok := v.(string); ok && s != "" {
				return s
			}
		}
	}
	return ""
}

// ── Known-error hint engine ───────────────────────────────────────────────────

// detectHints inspects the combined lowercased text of a log entry and prints
// an actionable hint block when a known error pattern is recognised.
// Each pattern is only printed once per rebuild via aw.hintSeen.
func (aw *appOutputWriter) detectHints(text string) {
	switch {
	case (strings.Contains(text, "connection refused") ||
		strings.Contains(text, "dial error") ||
		strings.Contains(text, "failed to connect") ||
		strings.Contains(text, "failed to initialize database")) &&
		!aw.hintSeen["db"]:

		aw.hintSeen["db"] = true
		printHint(aw.w,
			"Cannot connect to the database.",
			[]string{
				"# make sure your database is running:",
				"docker compose up -d",
				"# or check your DB_HOST / DB_PORT in .env",
			},
		)

	case strings.Contains(text, ".env") && strings.Contains(text, "not found") &&
		!aw.hintSeen["env"]:

		aw.hintSeen["env"] = true
		printHint(aw.w,
			"Environment file not found.",
			[]string{
				"cp .env.example .env",
				"# then edit .env with your database credentials",
			},
		)
	}
}

// printHint renders a styled actionable hint block to w.
func printHint(w *os.File, title string, steps []string) {
	fmt.Fprintln(w)
	fmt.Fprintf(w,
		"  %s  %s%s%s\n",
		badge(ansiBgYellow, "HINT"),
		ansiBold, title, ansiReset,
	)
	fmt.Fprintln(w)
	for _, s := range steps {
		if strings.HasPrefix(s, "#") {
			fmt.Fprintf(w, "    %s%s%s\n", ansiGray+ansiDim, s, ansiReset)
		} else {
			fmt.Fprintf(w, "    %s%s%s\n", ansiGreen, s, ansiReset)
		}
	}
	fmt.Fprintln(w)
}
