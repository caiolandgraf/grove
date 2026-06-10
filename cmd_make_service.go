package main

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

var makeServiceCmd = &cobra.Command{
	Use:   "make:service <Name>",
	Short: "Scaffold a new service",
	Long: bold(
		"make:service",
	) + ` scaffolds a service in ` + colorCyan + `internal/modules/<domain>/service.go` + colorReset + `.

` + colorGray + `Examples:` + colorReset + `
  grove make:service Post
  grove make:service BlogPost`,
	Args: cobra.ExactArgs(1),
	RunE: runMakeService,
}

func runMakeService(_ *cobra.Command, args []string) error {
	name := toPascalCase(toSingular(args[0]))
	modDir := moduleDir(name)

	fmt.Println()
	fmt.Printf(
		"  %sCreating service%s %s\n",
		colorGray, colorReset,
		bold(name),
	)
	fmt.Println()

	if err := scaffoldService(name); err != nil {
		return err
	}

	fmt.Println()
	fmt.Println(nextSteps())
	fmt.Printf(
		"    %s1.%s Implement your business logic in %s\n",
		colorGray,
		colorReset,
		colorCyan+filepath.Join(modDir, "service.go")+colorReset,
	)
	fmt.Println()

	return nil
}
