package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var makeResourceCmd = &cobra.Command{
	Use:   "make:resource <Name>",
	Short: "Scaffold model + migration + controller + DTO at once",
	Long: bold(
		"make:resource",
	) + ` scaffolds a model, migration, controller and DTO file in one shot.

This is equivalent to running ` + colorCyan + `grove make:model <Name> -r` + colorReset + `.
Every file respects the ` + colorCyan + `SKIPPED` + colorReset + ` rule — existing files are never overwritten.

The entity name is ` + colorBold + `automatically singularized` + colorReset + ` before generating files,
so you can pass the name in any form:

  ` + colorGray + `Post` + colorReset + `        → model ` + colorCyan + `Post` + colorReset + `, table ` + colorCyan + `posts` + colorReset + `, migration ` + colorCyan + `create_posts_table` + colorReset + `
  ` + colorGray + `Posts` + colorReset + `       → model ` + colorCyan + `Post` + colorReset + `, table ` + colorCyan + `posts` + colorReset + `, migration ` + colorCyan + `create_posts_table` + colorReset + `
  ` + colorGray + `order_items` + colorReset + ` → model ` + colorCyan + `OrderItem` + colorReset + `, table ` + colorCyan + `order_items` + colorReset + `, migration ` + colorCyan + `create_order_items_table` + colorReset + `

` + colorGray + `Examples:` + colorReset + `
  grove make:resource Post
  grove make:resource Posts
  grove make:resource BlogPost
  grove make:resource order_items`,
	Args: cobra.ExactArgs(1),
	RunE: runMakeResource,
}

func runMakeResource(_ *cobra.Command, args []string) error {
	name := toPascalCase(toSingular(args[0]))
	snake := toSnakeCase(name)
	tableName := toPlural(snake)
	migrationName := "create_" + strings.ToLower(tableName) + "_table"

	fmt.Println()
	fmt.Printf(
		"  %sScaffolding resource%s %s\n",
		colorGray, colorReset,
		bold(name),
	)
	fmt.Println()

	if err := scaffoldModel(name); err != nil {
		return err
	}

	if err := scaffoldController(name); err != nil {
		return err
	}

	if err := scaffoldRequest(name); err != nil {
		return err
	}

	// ── migration ────────────────────────────────────────────────────────────
	fmt.Println()
	fmt.Printf(
		"  %sGenerating migration%s %s %s\n",
		colorGray, colorReset,
		bold(migrationName),
		gray("(atlas migrate diff --env local)"),
	)
	fmt.Println()

	if _, err := exec.LookPath("atlas"); err != nil {
		fmt.Println(
			warn("atlas CLI not found — skipping migration generation"),
		)
		fmt.Printf(
			"  %sInstall it from: %s\n",
			colorGray,
			colorCyan+"https://atlasgo.io/docs"+colorReset,
		)
	} else {
		c := exec.Command(
			"atlas",
			"migrate",
			"diff",
			migrationName,
			"--env",
			"local",
		)
		c.Stdout = newIndentWriter(os.Stdout, "  ")
		c.Stderr = newIndentWriter(os.Stderr, "  ")
		c.Stdin = os.Stdin

		if err := c.Run(); err != nil {
			return fmt.Errorf("atlas migrate diff failed: %w", err)
		}

		fmt.Println()
		fmt.Println(done(
			"Migration " + bold(
				migrationName,
			) + " created in " + colorCyan + "migrations/" + colorReset,
		))
	}

	// ── next steps ───────────────────────────────────────────────────────────
	fmt.Println()
	fmt.Println(nextSteps())

	fmt.Printf(
		"    %s1.%s Add fields to the model in %s\n",
		colorGray, colorReset,
		colorCyan+"internal/models/"+snake+".go"+colorReset,
	)
	fmt.Printf(
		"    %s2.%s Fill in request/response fields in %s\n",
		colorGray, colorReset,
		colorCyan+"internal/dto/"+toKebabCase(name)+"-dto.go"+colorReset,
	)
	fmt.Printf(
		"    %s3.%s Run %s to apply the migration\n",
		colorGray, colorReset,
		colorGreen+"grove migrate"+colorReset,
	)
	fmt.Printf(
		"    %s4.%s Register routes in %s\n",
		colorGray, colorReset,
		colorCyan+"internal/routes/"+colorReset,
	)
	fmt.Println()

	return nil
}
