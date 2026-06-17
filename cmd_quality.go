package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

const qualityPackages = "./..."

var errQualitySkip = errors.New("quality step skipped")

var (
	fmtCmd = &cobra.Command{
		Use:   "fmt",
		Short: "Format code with golangci-lint and golines",
		Long: bold("fmt") + ` formats all Go packages in the project.

Runs ` + colorCyan + `golangci-lint fmt ./...` + colorReset + ` followed by
` + colorCyan + `golines -w -m 120 .` + colorReset + ` when golines is installed.

` + colorGray + `Examples:` + colorReset + `
  grove fmt`,
		RunE: runFmt,
	}

	lintCmd = &cobra.Command{
		Use:   "lint",
		Short: "Run golangci-lint",
		Long: bold("lint") + ` runs ` + colorCyan + `golangci-lint run ./...` + colorReset + `
using the project's ` + colorCyan + `.golangci.yml` + colorReset + ` config.

` + colorGray + `Examples:` + colorReset + `
  grove lint`,
		RunE: runLint,
	}

	lintFixCmd = &cobra.Command{
		Use:   "lint:fix",
		Short: "Auto-fix lint issues",
		Long: bold("lint:fix") + ` runs ` + colorCyan + `golangci-lint run --fix ./...` + colorReset + `
to auto-fix lint and format issues where possible.

` + colorGray + `Examples:` + colorReset + `
  grove lint:fix`,
		RunE: runLintFix,
	}

	checkCmd = &cobra.Command{
		Use:   "check",
		Short: "Run CI checks without modifying files",
		Long: bold("check") + ` runs lint, unit tests, and builds project binaries
without modifying source files — the read-only CI pipeline.

` + colorGray + `Pipeline:` + colorReset + `
  1. ` + colorCyan + `golangci-lint run ./...` + colorReset + `
  2. ` + colorCyan + `go test -race -count=1 -short ./...` + colorReset + `
  3. ` + colorCyan + `go build` + colorReset + ` → ` + colorCyan + `.grove/bin/api` + colorReset + ` and ` + colorCyan + `.grove/bin/atlas` + colorReset + `

` + colorGray + `Examples:` + colorReset + `
  grove check`,
		RunE: runCheck,
	}

	buildBinariesCmd = &cobra.Command{
		Use:   "build:binaries",
		Short: "Build api and atlas binaries to .grove/bin",
		Long: bold("build:binaries") + ` compiles ` + colorCyan + `cmd/api` + colorReset + ` and
` + colorCyan + `cmd/atlas` + colorReset + ` into ` + colorCyan + `.grove/bin/` + colorReset + `.

` + colorGray + `Examples:` + colorReset + `
  grove build:binaries`,
		RunE: runBuildBinaries,
	}
)

func requireGoModule() error {
	if _, err := os.Stat("go.mod"); err != nil {
		return fmt.Errorf(
			"not a Go module — run this command from your project root (go.mod not found)",
		)
	}
	return nil
}

func qualityTool(envVar, defaultName string) string {
	if v := os.Getenv(envVar); v != "" {
		return v
	}
	return defaultName
}

func requireQualityTool(envVar, defaultName, installHint string) (string, error) {
	name := qualityTool(envVar, defaultName)
	if _, err := exec.LookPath(name); err != nil {
		return "", fmt.Errorf(
			"%s not found in PATH — %s",
			colorCyan+name+colorReset,
			installHint,
		)
	}
	return name, nil
}

func runQualityCommand(name string, args ...string) error {
	c := exec.Command(name, args...)
	c.Stdout = newIndentWriter(os.Stdout, "    ")
	c.Stderr = newIndentWriter(os.Stderr, "    ")
	return c.Run()
}

type qualityStep struct {
	label     string
	run       func() error
	skippable bool
}

func runQualityPipeline(
	badgeLabel, subtitle, doneMsg string,
	steps []qualityStep,
) error {
	if err := requireGoModule(); err != nil {
		return err
	}

	fmt.Println()
	fmt.Printf("  %s  %s\n", badge(colorBgBlue, badgeLabel), gray(subtitle))
	fmt.Println()

	for i, step := range steps {
		fmt.Printf(
			"  %s%d/%d%s  %s\n",
			colorGray, i+1, len(steps), colorReset,
			step.label,
		)
		fmt.Println()

		if err := step.run(); err != nil {
			if step.skippable && errors.Is(err, errQualitySkip) {
				fmt.Println()
				continue
			}

			fmt.Println()
			fmt.Printf("  %s\n", badge(colorBgRed, badgeLabel+" FAILED"))
			fmt.Printf("  %s  %s\n", colorRed+"✕"+colorReset, step.label)
			fmt.Println()
			return err
		}

		fmt.Println(success(step.label + " passed"))
		fmt.Println()
	}

	fmt.Println(done(doneMsg))
	fmt.Println()

	return nil
}

func runProjectFmt() error {
	lint, err := requireQualityTool(
		"GOLANGCI_LINT", "golangci-lint",
		"install from https://golangci-lint.run/welcome/install/",
	)
	if err != nil {
		return err
	}
	return runQualityCommand(lint, "fmt", qualityPackages)
}

