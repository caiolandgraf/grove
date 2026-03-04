package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

// ──────────────────────────────────────────────
// make:test
// ──────────────────────────────────────────────

var makeTestCmd = &cobra.Command{
	Use:   "make:test <Name>",
	Short: "Scaffold a new gest test file",
	Long: bold(
		"make:test",
	) + ` scaffolds a new gest test file in ` + colorCyan + `internal/tests/` + colorReset + `.

The generated file follows the ` + colorCyan + `gest v2` + colorReset + ` convention: a standard
` + colorCyan + `*_test.go` + colorReset + ` file with a ` + colorCyan + `func Test<Name>(t *testing.T)` + colorReset + ` entry point
that calls ` + colorCyan + `s.Run(t)` + colorReset + ` — fully compatible with ` + colorGray + `go test` + colorReset + `.

On the first call, gest is added to the project's ` + colorCyan + `go.mod` + colorReset + ` automatically
via ` + colorGray + `go get` + colorReset + `.

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
		"  %sCreating test%s %s\n",
		colorGray, colorReset,
		bold(name),
	)
	fmt.Println()

	if err := scaffoldTestSpec(name); err != nil {
		return err
	}

	snake := toSnakeCase(name)

	fmt.Println()
	fmt.Println(nextSteps())
	fmt.Printf(
		"    %s1.%s Write your assertions in %s\n",
		colorGray, colorReset,
		colorCyan+"internal/tests/"+snake+"_test.go"+colorReset,
	)
	fmt.Printf(
		"    %s2.%s Run %s to execute all tests\n",
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
	Short: "Run all gest tests in internal/tests",
	Long: bold(
		"test",
	) + ` runs every ` + colorCyan + `*_test.go` + colorReset + ` file found in
` + colorCyan + `internal/tests/` + colorReset + ` using the ` + colorCyan + `gest` + colorReset + ` CLI for beautiful Jest-style output.

If the ` + colorCyan + `gest` + colorReset + ` CLI is not installed, grove falls back to ` + colorGray + `go test -v` + colorReset + `
automatically. Install it for the full experience:
  ` + colorGray + `go install github.com/caiolandgraf/gest/v2/cmd/gest@latest` + colorReset + `

Pass ` + colorGreen + `-c` + colorReset + ` to display a per-suite coverage report after the run.
Pass ` + colorGreen + `-w` + colorReset + ` to enter watch mode — tests re-run automatically on every
file change. Flags can be combined: ` + colorGreen + `-wc` + colorReset + `.

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
		"Re-run tests automatically on file changes",
	)
}

func runTest(_ *cobra.Command, _ []string) error {
	const testsDir = "./internal/tests"

	if _, err := os.Stat(testsDir); os.IsNotExist(err) {
		return fmt.Errorf(
			"tests directory not found: %s\n\n"+
				"  Create your first test with: %s",
			colorCyan+testsDir+colorReset,
			colorGreen+"grove make:test <Name>"+colorReset,
		)
	}

	if testWatch {
		return runTestWatch()
	}

	return runTestOnce()
}

// ──────────────────────────────────────────────
// One-shot run
// ──────────────────────────────────────────────

// runTestOnce runs the test suite once, preferring the gest CLI and falling
// back to plain `go test -v` when gest is not installed.
func runTestOnce() error {
	cmd, args := buildTestCommand()

	c := exec.Command(cmd, args...)
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
		// gest / go test exits with code 1 on failures — suppress the wrapper
		// and let the runner's own output speak for itself.
		return fmt.Errorf("one or more tests failed")
	}

	return nil
}

// ──────────────────────────────────────────────
// Watch mode
// ──────────────────────────────────────────────

// runTestWatch enters watch mode. When the gest CLI is available it delegates
// to `gest --watch [-c] ./internal/tests/...`. Otherwise it falls back to a
// simple polling loop that re-runs `go test -v` on every .go file change.
func runTestWatch() error {
	gestPath, gestAvailable := resolveGestCLI()
	if gestAvailable {
		return runGestWatch(gestPath)
	}
	return runGoTestWatchLoop()
}

