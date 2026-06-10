package main

import (
	"fmt"
	"path/filepath"
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
	) + ` scaffolds a new GORM model in ` + colorCyan + `internal/modules/<domain>/model.go` + colorReset + `.

The entity name is ` + colorBold + `automatically singularized` + colorReset + ` before generating files,
so you can pass the name in any form:

  ` + colorGray + `Post` + colorReset + `        → model ` + colorCyan + `Post` + colorReset + `, table ` + colorCyan + `posts` + colorReset + `, package ` + colorCyan + `posts` + colorReset + `
  ` + colorGray + `Posts` + colorReset + `       → model ` + colorCyan + `Post` + colorReset + `, table ` + colorCyan + `posts` + colorReset + `
  ` + colorGray + `order_items` + colorReset + ` → model ` + colorCyan + `OrderItem` + colorReset + `, table ` + colorCyan + `order_items` + colorReset + `

Combine flags to scaffold additional layers in one shot:

  ` + colorGreen + `-c` + colorReset + `  also scaffold controller + OpenAPI docs
  ` + colorGreen + `-d` + colorReset + `  also scaffold DTO file
  ` + colorGreen + `-r` + colorReset + `  full resource — service, controller, DTO, docs + module registration

` + colorYellow + `Migration workflow:` + colorReset + `
  Migrations are NOT generated automatically. After adding fields to your model,
  run ` + colorGreen + `grove make:migration <name>` + colorReset + ` to generate the SQL diff via Atlas.

` + colorGray + `Examples:` + colorReset + `
  grove make:model Post
  grove make:model Post -c
  grove make:model Post -cd
  grove make:model Post -r`,
	Args: cobra.ExactArgs(1),
	RunE: runMakeModel,
}

func init() {
	makeModelCmd.Flags().BoolVarP(
		&makeModelWithController,
		"controller", "c", false,
		"Also scaffold a controller and docs",
	)
	makeModelCmd.Flags().BoolVarP(
		&makeModelWithDTO,
		"dto", "d", false,
		"Also scaffold a DTO file",
	)
	makeModelCmd.Flags().BoolVarP(
		&makeModelResource,
		"resource", "r", false,
		"Full resource — service, controller, DTO, docs + module registration",
	)
}

func runMakeModel(_ *cobra.Command, args []string) error {
	name := toPascalCase(toSingular(args[0]))
	tableName := toPlural(toSnakeCase(name))

	if makeModelResource {
		makeModelWithController = true
		makeModelWithDTO = true
	}

	fmt.Println()
	fmt.Printf("  %sCreating model%s %s\n", colorGray, colorReset, bold(name))
	fmt.Println()

	if err := scaffoldModel(name); err != nil {
		return err
	}

	if makeModelResource {
		if err := scaffoldService(name); err != nil {
			return err
		}
	}

	if makeModelWithController {
		if err := scaffoldController(name, false); err != nil {
			return err
		}
		if err := scaffoldDocs(name); err != nil {
			return err
		}
	}

	if makeModelWithDTO {
		if err := scaffoldRequest(name); err != nil {
			return err
		}
	}

	if makeModelResource {
		if err := wireModule(name); err != nil {
			return err
		}
	}

	modDir := moduleDir(name)

	fmt.Println()
	fmt.Println(nextSteps())

	step := 1

	fmt.Printf(
		"    %s%d.%s Add your fields to the model in %s\n",
		colorGray, step, colorReset,
		colorCyan+filepath.Join(modDir, "model.go")+colorReset,
	)
	step++

	if makeModelWithDTO {
		fmt.Printf(
			"    %s%d.%s Fill in request/response fields in %s\n",
			colorGray, step, colorReset,
			colorCyan+filepath.Join(modDir, "dto.go")+colorReset,
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

	if makeModelResource {
		fmt.Printf(
			"    %s%d.%s Verify module registered in %s\n",
			colorGray, step, colorReset,
			colorCyan+"internal/modules/register.go"+colorReset,
		)
	} else if makeModelWithController {
		fmt.Printf(
			"    %s%d.%s Register the module in %s\n",
			colorGray, step, colorReset,
			colorCyan+"internal/modules/register.go"+colorReset,
		)
	}

	fmt.Println()

	return nil
}
