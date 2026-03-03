package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var makeDtoCmd = &cobra.Command{
	Use:   "make:dto <Name>",
	Short: "Scaffold a new DTO request/response file",
	Long: bold(
		"make:dto",
	) + ` scaffolds a new DTO request/response file in ` + colorCyan + `internal/dto/` + colorReset + `.

The entity name is ` + colorBold + `automatically singularized` + colorReset + ` before generating files,
so ` + colorCyan + `Posts` + colorReset + ` and ` + colorCyan + `Post` + colorReset + ` both produce the same ` + colorCyan + `post-dto.go` + colorReset + `.

` + colorGray + `Examples:` + colorReset + `
  grove make:dto Post
  grove make:dto Posts           # same as Post (singularized)
  grove make:dto BlogPost
  grove make:dto order_items`,
	Args: cobra.ExactArgs(1),
	RunE: runMakeDto,
}

func runMakeDto(_ *cobra.Command, args []string) error {
	name := toPascalCase(toSingular(args[0]))

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
		colorCyan+"internal/dto/"+toKebabCase(name)+"-dto.go"+colorReset,
	)
	fmt.Printf(
		"    %s2.%s Add response fields to %s as needed\n",
		colorGray, colorReset,
		colorCyan+name+"Response"+colorReset,
	)
	fmt.Println()

	return nil
}

// ──────────────────────────────────────────────
// make:request — backward-compat alias
// ──────────────────────────────────────────────

// makeRequestCmd is a hidden alias kept for backward compatibility.
// Prefer make:dto for new usage.
var makeRequestCmd = &cobra.Command{
	Use:    "make:request <Name>",
	Short:  "Scaffold a new DTO request/response file",
	Long:   "Deprecated alias for " + colorGreen + "make:dto" + colorReset + ". Use that instead.",
	Hidden: true,
	Args:   cobra.ExactArgs(1),
	RunE:   runMakeDto,
}
