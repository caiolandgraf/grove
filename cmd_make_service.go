package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var makeServiceCmd = &cobra.Command{
	Use:   "make:service <Name>",
	Short: "Scaffold a new service",
	Long: bold(
		"make:service",
	) + ` scaffolds a new service in ` + colorCyan + `internal/services/` + colorReset + `.

The entity name is ` + colorBold + `automatically singularized` + colorReset + ` before generating files,
so ` + colorCyan + `Posts` + colorReset + ` and ` + colorCyan + `Post` + colorReset + ` both produce the same ` + colorCyan + `PostService` + colorReset + `.

` + colorGray + `Examples:` + colorReset + `
  grove make:service Post
  grove make:service Posts        # same as Post (singularized)
  grove make:service BlogPost
  grove make:service user_profile`,
	Args: cobra.ExactArgs(1),
	RunE: runMakeService,
}

func runMakeService(_ *cobra.Command, args []string) error {
	name := toPascalCase(toSingular(args[0]))
	snake := toSnakeCase(name)

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
		colorCyan+"internal/services/"+snake+"_service.go"+colorReset,
	)
	fmt.Println()

	return nil
}
