package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

// testBin is the path to the compiled test binary.
const testBin = ".grove/tmp/tests"

// ──────────────────────────────────────────────
// make:test
// ──────────────────────────────────────────────

var makeTestCmd = &cobra.Command{
	Use:   "make:test <Name>",
	Short: "Scaffold a new gest spec file",
	Long: bold(
		"make:test",
	) + ` scaffolds a new gest spec file in ` + colorCyan + `internal/tests/` + colorReset + `.

If ` + colorCyan + `internal/tests/main.go` + colorReset + ` does not exist it will be created
automatically as the gest entrypoint.

` + colorGray + `Examples:` + colorReset + `
  grove make:test User
  grove make:test AuthService
  grove make:test order_calculations`,
	Args: cobra.ExactArgs(1),
	RunE: runMakeTest,
}

func runMakeTest(_ *cobra.Command, args []string) error {
	name := toPascalCase(args[0])

	fmt.Println()
	fmt.Printf(
		"  %sCreating test spec%s %s\n",
		colorGray, colorReset,
		bold(name),
	)
	fmt.Println()

	// Ensure the gest entrypoint exists before creating the spec.
	if err := scaffoldTestMain(); err != nil {
		return err
	}

	if err := scaffoldTestSpec(name); err != nil {
		return err
	}

	snake := toSnakeCase(name)

	fmt.Println()
	fmt.Println(nextSteps())
	fmt.Printf(
		"    %s1.%s Write your assertions in %s\n",
		colorGray, colorReset,
		colorCyan+"internal/tests/"+snake+"_spec.go"+colorReset,
	)
	fmt.Printf(
		"    %s2.%s Run %s to execute all specs\n",
		colorGray, colorReset,
		colorGreen+"grove test"+colorReset,
	)
	fmt.Printf(
		"    %s3.%s Run %s to watch and re-run on every save\n",
		colorGray, colorReset,
		colorGreen+"grove test -w"+colorReset,
	)
	fmt.Println()

	return nil
}

// ──────────────────────────────────────────────
// test
// ──────────────────────────────────────────────

var (
	testCoverage bool
	testWatch    bool
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run all gest specs in internal/tests",
	Long: bold(
		"test",
	) + ` compiles and runs every ` + colorCyan + `*_spec.go` + colorReset + ` file found in
` + colorCyan + `internal/tests/` + colorReset + ` using the ` + colorCyan + `gest` + colorReset + ` testing framework.

The test suite is compiled once to ` + colorCyan + `.grove/tmp/tests` + colorReset + ` and executed
directly — subsequent runs (and watch-mode rebuilds) only recompile what
changed, making them significantly faster than ` + colorGray + `go run` + colorReset + `.

Pass ` + colorGreen + `-c` + colorReset + ` to display a per-suite coverage report.
Pass ` + colorGreen + `-w` + colorReset + ` to enter watch mode — specs re-run automatically on every
file change. No external tools required.
Flags can be combined: ` + colorGreen + `-wc` + colorReset + ` runs watch mode with the coverage report.

` + colorGray + `Examples:` + colorReset + `
  grove test
  grove test -c
  grove test -w
  grove test -wc`,
	RunE: runTest,
}

func init() {
	testCmd.Flags().BoolVarP(
		&testCoverage,
		"coverage", "c", false,
		"Display a per-suite coverage report",
	)
	testCmd.Flags().BoolVarP(
		&testWatch,
		"watch", "w", false,
		"Re-run specs automatically on file changes (built-in, no external tools)",
	)
}

func runTest(_ *cobra.Command, _ []string) error {
	const testsDir = "./internal/tests"

	if _, err := os.Stat(testsDir); os.IsNotExist(err) {
		return fmt.Errorf(
			"tests directory not found: %s\n\n"+
				"  Create your first spec with: %s",
			colorCyan+testsDir+colorReset,
			colorGreen+"grove make:test <Name>"+colorReset,
		)
	}

	// Ensure the tmp directory exists for the compiled binary.
	if err := os.MkdirAll(filepath.Dir(testBin), 0o755); err != nil {
		return fmt.Errorf("cannot create tmp dir: %w", err)
	}

	if testWatch {
		return runTestWatch(testsDir)
	}

	return runTestOnce(testsDir)
}

