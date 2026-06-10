package main

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

var makeControllerNoAuth bool

var makeControllerCmd = &cobra.Command{
	Use:   "make:controller <Name>",
	Short: "Scaffold a new fuego controller",
	Long: bold(
		"make:controller",
	) + ` scaffolds a module controller in ` + colorCyan + `internal/modules/<domain>/controller.go` + colorReset + `.

The generated controller includes ` + colorCyan + `Mount` + colorReset + ` and ` + colorCyan + `Wire` + colorReset + ` methods
matching the modular MSC architecture. OpenAPI docs are scaffolded in ` + colorCyan + `docs.go` + colorReset + `.

Pass ` + colorGreen + `--no-auth` + colorReset + ` to generate legacy function-based handlers instead.

` + colorGray + `Examples:` + colorReset + `
  grove make:controller Post
  grove make:controller Post --no-auth`,
	Args: cobra.ExactArgs(1),
	RunE: runMakeController,
}

func init() {
	makeControllerCmd.Flags().BoolVar(
		&makeControllerNoAuth,
		"no-auth",
		false,
		"Generate legacy function-based handlers instead of a struct controller",
	)
}

func runMakeController(_ *cobra.Command, args []string) error {
	name := toPascalCase(toSingular(args[0]))
	modDir := moduleDir(name)

	fmt.Println()
	fmt.Printf(
		"  %sCreating controller%s %s\n",
		colorGray, colorReset,
		bold(name),
	)
	fmt.Println()

	if err := scaffoldController(name, makeControllerNoAuth); err != nil {
		return err
	}

	if !makeControllerNoAuth {
		if err := scaffoldDocs(name); err != nil {
			return err
		}
	}

	fmt.Println()
	fmt.Println(nextSteps())
	fmt.Printf(
		"    %s1.%s Implement handlers in %s\n",
		colorGray, colorReset,
		colorCyan+filepath.Join(modDir, "controller.go")+colorReset,
	)
	if !makeControllerNoAuth {
		fmt.Printf(
			"    %s2.%s Register the module in %s\n",
			colorGray, colorReset,
			colorCyan+"internal/modules/register.go"+colorReset,
		)
	}
	fmt.Println()

	return nil
}
