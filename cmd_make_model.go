package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
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

The entity name is ` + colorBold + `automatically singularized` + colorReset + ` before generating files,
so you can pass the name in any form:

  ` + colorGray + `Post` + colorReset + `        → model ` + colorCyan + `Post` + colorReset + `, table ` + colorCyan + `posts` + colorReset + `
  ` + colorGray + `Posts` + colorReset + `       → model ` + colorCyan + `Post` + colorReset + `, table ` + colorCyan + `posts` + colorReset + `
  ` + colorGray + `order_items` + colorReset + ` → model ` + colorCyan + `OrderItem` + colorReset + `, table ` + colorCyan + `order_items` + colorReset + `

Combine flags to scaffold additional layers in one shot:

  ` + colorGreen + `-c` + colorReset + `  also scaffold a fuego controller
  ` + colorGreen + `-d` + colorReset + `  also scaffold a DTO request/response file
  ` + colorGreen + `-r` + colorReset + `  full resource — shorthand for ` + colorGreen + `-c -d` + colorReset + ` combined

` + colorYellow + `Migration workflow:` + colorReset + `
  Migrations are NOT generated automatically. After adding fields to your model,
  run ` + colorGreen + `grove make:migration <name>` + colorReset + ` to generate the SQL diff via Atlas.

  This ensures your migration reflects the actual fields you defined — not an
  empty struct scaffolded before you had a chance to edit it.

` + colorGray + `Examples:` + colorReset + `
  grove make:model Post
  grove make:model Posts        # same as Post (singularized)
  grove make:model Post -c
  grove make:model Post -cd
  grove make:model Post -r
  grove make:model BlogPost -c
  grove make:model BlogPost -d
  grove make:model order_items --resource`,
	Args: cobra.ExactArgs(1),
	RunE: runMakeModel,
}

func init() {
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
		"Full resource — shorthand for -c -d",
	)
}

func runMakeModel(_ *cobra.Command, args []string) error {
	name := toPascalCase(toSingular(args[0]))
	snake := toSnakeCase(name)
	tableName := toPlural(snake)

	// -r expands to -c -d
	if makeModelResource {
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

	// ── next steps ───────────────────────────────────────────────────────────
	fmt.Println()
	fmt.Println(nextSteps())

	step := 1

	fmt.Printf(
		"    %s%d.%s Add your fields to the model in %s\n",
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

	fmt.Printf(
		"    %s%d.%s Run %s to generate the migration\n",
		colorGray,
		step,
		colorReset,
		colorGreen+"grove make:migration create_"+strings.ToLower(
			tableName,
		)+"_table"+colorReset,
	)
	step++

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
