package main

import (
	"github.com/spf13/cobra"
)

var prepareCmd = &cobra.Command{
	Use:   "prepare",
	Short: "Format, lint, test, and build before committing",
	Long: bold("prepare") + ` runs the full pre-commit quality pipeline:
format code, auto-fix lint issues, verify lint, run tests, and build binaries.

Run it from your Grove project root before every commit.

` + colorGray + `Pipeline:` + colorReset + `
  1. ` + colorCyan + `golangci-lint fmt ./...` + colorReset + `
  2. ` + colorCyan + `golines -w -m 120 .` + colorReset + ` (skipped when not installed)
  3. ` + colorCyan + `golangci-lint run --fix ./...` + colorReset + `
  4. ` + colorCyan + `golangci-lint run ./...` + colorReset + `
  5. ` + colorCyan + `go test -race -count=1 -short ./...` + colorReset + `
  6. ` + colorCyan + `go build` + colorReset + ` → ` + colorCyan + `.grove/bin/api` + colorReset + ` and ` + colorCyan + `.grove/bin/atlas` + colorReset + `

Environment variables ` + colorGray + `GO` + colorReset + `, ` + colorGray + `GOLANGCI_LINT` + colorReset + `, and
` + colorGray + `GOLINES` + colorReset + ` override the default tool binaries.

` + colorGray + `Examples:` + colorReset + `
  grove prepare`,
	RunE: runPrepare,
}

func runPrepare(_ *cobra.Command, _ []string) error {
	return runQualityPipeline(
		"PREPARE",
		"format · lint · test · build",
		"Project is ready to commit.",
		[]qualityStep{
			{"Format (golangci-lint fmt)", runProjectFmt, false},
			{"Format (golines)", runProjectGolines, true},
			{"Lint fix (golangci-lint run --fix)", runProjectLintFix, false},
			{"Lint (golangci-lint run)", runProjectLint, false},
			{"Test (go test -race -short)", func() error { return runProjectGoTest(true) }, false},
			{"Build (.grove/bin/api, .grove/bin/atlas)", runProjectBuild, false},
		},
	)
}
