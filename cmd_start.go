package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

var startOutput string

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Build and start the application binary",
	Long: bold(
		"start",
	) + ` compiles the application and runs the resulting binary.

` + colorGray + `Examples:` + colorReset + `
  grove start
  grove start -o ./bin/my-api`,
	RunE: runStart,
}

func init() {
	startCmd.Flags().StringVarP(
		&startOutput,
		"output", "o", "./bin/app",
		"Output path for the compiled binary",
	)
}

func runStart(_ *cobra.Command, _ []string) error {
	fmt.Println()
	fmt.Printf(
		"  %s  %s\n",
		badge(colorBgBlue, "BUILDING"),
		gray("go build -o "+startOutput+" ./cmd/api/"),
	)
	fmt.Println()

	elapsed, err := buildBinary(startOutput)
	if err != nil {
		fmt.Println()
		fmt.Printf("  %s\n", badge(colorBgRed, "BUILD FAILED"))
		fmt.Println()
		return fmt.Errorf("")
	}

	fmt.Println()
	fmt.Println(done(
		"Binary compiled to " + colorCyan + startOutput + colorReset +
			"  " + gray("("+fmtDuration(elapsed)+")"),
	))
	fmt.Println()

	fmt.Printf(
		"  %s  %s\n",
		badge(colorBgBlue, "STARTING"),
		gray(startOutput),
	)
	fmt.Println()

	c := exec.Command(startOutput)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	if err := c.Start(); err != nil {
		return fmt.Errorf("failed to start binary: %w", err)
	}

	go func() {
		sig := <-sigCh
		if c.Process != nil {
			_ = c.Process.Signal(sig)
		}
	}()

	if err := c.Wait(); err != nil {
		if c.ProcessState != nil && !c.ProcessState.Success() {
			if isSignalError(err) {
				fmt.Println()
				fmt.Println(gray("  Server stopped."))
				fmt.Println()
				return nil
			}
		}
		return fmt.Errorf("binary exited with error: %w", err)
	}

	return nil
}
