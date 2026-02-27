package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// ──────────────────────────────────────────────
// Shared flag
// ──────────────────────────────────────────────

var migrateEnv string

func addEnvFlag(cmd *cobra.Command) {
	cmd.Flags().StringVar(
		&migrateEnv,
		"env", "local",
		"Atlas environment to use (local, dev, production)",
	)
}

// ──────────────────────────────────────────────
// migrate
// ──────────────────────────────────────────────

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Apply all pending migrations",
	Long: bold("migrate") + ` applies all pending migrations to the database
using Atlas migrate apply.

` + colorGray + `Examples:` + colorReset + `
  grove migrate
  grove migrate --env dev
  grove migrate --env production`,
	RunE: runMigrate,
}

func init() {
	addEnvFlag(migrateCmd)
}

func runMigrate(cmd *cobra.Command, args []string) error {
	fmt.Println()
	fmt.Printf(
		"  %sRunning migrations%s %s\n",
		colorGray, colorReset,
		gray("(atlas migrate apply --env "+migrateEnv+")"),
	)
	fmt.Println()

	return runAtlas("apply pending migrations",
		"migrate", "apply",
		"--env", migrateEnv,
	)
}

// ──────────────────────────────────────────────
// migrate:rollback
// ──────────────────────────────────────────────

var (
	migrateRollbackEnv    string
	migrateRollbackAmount int
)

var migrateRollbackCmd = &cobra.Command{
	Use:   "migrate:rollback",
	Short: "Rollback the last applied migration",
	Long: bold("migrate:rollback") + ` rolls back the last applied migration
using Atlas migrate down.

` + colorGray + `Examples:` + colorReset + `
  grove migrate:rollback
  grove migrate:rollback --amount 3
  grove migrate:rollback --env dev`,
	RunE: runMigrateRollback,
}

func init() {
	migrateRollbackCmd.Flags().StringVar(
		&migrateRollbackEnv,
		"env", "local",
		"Atlas environment to use (local, dev, production)",
	)
	migrateRollbackCmd.Flags().IntVar(
		&migrateRollbackAmount,
		"amount", 1,
		"Number of migrations to roll back",
	)
}

func runMigrateRollback(cmd *cobra.Command, args []string) error {
	amount := fmt.Sprintf("%d", migrateRollbackAmount)

	fmt.Println()
	fmt.Printf(
		"  %sRolling back%s %s %s\n",
		colorGray, colorReset,
		bold(amount),
		gray("migration(s) (atlas migrate down --env "+migrateRollbackEnv+")"),
	)
	fmt.Println()

	return runAtlas("rollback migration",
		"migrate", "down",
		"--amount", amount,
		"--env", migrateRollbackEnv,
	)
}

// ──────────────────────────────────────────────
// migrate:status
// ──────────────────────────────────────────────

var migrateStatusEnv string

var migrateStatusCmd = &cobra.Command{
	Use:   "migrate:status",
	Short: "Show the status of all migrations",
	Long: bold("migrate:status") + ` shows which migrations have been applied
and which are still pending using Atlas migrate status.

` + colorGray + `Examples:` + colorReset + `
  grove migrate:status
  grove migrate:status --env dev`,
	RunE: runMigrateStatus,
}

func init() {
	migrateStatusCmd.Flags().StringVar(
		&migrateStatusEnv,
		"env", "local",
		"Atlas environment to use (local, dev, production)",
	)
}

func runMigrateStatus(cmd *cobra.Command, args []string) error {
	fmt.Println()
	fmt.Printf(
		"  %sChecking migration status%s %s\n",
		colorGray, colorReset,
		gray("(atlas migrate status --env "+migrateStatusEnv+")"),
	)
	fmt.Println()

	return runAtlas("check migration status",
		"migrate", "status",
		"--env", migrateStatusEnv,
	)
}

// ──────────────────────────────────────────────
// migrate:fresh
// ──────────────────────────────────────────────

var (
	migrateFreshEnv   string
	migrateFreshForce bool
)

