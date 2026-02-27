package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var makeResourceCmd = &cobra.Command{
	Use:   "make:resource <Name>",
	Short: "Scaffold model + controller + DTO at once",
	Long: bold(
		"make:resource",
	) + ` scaffolds a model, controller and DTO file in one shot.

This is equivalent to running:
  grove make:model      <Name>
  grove make:controller <Name>
  grove make:dto        <Name>

` + colorGray + `Examples:` + colorReset + `
  grove make:resource Post
  grove make:resource BlogPost
  grove make:resource order_item`,
	Args: cobra.ExactArgs(1),
	RunE: runMakeResource,
}

func runMakeResource(_ *cobra.Command, args []string) error {
	name := toPascalCase(args[0])

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

	fmt.Println()
	fmt.Println(nextSteps())

	snake := toSnakeCase(name)
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
		"    %s3.%s Run %s to generate the migration\n",
		colorGray,
		colorReset,
		colorGreen+"grove make:migration create_"+toPlural(
			snake,
		)+"_table"+colorReset,
	)
	fmt.Printf(
		"    %s4.%s Register routes in %s\n",
		colorGray, colorReset,
		colorCyan+"internal/routes/"+colorReset,
	)
	fmt.Println()

	return nil
}