func runProjectGolines() error {
	name := qualityTool("GOLINES", "golines")
	if _, err := exec.LookPath(name); err != nil {
		fmt.Printf(
			"    %s\n",
			gray("golines not found; skipped (install or add to PATH to match editor formatting)"),
		)
		return errQualitySkip
	}
	return runQualityCommand(name, "-w", "-m", "120", ".")
}

func runProjectLint() error {
	lint, err := requireQualityTool(
		"GOLANGCI_LINT", "golangci-lint",
		"install from https://golangci-lint.run/welcome/install/",
	)
	if err != nil {
		return err
	}
	return runQualityCommand(lint, "run", qualityPackages)
}

func runProjectLintFix() error {
	lint, err := requireQualityTool(
		"GOLANGCI_LINT", "golangci-lint",
		"install from https://golangci-lint.run/welcome/install/",
	)
	if err != nil {
		return err
	}
	return runQualityCommand(lint, "run", "--fix", qualityPackages)
}

func runProjectGoTest(short bool) error {
	goBin, err := requireQualityTool("GO", "go", "install Go from https://go.dev/dl/")
	if err != nil {
		return err
	}

	args := []string{"test", "-race", "-count=1", qualityPackages}
	if short {
		args = append(args, "-short")
	}
	return runQualityCommand(goBin, args...)
}

func runProjectBuild() error {
	goBin, err := requireQualityTool("GO", "go", "install Go from https://go.dev/dl/")
	if err != nil {
		return err
	}

	if err := ensureDir(".grove/bin"); err != nil {
		return fmt.Errorf("failed to create .grove/bin: %w", err)
	}

	ldflags := "-s -w"
	targets := []struct {
		out string
		pkg string
	}{
		{".grove/bin/api", "./cmd/api"},
		{".grove/bin/atlas", "./cmd/atlas"},
	}

	for _, t := range targets {
		if err := runQualityCommand(
			goBin, "build", "-trimpath", "-ldflags="+ldflags, "-o", t.out, t.pkg,
		); err != nil {
			return fmt.Errorf("build %s: %w", t.out, err)
		}
	}

	return nil
}

func runFmt(_ *cobra.Command, _ []string) error {
	return runQualityPipeline(
		"FORMAT",
		"golangci-lint fmt · golines",
		"Code formatted.",
		[]qualityStep{
			{"Format (golangci-lint fmt)", runProjectFmt, false},
			{"Format (golines)", runProjectGolines, true},
		},
	)
}

func runLint(_ *cobra.Command, _ []string) error {
	if err := requireGoModule(); err != nil {
		return err
	}

	fmt.Println()
	fmt.Printf("  %s  %s\n", badge(colorBgBlue, "LINT"), gray("golangci-lint run ./..."))
	fmt.Println()

	if err := runProjectLint(); err != nil {
		fmt.Println()
		fmt.Printf("  %s\n", badge(colorBgRed, "LINT FAILED"))
		fmt.Println()
		return err
	}

	fmt.Println()
	fmt.Println(done("Lint passed."))
	fmt.Println()

	return nil
}

func runLintFix(_ *cobra.Command, _ []string) error {
	if err := requireGoModule(); err != nil {
		return err
	}

	fmt.Println()
	fmt.Printf("  %s  %s\n", badge(colorBgBlue, "LINT FIX"), gray("golangci-lint run --fix ./..."))
	fmt.Println()

	if err := runProjectLintFix(); err != nil {
		fmt.Println()
		fmt.Printf("  %s\n", badge(colorBgRed, "LINT FIX FAILED"))
		fmt.Println()
		return err
	}

	fmt.Println()
	fmt.Println(done("Lint issues auto-fixed."))
	fmt.Println()

	return nil
}

func runCheck(_ *cobra.Command, _ []string) error {
	return runQualityPipeline(
		"CHECK",
		"lint · test · build",
		"All CI checks passed.",
		[]qualityStep{
			{"Lint (golangci-lint run)", runProjectLint, false},
			{"Test (go test -race -short)", func() error { return runProjectGoTest(true) }, false},
			{"Build (.grove/bin/api, .grove/bin/atlas)", runProjectBuild, false},
		},
	)
}

func runBuildBinaries(_ *cobra.Command, _ []string) error {
	if err := requireGoModule(); err != nil {
		return err
	}

	fmt.Println()
	fmt.Printf(
		"  %s  %s\n",
		badge(colorBgBlue, "BUILDING"),
		gray("go build → .grove/bin/api, .grove/bin/atlas"),
	)
	fmt.Println()

	if err := runProjectBuild(); err != nil {
		fmt.Println()
		fmt.Printf("  %s\n", badge(colorBgRed, "BUILD FAILED"))
		fmt.Println()
		return err
	}

	fmt.Println()
	fmt.Println(done("Binaries compiled to " + colorCyan + ".grove/bin/" + colorReset + "."))
	fmt.Println()

	return nil
}
