package main

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

var makeDtoCmd = &cobra.Command{
	Use:   "make:dto <Name>",
	Short: "Scaffold a new DTO request/response file",
	Long: bold(
		"make:dto",
	) + ` scaffolds DTO types in ` + colorCyan + `internal/modules/<domain>/dto.go` + colorReset + `.

` + colorGray + `Examples:` + colorReset + `
  grove make:dto Post
  grove make:dto BlogPost`,
	Args: cobra.ExactArgs(1),
	RunE: runMakeDto,
}

func runMakeDto(_ *cobra.Command, args []string) error {
	name := toPascalCase(toSingular(args[0]))
	modDir := moduleDir(name)

	fmt.Println()
	fmt.Printf(
		"  %sCreating DTO%s %s\n",
		colorGray, colorReset,
		bold(name),
	)
	fmt.Println()

	if err := scaffoldRequest(name); err != nil {
		return err
	}

	fmt.Println()
	fmt.Println(nextSteps())
	fmt.Printf(
		"    %s1.%s Add your fields to %s and %s in %s\n",
		colorGray, colorReset,
		colorCyan+"Create"+name+"Request"+colorReset,
		colorCyan+"Update"+name+"Request"+colorReset,
		colorCyan+filepath.Join(modDir, "dto.go")+colorReset,
	)
	fmt.Println()

	return nil
}

var makeRequestCmd = &cobra.Command{
	Use:    "make:request <Name>",
	Short:  "Scaffold a new DTO request/response file",
	Long:   "Deprecated alias for " + colorGreen + "make:dto" + colorReset + ". Use that instead.",
	Hidden: true,
	Args:   cobra.ExactArgs(1),
	RunE:   runMakeDto,
}
