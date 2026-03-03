package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

var devAirCmd = &cobra.Command{
	Use:   "dev:air",
	Short: "Start the development server using Air for hot-reload",
	Long: bold(
		"dev:air",
	) + ` starts the development HTTP server using ` + colorCyan + `air` + colorReset + ` for hot-reload.

If ` + colorCyan + `air` + colorReset + ` is not installed it falls back to ` + colorCyan + `go run ./cmd/api/main.go` + colorReset + `.

For a zero-dependency hot-reload experience use ` + colorGreen + `grove dev` + colorReset + ` instead —
it has a built-in watcher that requires no external tools.

` + colorGray + `Examples:` + colorReset + `
  grove dev:air`,
	RunE: runDevAir,
}

func runDevAir(_ *cobra.Command, _ []string) error {
	fmt.Println()

	var c *exec.Cmd

	if _, err := exec.LookPath("air"); err == nil {
		fmt.Printf(
			"  %s AIR %s  Starting server with %s (hot-reload enabled)\n",
			colorBgGreen, colorReset,
			bold("air"),
		)
		fmt.Println()
		c = exec.Command("air")
	} else {
		fmt.Printf(
			"  %s INFO %s  %s not found — starting with %s\n",
			colorBgBlue, colorReset,
			colorCyan+"air"+colorReset,
			bold("go run ./cmd/api/main.go"),
		)
		fmt.Printf(
			"  %sTip: install air for hot-reload → %s\n",
			colorGray,
			colorCyan+"go install github.com/air-verse/air@latest"+colorReset,
		)
		fmt.Printf(
			"  %sTip: use %s for built-in hot-reload with no external tools\n",
			colorGray,
			colorGreen+"grove dev"+colorReset,
		)
		fmt.Println()
		c = exec.Command("go", "run", "./cmd/api/main.go")
	}

	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	if err := c.Start(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
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
		return fmt.Errorf("server exited with error: %w", err)
	}

	return nil
}

// isSignalError reports whether err was caused by a signal (e.g. SIGINT).
func isSignalError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return msg == "signal: interrupt" ||
		msg == "signal: terminated" ||
		msg == "signal: killed"
}