// ──────────────────────────────────────────────
// Build helper
// ──────────────────────────────────────────────

// buildTests compiles internal/tests into the testBin binary.
// It returns the stderr output on failure so the caller can display it.
func buildTests(testsDir string) error {
	var stderr bytes.Buffer
	c := exec.Command("go", "build", "-o", testBin, testsDir)
	c.Stderr = &stderr
	c.Stdout = nil

	if err := c.Run(); err != nil {
		// Print compiler errors through the build writer for coloured output.
		bw := newBuildOutputWriter(os.Stderr)
		_, _ = bw.Write(stderr.Bytes())
		return fmt.Errorf("build failed")
	}
	return nil
}

// ──────────────────────────────────────────────
// One-shot run
// ──────────────────────────────────────────────

func runTestOnce(testsDir string) error {
	fmt.Println()
	fmt.Printf(
		"  %sBuilding specs%s %s\n",
		colorGray, colorReset,
		gray("(go build "+testsDir+" → "+testBin+")"),
	)
	fmt.Println()

	start := time.Now()
	if err := buildTests(testsDir); err != nil {
		fmt.Println()
		fmt.Printf("  %s\n", badge(colorBgRed, "BUILD FAILED"))
		fmt.Println()
		return err
	}
	buildElapsed := time.Since(start)

	fmt.Printf(
		"  %sRunning specs%s %s\n",
		colorGray, colorReset,
		gray(fmt.Sprintf("(built in %s)", fmtTestElapsed(buildElapsed))),
	)
	fmt.Println()

	var runArgs []string
	if testCoverage {
		runArgs = append(runArgs, "-c")
	}

	c := exec.Command(testBin, runArgs...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	if err := c.Start(); err != nil {
		return fmt.Errorf("failed to start tests: %w", err)
	}

	go func() {
		sig := <-sigCh
		if c.Process != nil {
			_ = c.Process.Signal(sig)
		}
	}()

	if err := c.Wait(); err != nil {
		if isSignalError(err) {
			fmt.Println()
			fmt.Println(gray("  Tests stopped."))
			fmt.Println()
			return nil
		}
		// gest exits with code 1 when tests fail — that's expected behaviour,
		// so we suppress the generic error wrapper and let gest's own output
		// speak for itself.
		return fmt.Errorf("one or more specs failed")
	}

	return nil
}

// ──────────────────────────────────────────────
// Watch mode (built-in, no external tools)
// ──────────────────────────────────────────────

// testWatcher holds the state for the built-in test watch loop.
type testWatcher struct {
	testsDir string

	mu       sync.Mutex
	debounce *time.Timer
	proc     *exec.Cmd     // currently running test binary
	procDone chan struct{} // closed when proc exits
	buildCh  chan struct{} // signals the build worker
}

func runTestWatch(testsDir string) error {
	fsw, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("cannot create watcher: %w", err)
	}
	defer fsw.Close() //nolint:errcheck

	// Watch the whole project for .go file changes, not just internal/tests,
	// so that changes to models, services etc. also trigger a re-run.
	if err := addTestDirRecursive(fsw, "."); err != nil {
		return err
	}

	tw := &testWatcher{
		testsDir: testsDir,
		buildCh:  make(chan struct{}, 1),
	}

	fmt.Println()
	fmt.Printf(
		"  %s WATCH %s  Watching for changes — press %s to stop\n",
		colorBgGreen, colorReset,
		bold("Ctrl+C"),
	)
	if testCoverage {
		fmt.Printf("  %sCoverage report enabled%s\n", colorGray, colorReset)
	}
	fmt.Println()

	// Trigger an initial build+run immediately.
	tw.buildCh <- struct{}{}

	// Build worker — serialises all rebuild+run cycles.
	go func() {
		for range tw.buildCh {
			tw.buildAndRun()
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	for {
		select {
		case event, ok := <-fsw.Events:
			if !ok {
				return nil
			}
			// Auto-watch newly created subdirectories.
			if event.Has(fsnotify.Create) {
				if info, err := os.Stat(
					event.Name,
				); err == nil &&
					info.IsDir() {
					_ = addTestDirRecursive(fsw, event.Name)
				}
			}
			if shouldHandleTestEvent(event) {
				tw.scheduleRebuild()
			}

		case err, ok := <-fsw.Errors:
			if !ok {
				return nil
			}
			fmt.Printf(
				"  %swatcher error: %s%s\n",
				colorYellow,
				err.Error(),
				colorReset,
			)

		case <-sigCh:
			tw.stopProc()
			fmt.Println()
			fmt.Println(gray("  Watch stopped."))
			fmt.Println()
			return nil
		}
	}
}

// scheduleRebuild debounces rapid saves into a single rebuild.
func (tw *testWatcher) scheduleRebuild() {
	tw.mu.Lock()
	defer tw.mu.Unlock()

	if tw.debounce != nil {
		tw.debounce.Stop()
	}
	tw.debounce = time.AfterFunc(100*time.Millisecond, func() {
		select {
		case tw.buildCh <- struct{}{}:
		default:
		}
	})
}

// buildAndRun stops any running test binary, rebuilds, and runs.
func (tw *testWatcher) buildAndRun() {
	tw.stopProc()

	fmt.Println()
	fmt.Printf(
		"  %s%s%s  %s\n",
		colorBgBlue, " BUILDING ", colorReset,
		gray("(go build "+tw.testsDir+" → "+testBin+")"),
	)

	start := time.Now()
	if err := buildTests(tw.testsDir); err != nil {
		fmt.Println()
		fmt.Printf("  %s\n", badge(colorBgRed, "BUILD FAILED"))
		fmt.Println()
		return
	}
	buildElapsed := time.Since(start)

	fmt.Printf(
		"  %s%s%s  %s\n\n",
		colorBgGreen, " RUNNING ", colorReset,
		gray(fmt.Sprintf("built in %s", fmtTestElapsed(buildElapsed))),
	)

	var runArgs []string
	if testCoverage {
		runArgs = append(runArgs, "-c")
	}

	proc := exec.Command(testBin, runArgs...)
	proc.Stdout = os.Stdout
	proc.Stderr = os.Stderr
	proc.Stdin = os.Stdin

	if err := proc.Start(); err != nil {
		fmt.Printf(
			"  %sfailed to start tests: %s%s\n",
			colorRed,
			err.Error(),
			colorReset,
		)
		return
	}

	done := make(chan struct{})

	tw.mu.Lock()
	tw.proc = proc
	tw.procDone = done
	tw.mu.Unlock()

	go func() {
		_ = proc.Wait()
		close(done)
	}()
}

// stopProc sends SIGINT to the running test binary and waits for it to exit.
func (tw *testWatcher) stopProc() {
	tw.mu.Lock()
	proc := tw.proc
	done := tw.procDone
	tw.proc = nil
	tw.procDone = nil
	tw.mu.Unlock()

	if proc == nil || proc.Process == nil {
		return
	}

	_ = proc.Process.Signal(os.Interrupt)

	select {
	case <-done:
	case <-time.After(3 * time.Second):
		_ = proc.Process.Kill()
		if done != nil {
			<-done
		}
	}
}

// ──────────────────────────────────────────────
// fsnotify helpers
// ──────────────────────────────────────────────

var testExcludeDirs = map[string]bool{
	".git":         true,
	".grove":       true,
	"vendor":       true,
	"node_modules": true,
}

func addTestDirRecursive(fsw *fsnotify.Watcher, root string) error {
	return filepath.WalkDir(
		root,
		func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if d.IsDir() {
				if testExcludeDirs[d.Name()] {
					return filepath.SkipDir
				}
				return fsw.Add(path)
			}
			return nil
		},
	)
}

func shouldHandleTestEvent(event fsnotify.Event) bool {
	if !event.Has(fsnotify.Write) && !event.Has(fsnotify.Create) {
		return false
	}
	return filepath.Ext(event.Name) == ".go"
}

// ──────────────────────────────────────────────
// Formatting helpers
// ──────────────────────────────────────────────

func fmtTestElapsed(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%.0fms", float64(d)/float64(time.Millisecond))
	}
	return fmt.Sprintf("%.1fs", d.Seconds())
}
