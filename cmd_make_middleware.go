package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var makeMiddlewareCmd = &cobra.Command{
	Use:   "make:middleware <Name>",
	Short: "Scaffold a new HTTP middleware",
	Long: bold(
		"make:middleware",
	) + ` scaffolds a new HTTP middleware in ` + colorCyan + `internal/middleware/` + colorReset + `.

` + colorGray + `Examples:` + colorReset + `
  grove make:middleware Auth
  grove make:middleware RateLimit
  grove make:middleware cors_headers`,
	Args: cobra.ExactArgs(1),
	RunE: runMakeMiddleware,
}

func runMakeMiddleware(_ *cobra.Command, args []string) error {
	name := toPascalCase(args[0])
	kebab := toKebabCase(name)

	fmt.Println()
	fmt.Printf(
		"  %sCreating middleware%s %s\n",
		colorGray, colorReset,
		bold(name),
	)
	fmt.Println()

	if err := scaffoldMiddleware(name); err != nil {
		return err
	}

	fmt.Println()
	fmt.Println(nextSteps())
	fmt.Printf(
		"    %s1.%s Implement your logic in %s\n",
		colorGray, colorReset,
		colorCyan+"internal/middleware/"+kebab+"-middleware.go"+colorReset,
	)
	fmt.Printf(
		"    %s2.%s Register it in your server setup, e.g.:\n",
		colorGray, colorReset,
	)
	fmt.Printf(
		"             %ss.Use(middleware.%s)\n",
		colorGray,
		name+colorReset,
	)
	fmt.Println()

	return nil
}
