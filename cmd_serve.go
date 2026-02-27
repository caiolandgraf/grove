package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the development server",
	Long: bold("serve") + ` starts the development HTTP server.

If ` + colorCyan + `air` + colorReset + ` is installed it will be used for hot-reload.
Otherwise falls back to ` + colorCyan + `go run ./cmd/api/main.go` + colorReset + `.

` + colorGray + `Examples:` + colorReset + `
  grove serve`,
	RunE: runServe,
}

func runServe(_ *cobra.Command, _ []string) error {
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
		fmt.Println()
		c = exec.Command("go", "run", "./cmd/api/main.go")
	}

	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin

	// Forward SIGINT / SIGTERM so the child process gets a chance to
	// shut down gracefully before grove itself exits.
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
		// Ignore the error produced by a signal-triggered exit so we
		// don't print a confusing message when the user hits Ctrl-C.
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
	// exec.ExitError wraps the exit status; a signal-killed process has no
	// numeric exit code on Unix — checking the string is the portable way.
	msg := err.Error()
	return msg == "signal: interrupt" ||
		msg == "signal: terminated" ||
		msg == "signal: killed"
}
