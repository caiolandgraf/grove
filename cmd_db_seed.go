package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var (
	dbSeedPackage string
	dbSeedEnvFile string
)

var dbSeedCmd = &cobra.Command{
	Use:   "db:seed",
	Short: "Run database seeders",
	Long: bold(
		"db:seed",
	) + ` runs your project's database seeders by executing a dedicated seed runner.

By default, Grove runs:

  go run ./cmd/seed

This expects your project to provide a seed entrypoint at ` + colorCyan + `cmd/seed/` + colorReset + ` (option 1).
That runner should load config, connect to the database, and call:

  internal/app/database/seeders.Run(db)

` + colorGray + `Examples:` + colorReset + `
  grove db:seed
  grove db:seed --package ./cmd/seed
  grove db:seed --env-file .env`,
	Args: cobra.NoArgs,
	RunE: runDBSeed,
}

func init() {
	dbSeedCmd.Flags().StringVar(
		&dbSeedPackage,
		"package",
		"./cmd/seed",
		"Go package/path to run for seeding (e.g. ./cmd/seed)",
	)
	dbSeedCmd.Flags().StringVar(
		&dbSeedEnvFile,
		"env-file",
		".env",
		"Env file to load before running seeders (only if present)",
	)
}

func runDBSeed(_ *cobra.Command, _ []string) error {
	fmt.Println()
	fmt.Printf(
		"  %sSeeding database%s %s\n",
		colorGray,
		colorReset,
		gray("(go run "+dbSeedPackage+")"),
	)
	fmt.Println()

	// Load env vars into the current process, so they propagate to `go run`.
	// (Does not override existing environment variables.)
	if dbSeedEnvFile != "" {
		if err := loadDotEnv(dbSeedEnvFile); err != nil {
			return fmt.Errorf("failed to load %s: %w", dbSeedEnvFile, err)
		}
	}

	c := exec.Command("go", "run", dbSeedPackage)
	c.Stdout = newIndentWriter(os.Stdout, "  ")
	c.Stderr = newIndentWriter(os.Stderr, "  ")
	c.Stdin = os.Stdin

	if err := c.Run(); err != nil {
		return fmt.Errorf("failed to run seeders: %w", err)
	}

	fmt.Println()
	fmt.Println(done("Seeders executed successfully."))
	fmt.Println()

	return nil
}
