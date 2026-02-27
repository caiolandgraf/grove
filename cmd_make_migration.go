package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var makeMigrationEnv string

var makeMigrationCmd = &cobra.Command{
	Use:   "make:migration <name>",
	Short: "Generate a new migration via atlas migrate diff",
	Long: bold(
		"make:migration",
	) + ` generates a new SQL migration file by diffing your GORM
models against the current database schema using Atlas.

Make sure you have added or updated your models in ` + colorCyan + `internal/models/` + colorReset + `
before running this command.

` + colorGray + `Examples:` + colorReset + `
  grove make:migration add_posts_table
  grove make:migration add_email_to_users
  grove make:migration create_orders_table --env dev`,
	Args: cobra.ExactArgs(1),
	RunE: runMakeMigration,
}

func init() {
	makeMigrationCmd.Flags().StringVar(
		&makeMigrationEnv,
		"env", "local",
		"Atlas environment to use (local, dev, production)",
	)
}

func runMakeMigration(cmd *cobra.Command, args []string) error {
	name := args[0]
	// Normalize: spaces → underscores, lowercase
	name = strings.ToLower(strings.ReplaceAll(name, " ", "_"))
	name = strings.ReplaceAll(name, "-", "_")

	fmt.Println()
	fmt.Printf(
		"  %sGenerating migration%s %s %s\n",
		colorGray, colorReset,
		bold(name),
		gray("(atlas migrate diff --env "+makeMigrationEnv+")"),
	)
	fmt.Println()

	// Check atlas is available
	if _, err := exec.LookPath("atlas"); err != nil {
		return fmt.Errorf(
			"atlas CLI not found in PATH\n\n  Install it from: %s",
			colorCyan+"https://atlasgo.io/docs"+colorReset,
		)
	}

	atlasArgs := []string{
		"migrate", "diff", name,
		"--env", makeMigrationEnv,
	}

	c := exec.Command("atlas", atlasArgs...)
	c.Stdout = newIndentWriter(os.Stdout, "  ")
	c.Stderr = newIndentWriter(os.Stderr, "  ")
	c.Stdin = os.Stdin

	if err := c.Run(); err != nil {
		return fmt.Errorf("atlas migrate diff failed: %w", err)
	}

	fmt.Println()
	fmt.Println(done(
		"Migration " + bold(
			name,
		) + " created in " + colorCyan + "migrations/" + colorReset,
	))
	fmt.Println()
	fmt.Println(nextSteps())
	fmt.Printf(
		"    %s1.%s Review the generated SQL in %s\n",
		colorGray, colorReset,
		colorCyan+"migrations/"+colorReset,
	)
	fmt.Printf(
		"    %s2.%s Run %s to apply it\n",
		colorGray, colorReset,
		colorGreen+"grove migrate"+colorReset,
	)
	fmt.Println()

	return nil
}
