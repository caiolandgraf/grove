package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
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
		"    %s3.%s Run %s to include a coverage report\n",
		colorGray, colorReset,
		colorGreen+"grove test -c"+colorReset,
	)
	fmt.Println()

	return nil
}

// ──────────────────────────────────────────────
// test
// ──────────────────────────────────────────────

var testCoverage bool

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run all gest specs in internal/tests",
	Long: bold(
		"test",
	) + ` compiles and runs every ` + colorCyan + `*_spec.go` + colorReset + ` file found in
` + colorCyan + `internal/tests/` + colorReset + ` using the ` + colorCyan + `gest` + colorReset + ` testing framework.

Pass ` + colorGreen + `-c` + colorReset + ` to display a per-suite coverage report after the run.

` + colorGray + `Examples:` + colorReset + `
  grove test
  grove test -c
  grove test --coverage`,
	RunE: runTest,
}

func init() {
	testCmd.Flags().BoolVarP(
		&testCoverage,
		"coverage", "c", false,
		"Display a per-suite coverage report",
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
