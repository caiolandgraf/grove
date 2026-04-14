package main

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	makeSeedRunnerPath string
)

var makeSeedRunnerCmd = &cobra.Command{
	Use:   "make:seed",
	Short: "Scaffold a seed runner entrypoint (cmd/seed)",
	Long: bold(
		"make:seed",
	) + ` scaffolds a seed runner entrypoint at ` + colorCyan + `cmd/seed/main.go` + colorReset + `.

This is used by ` + colorGreen + `grove db:seed` + colorReset + ` to execute your project's seeders.

The generated runner should initialize app globals and run:

  internal/database/seeders.Run(app.DB)

` + colorGray + `Examples:` + colorReset + `
  grove make:seed
  grove make:seed --path ./cmd/seed`,
	Args: cobra.NoArgs,
	RunE: runMakeSeedRunner,
}

func init() {
	makeSeedRunnerCmd.Flags().StringVar(
		&makeSeedRunnerPath,
		"path",
		filepath.Join("cmd", "seed"),
		"Target directory for the seed runner (main.go will be created inside it)",
	)
}

func runMakeSeedRunner(_ *cobra.Command, _ []string) error {
	destDir := makeSeedRunnerPath
	destPath := filepath.Join(destDir, "main.go")

	fmt.Println()
	fmt.Printf(
		"  %sCreating seed runner%s %s\n",
		colorGray, colorReset,
		gray("→ "+destPath),
	)
	fmt.Println()

	if fileExists(destPath) {
		printSkipped("SeedRunner", "cmd/seed", destPath)
		fmt.Println()
		return nil
	}

	if err := ensureDir(destDir); err != nil {
		return fmt.Errorf("failed to create %s: %w", destDir, err)
	}

	module := getModuleName()

	data := struct {
		Module string
	}{
		Module: module,
	}

	content, err := renderStub(seedRunnerStub, "seed_runner", data)
	if err != nil {
		return err
	}

	if err := writeFile(destPath, content); err != nil {
		return err
	}

	printCreated("SeedRunner", "cmd/seed", destPath)

	fmt.Println()
	fmt.Println(nextSteps())
	fmt.Printf(
		"    %s1.%s Ensure you have seeders in %s\n",
		colorGray, colorReset,
		colorCyan+"internal/database/seeders/"+colorReset,
	)
	fmt.Printf(
		"    %s2.%s Run %s to execute them\n",
		colorGray, colorReset,
		colorGreen+"grove db:seed"+colorReset,
	)
	fmt.Println()

	return nil
}
