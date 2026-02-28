package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
)

var buildOutput string

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Compile the application to a binary",
	Long: bold("build") + ` compiles the application and outputs the binary.

` + colorGray + `Examples:` + colorReset + `
  grove build
  grove build -o ./bin/my-api`,
	RunE: runBuild,
}

func init() {
	buildCmd.Flags().StringVarP(
		&buildOutput,
		"output", "o", "./bin/app",
		"Output path for the compiled binary",
	)
}

func runBuild(_ *cobra.Command, _ []string) error {
	fmt.Println()
	fmt.Printf(
		"  %s  %s\n",
		badge(colorBgBlue, "BUILDING"),
		gray("go build -o "+buildOutput+" ./cmd/api/"),
	)
	fmt.Println()

	if err := ensureDir("bin"); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	start := time.Now()

	bw := newBuildOutputWriter(os.Stderr)
	c := exec.Command("go", "build", "-o", buildOutput, "./cmd/api/")
	c.Stdout = bw
	c.Stderr = bw

	if err := c.Run(); err != nil {
		fmt.Println()
		fmt.Printf("  %s\n", badge(colorBgRed, "BUILD FAILED"))
		fmt.Println()
		return fmt.Errorf("")
	}

	elapsed := time.Since(start)

	fmt.Println()
	fmt.Println(done(
		"Binary compiled to " + colorCyan + buildOutput + colorReset +
			"  " + gray("("+fmtDuration(elapsed)+")"),
	))
	fmt.Println()

	return nil
}