// runGestWatch delegates watch mode to the gest CLI binary.
func runGestWatch(gestPath string) error {
	args := []string{"--watch"}
	if testCoverage {
		args = append(args, "-c")
	}
	args = append(args, "./internal/tests/...")

	c := exec.Command(gestPath, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	if err := c.Start(); err != nil {
		return fmt.Errorf("failed to start gest watch: %w", err)
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
			fmt.Println(gray("  Watch stopped."))
			fmt.Println()
			return nil
		}
		// test failures inside watch mode are not CLI errors
		return nil
	}

	return nil
}

// runGoTestWatchLoop is the fallback watch implementation used when the gest
// CLI is not installed. It polls for .go file changes every 500 ms and
// re-runs `go test -v ./internal/tests/...` on each change.
func runGoTestWatchLoop() error {
	fmt.Println()
	fmt.Printf(
		"  %s WATCH %s  Watching for changes — press %s to stop\n",
		colorBgGreen, colorReset,
		bold("Ctrl+C"),
	)
	fmt.Printf(
		"  %sTip: install the gest CLI for a better experience:%s\n",
		colorGray, colorReset,
	)
	fmt.Printf(
		"       %sgo install github.com/caiolandgraf/gest/v2/cmd/gest@latest%s\n\n",
		colorCyan,
		colorReset,
	)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	// Run once immediately.
	runGoTestOnce()

	snapshots := snapshotGoFiles(".")

	for {
		select {
		case <-sigCh:
			fmt.Println()
			fmt.Println(gray("  Watch stopped."))
			fmt.Println()
			return nil
		case <-time.After(500 * time.Millisecond):
			current := snapshotGoFiles(".")
			if snapshotsChanged(snapshots, current) {
				snapshots = current
				fmt.Print("\033[2J\033[3J\033[H")
				runGoTestOnce()
			}
		}
	}
}

// runGoTestOnce executes `go test -v ./internal/tests/...` synchronously.
func runGoTestOnce() {
	args := []string{"test", "-v"}
	if testCoverage {
		args = append(args, "-cover")
	}
	args = append(args, "./internal/tests/...")

	c := exec.Command("go", args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	_ = c.Run()
}

// ──────────────────────────────────────────────
// Snapshot helpers (polling watcher)
// ──────────────────────────────────────────────

var watchExcludeDirs = map[string]bool{
	".git":         true,
	".grove":       true,
	"vendor":       true,
	"node_modules": true,
}

// snapshotGoFiles returns a map of path → mtime (nanoseconds) for every .go
// file under root, skipping directories in watchExcludeDirs.
func snapshotGoFiles(root string) map[string]int64 {
	snap := map[string]int64{}
	_ = filepath.WalkDir(
		root,
		func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if d.IsDir() {
				if watchExcludeDirs[d.Name()] {
					return filepath.SkipDir
				}
				return nil
			}
			if filepath.Ext(path) == ".go" {
				if info, e := d.Info(); e == nil {
					snap[path] = info.ModTime().UnixNano()
				}
			}
			return nil
		},
	)
	return snap
}

// snapshotsChanged reports whether current differs from prev (new, deleted or
// modified files).
func snapshotsChanged(prev, current map[string]int64) bool {
	if len(prev) != len(current) {
		return true
	}
	for path, mtime := range current {
		if prev[path] != mtime {
			return true
		}
	}
	return false
}

// ──────────────────────────────────────────────
// Helpers
// ──────────────────────────────────────────────

// buildTestCommand returns the command and arguments to run the test suite
// once. Prefers the gest CLI; falls back to `go test -v`.
func buildTestCommand() (string, []string) {
	gestPath, ok := resolveGestCLI()
	if ok {
		args := []string{}
		if testCoverage {
			args = append(args, "-c")
		}
		args = append(args, "./internal/tests/...")
		return gestPath, args
	}

	// Fallback: plain go test
	args := []string{"test", "-v"}
	if testCoverage {
		args = append(args, "-cover")
	}
	args = append(args, "./internal/tests/...")
	return "go", args
}

// resolveGestCLI returns the path to the gest CLI binary and true when it is
// available on PATH.
func resolveGestCLI() (string, bool) {
	path, err := exec.LookPath("gest")
	if err != nil {
		return "", false
	}
	return path, true
}

// joinArgs joins a slice of strings into a single space-separated string for
// display purposes.
func joinArgs(args []string) string {
	result := ""
	for i, a := range args {
		if i > 0 {
			result += " "
		}
		result += a
	}
	return result
}
