package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Grove project dependencies",
	Long: bold("update") + ` updates the Grove project dependencies to their
latest versions and runs ` + colorCyan + `go mod tidy` + colorReset + ` to clean up the module graph.

Currently updates:
  ` + colorCyan + `github.com/caiolandgraf/gest` + colorReset + `  — the testing framework used by ` + colorGreen + `grove test` + colorReset + `

` + colorGray + `Examples:` + colorReset + `
  grove update`,
	RunE: runUpdate,
}

func runUpdate(_ *cobra.Command, _ []string) error {
	fmt.Println()

	// ── gest ──────────────────────────────────────────────────────────────────
	fmt.Printf(
		"  %sUpdating gest%s %s\n",
		colorGray, colorReset,
		gray("(go get "+gestModule+")"),
	)
	fmt.Println()

	gestCmd := exec.Command("go", "get", gestModule)
	gestCmd.Stdout = newIndentWriter(os.Stdout, "    ")
	gestCmd.Stderr = newIndentWriter(os.Stderr, "    ")

	if err := gestCmd.Run(); err != nil {
		fmt.Println(warn("Failed to update gest."))
		fmt.Printf(
			"  %sRun manually: %s\n",
			colorGray,
			colorGreen+"go get "+gestModule+colorReset,
		)
		fmt.Println()
		return fmt.Errorf("go get %s: %w", gestModule, err)
	}

	fmt.Println(success("gest updated"))
	fmt.Println()

	// ── go mod tidy ───────────────────────────────────────────────────────────
	fmt.Printf(
		"  %sCleaning module graph%s %s\n",
		colorGray, colorReset,
		gray("(go mod tidy)"),
	)
	fmt.Println()

	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Stdout = newIndentWriter(os.Stdout, "    ")
	tidyCmd.Stderr = newIndentWriter(os.Stderr, "    ")

	if err := tidyCmd.Run(); err != nil {
		fmt.Println(
			warn("go mod tidy failed — your go.mod may need manual attention."),
		)
		fmt.Println()
		return fmt.Errorf("go mod tidy: %w", err)
	}

	fmt.Println(success("go.mod tidied"))
	fmt.Println()
	fmt.Println(done("All dependencies are up to date."))
	fmt.Println()

	return nil
}
