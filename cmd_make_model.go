package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var (
	makeModelWithMigration  bool
	makeModelWithController bool
	makeModelWithDTO        bool
	makeModelResource       bool
)

var makeModelCmd = &cobra.Command{
	Use:   "make:model <Name>",
	Short: "Scaffold a new GORM model",
	Long: bold(
		"make:model",
	) + ` scaffolds a new GORM model in ` + colorCyan + `internal/models/` + colorReset + `.

Combine flags to scaffold additional layers in one shot:

  ` + colorGreen + `-m` + colorReset + `  also generate a migration via atlas migrate diff
  ` + colorGreen + `-c` + colorReset + `  also scaffold a fuego controller
  ` + colorGreen + `-d` + colorReset + `  also scaffold a DTO request/response file
  ` + colorGreen + `-r` + colorReset + `  full resource — shorthand for ` + colorGreen + `-m -c -d` + colorReset + ` combined

` + colorGray + `Examples:` + colorReset + `
  grove make:model Post
  grove make:model Post -m
  grove make:model Post -mc
  grove make:model Post -mcd
  grove make:model Post -r
  grove make:model BlogPost -c
  grove make:model BlogPost -d
  grove make:model user_profile --resource`,
	Args: cobra.ExactArgs(1),
	RunE: runMakeModel,
}

func init() {
	makeModelCmd.Flags().BoolVarP(
		&makeModelWithMigration,
		"migration", "m", false,
		"Also generate a migration via atlas migrate diff",
	)
	makeModelCmd.Flags().BoolVarP(
		&makeModelWithController,
		"controller", "c", false,
		"Also scaffold a controller",
	)
	makeModelCmd.Flags().BoolVarP(
		&makeModelWithDTO,
		"dto", "d", false,
		"Also scaffold a DTO request/response file",
	)
	makeModelCmd.Flags().BoolVarP(
		&makeModelResource,
		"resource", "r", false,
		"Full resource — shorthand for -m -c -d",
	)
}

func runMakeModel(_ *cobra.Command, args []string) error {
	name := toPascalCase(args[0])
	snake := toSnakeCase(name)
	tableName := toPlural(snake)

	// -r expands to -m -c -d
	if makeModelResource {
		makeModelWithMigration = true
		makeModelWithController = true
		makeModelWithDTO = true
	}

	fmt.Println()
	fmt.Printf("  %sCreating model%s %s\n", colorGray, colorReset, bold(name))
	fmt.Println()

	// ── model ────────────────────────────────────────────────────────────────
	if err := scaffoldModel(name); err != nil {
		return err
	}

	// ── controller ───────────────────────────────────────────────────────────
	if makeModelWithController {
		if err := scaffoldController(name); err != nil {
			return err
		}
	}

	// ── DTO ──────────────────────────────────────────────────────────────────
	if makeModelWithDTO {
		if err := scaffoldRequest(name); err != nil {
			return err
		}
	}

	// ── migration ────────────────────────────────────────────────────────────
	if makeModelWithMigration {
		migrationName := "create_" + strings.ToLower(tableName) + "_table"

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
	}

	// ── next steps ───────────────────────────────────────────────────────────
	fmt.Println()
	fmt.Println(nextSteps())

	step := 1

	fmt.Printf(
		"    %s%d.%s Add your fields to the struct in %s\n",
		colorGray, step, colorReset,
		colorCyan+"internal/models/"+snake+".go"+colorReset,
	)
	step++

	if makeModelWithDTO {
		fmt.Printf(
			"    %s%d.%s Fill in request/response fields in %s\n",
			colorGray, step, colorReset,
			colorCyan+"internal/dto/"+toKebabCase(name)+"-dto.go"+colorReset,
		)
		step++
	}

	if !makeModelWithMigration {
		fmt.Printf(
			"    %s%d.%s Run %s to generate the migration\n",
			colorGray, step, colorReset,
			colorGreen+"grove make:migration create_"+strings.ToLower(
				tableName,
			)+"_table"+colorReset,
		)
		step++
	}

	fmt.Printf(
		"    %s%d.%s Run %s to apply it\n",
		colorGray, step, colorReset,
		colorGreen+"grove migrate"+colorReset,
	)
	step++

	if makeModelWithController {
		fmt.Printf(
			"    %s%d.%s Register routes in %s\n",
			colorGray, step, colorReset,
			colorCyan+"internal/routes/"+colorReset,
		)
	}

	fmt.Println()

	return nil
}