var migrateFreshCmd = &cobra.Command{
	Use:   "migrate:fresh",
	Short: "Wipe the database and re-run all migrations",
	Long: bold("migrate:fresh") + ` drops all tables and re-applies every
migration from scratch. ` + colorRed + `This is a destructive operation.` + colorReset + `
Only use it on development databases.

` + colorGray + `Examples:` + colorReset + `
  grove migrate:fresh
  grove migrate:fresh --force
  grove migrate:fresh --env dev`,
	RunE: runMigrateFresh,
}

func init() {
	migrateFreshCmd.Flags().StringVar(
		&migrateFreshEnv,
		"env", "local",
		"Atlas environment to use (local, dev, production)",
	)
	migrateFreshCmd.Flags().BoolVar(
		&migrateFreshForce,
		"force", false,
		"Skip the confirmation prompt",
	)
}

func runMigrateFresh(cmd *cobra.Command, args []string) error {
	if !migrateFreshForce {
		fmt.Println()
		fmt.Printf(
			"  %s WARNING %s  This will %sdrop all tables%s and re-run every migration.\n",
			colorBgYellow,
			colorReset,
			colorRed,
			colorReset,
		)
		fmt.Printf(
			"  Environment: %s\n", bold(migrateFreshEnv),
		)
		fmt.Println()
		fmt.Print("  Are you sure? [y/N] ")

		var answer string
		if _, err := fmt.Scanln(
			&answer,
		); err != nil ||
			(answer != "y" && answer != "Y") {
			fmt.Println()
			fmt.Println(warn("Aborted."))
			fmt.Println()
			return nil
		}
	}

	fmt.Println()
	fmt.Printf(
		"  %sDropping all tables%s %s\n",
		colorGray, colorReset,
		gray("(atlas schema clean --env "+migrateFreshEnv+")"),
	)
	fmt.Println()

	if err := runAtlas("clean schema",
		"schema", "clean",
		"--env", migrateFreshEnv,
		"--auto-approve",
	); err != nil {
		return err
	}

	fmt.Println()
	fmt.Printf(
		"  %sRe-applying all migrations%s %s\n",
		colorGray, colorReset,
		gray("(atlas migrate apply --env "+migrateFreshEnv+")"),
	)
	fmt.Println()

	if err := runAtlas("apply all migrations",
		"migrate", "apply",
		"--env", migrateFreshEnv,
	); err != nil {
		return err
	}

	fmt.Println()
	fmt.Printf(
		"  %s DONE %s  Database wiped and migrations re-applied.\n",
		colorBgGreen, colorReset,
	)
	fmt.Println()

	return nil
}

// ──────────────────────────────────────────────
// migrate:hash
// ──────────────────────────────────────────────

var migrateHashCmd = &cobra.Command{
	Use:   "migrate:hash",
	Short: "Rehash the migrations directory (atlas migrate hash)",
	Long: bold("migrate:hash") + ` recalculates the atlas.sum file for the
migrations directory. Run this if you edited a migration file manually.

` + colorGray + `Examples:` + colorReset + `
  grove migrate:hash`,
	RunE: runMigrateHash,
}

func runMigrateHash(cmd *cobra.Command, args []string) error {
	fmt.Println()
	fmt.Printf(
		"  %sRehashing migrations directory%s %s\n",
		colorGray, colorReset,
		gray("(atlas migrate hash)"),
	)
	fmt.Println()

	return runAtlas("hash migrations directory", "migrate", "hash")
}

// ──────────────────────────────────────────────
// Shared atlas runner
// ──────────────────────────────────────────────

// runAtlas checks for the atlas binary, then runs it with the given arguments,
// forwarding stdout/stderr and stdin to the terminal.
func runAtlas(description string, atlasArgs ...string) error {
	if _, err := exec.LookPath("atlas"); err != nil {
		return fmt.Errorf(
			"atlas CLI not found in PATH\n\n  Install it from: %s",
			colorCyan+"https://atlasgo.io/docs"+colorReset,
		)
	}

	c := exec.Command("atlas", atlasArgs...)
	c.Stdout = newIndentWriter(os.Stdout, "  ")
	c.Stderr = newIndentWriter(os.Stderr, "  ")
	c.Stdin = os.Stdin

	if err := c.Run(); err != nil {
		return fmt.Errorf("failed to %s: %w", description, err)
	}

	return nil
}
