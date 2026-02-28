package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"
)

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

Pass ` + colorGreen + `-c` + colorReset + ` to display a per-suite coverage report.
Pass ` + colorGreen + `-w` + colorReset + ` to enter watch mode — specs re-run automatically on every
file change. Requires ` + colorCyan + `air` + colorReset + ` (` + colorGray + `go install github.com/air-verse/air@latest` + colorReset + `).
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
		"Re-run specs automatically on file changes (requires air)",
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

	if testWatch {
		return runTestWatch(testsDir)
	}

	return runTestOnce(testsDir)
}

// ──────────────────────────────────────────────
// One-shot run
// ──────────────────────────────────────────────

func runTestOnce(testsDir string) error {
	fmt.Println()
	fmt.Printf(
		"  %sUpdating gest%s %s\n",
		colorGray, colorReset,
		gray("(go get "+gestModule+")"),
	)
	fmt.Println()

	if err := ensureGest(); err != nil {
		fmt.Println(
			warn("Could not update gest — running with current version."),
		)
		fmt.Printf(
			"  %sRun manually: %s\n",
			colorGray,
			colorGreen+"go get "+gestModule+colorReset,
		)
	}

	goArgs := []string{"run", testsDir}
	if testCoverage {
		goArgs = append(goArgs, "-c")
	}

	label := "go run " + testsDir
	if testCoverage {
		label += " -c"
	}

	fmt.Println()
	fmt.Printf(
		"  %sRunning specs%s %s\n",
		colorGray, colorReset,
		gray("("+label+")"),
	)
	fmt.Println()

	c := exec.Command("go", goArgs...)
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
		// so we suppress the generic "exited with error" wrapper and let gest's
		// own output speak for itself.
		return fmt.Errorf("one or more specs failed")
	}

	return nil
}

// ──────────────────────────────────────────────
// Watch mode
// ──────────────────────────────────────────────

func runTestWatch(testsDir string) error {
	if _, err := exec.LookPath("air"); err != nil {
		return fmt.Errorf(
			"air not found in PATH\n\n"+
				"  Install it with: %s\n\n"+
				"  Then re-run: %s",
			colorCyan+"go install github.com/air-verse/air@latest"+colorReset,
			colorGreen+"grove test -w"+colorReset,
		)
	}

	// Build output goes into .grove_tmp so it is always gitignored.
	bin := filepath.Join(".grove_tmp", "grove_tests")
	fullBin := bin
	if testCoverage {
		fullBin = bin + " -c"
	}

	airCfg := buildAirConfig(bin, fullBin)

	// Write config to a temp file so we don't litter the project root.
	tmpFile, err := os.CreateTemp("", "grove-air-*.toml")
	if err != nil {
		return fmt.Errorf("failed to create temp air config: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(airCfg); err != nil {
		tmpFile.Close()
		return fmt.Errorf("failed to write air config: %w", err)
	}
	tmpFile.Close()

	fmt.Println()
	fmt.Printf(
		"  %s WATCH %s  Watching for changes — press %s to stop\n",
		colorBgGreen, colorReset,
		bold("Ctrl+C"),
	)
	if testCoverage {
		fmt.Printf(
			"  %sCoverage report enabled%s\n",
			colorGray, colorReset,
		)
	}
	fmt.Println()

	c := exec.Command("air", "-c", tmpFile.Name())
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	if err := c.Start(); err != nil {
		return fmt.Errorf("failed to start air: %w", err)
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
		return fmt.Errorf("air exited with error: %w", err)
	}

	return nil
}

// buildAirConfig returns an air TOML configuration that builds the tests
// binary on every change and re-runs it.
func buildAirConfig(bin, fullBin string) string {
	return `root = "."
tmp_dir = ".grove_tmp"

[build]
  cmd        = "go build -o ` + bin + ` ./internal/tests/"
  bin        = "` + bin + `"
  full_bin   = "` + fullBin + `"
  include_ext = ["go"]
  exclude_dir = ["tmp", "bin", "vendor", ".git", ".grove_tmp"]
  delay      = 400
  kill_delay = 200

[log]
  time = false

[color]
  main    = "magenta"
  watcher = "cyan"
  build   = "yellow"
  runner  = "green"

[misc]
  clean_on_exit = true
`
}
