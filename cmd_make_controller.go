package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var makeControllerCmd = &cobra.Command{
	Use:   "make:controller <Name>",
	Short: "Scaffold a new fuego controller",
	Long: bold(
		"make:controller",
	) + ` scaffolds a new fuego controller in ` + colorCyan + `internal/controllers/` + colorReset + `.

` + colorGray + `Examples:` + colorReset + `
  grove make:controller Post
  grove make:controller BlogPost
  grove make:controller user_profile`,
	Args: cobra.ExactArgs(1),
	RunE: runMakeController,
}

func runMakeController(_ *cobra.Command, args []string) error {
	name := toPascalCase(args[0])
	snake := toSnakeCase(name)

	fmt.Println()
	fmt.Printf(
		"  %sCreating controller%s %s\n",
		colorGray, colorReset,
		bold(name),
	)
	fmt.Println()

	if err := scaffoldController(name); err != nil {
		return err
	}

	fmt.Println()
	fmt.Println(nextSteps())
	fmt.Printf(
		"    %s1.%s Fill in the field mappings marked with %s in %s\n",
		colorGray,
		colorReset,
		colorYellow+"// TODO"+colorReset,
		colorCyan+"internal/controllers/"+toKebabCase(
			name,
		)+"-controller.go"+colorReset,
	)
	fmt.Printf(
		"    %s2.%s Register routes in %s:\n",
		colorGray, colorReset,
		colorCyan+"internal/routes/"+colorReset,
	)
	fmt.Printf(
		"             %sfuego.Get(s, \"/%ss/{%s_id}\", controllers.Get%s)\n",
		colorGray,
		snake,
		snake,
		name+colorReset,
	)
	fmt.Printf(
		"             %sfuego.Get(s, \"/%ss\", controllers.List%ss)\n",
		colorGray,
		snake,
		name+colorReset,
	)
	fmt.Printf(
		"             %sfuego.Post(s, \"/%ss\", controllers.Create%s)\n",
		colorGray,
		snake,
		name+colorReset,
	)
	fmt.Printf(
		"             %sfuego.Put(s, \"/%ss/{%s_id}\", controllers.Update%s)\n",
		colorGray,
		snake,
		snake,
		name+colorReset,
	)
	fmt.Printf(
		"             %sfuego.Delete(s, \"/%ss/{%s_id}\", controllers.Delete%s)\n",
		colorGray,
		snake,
		snake,
		name+colorReset,
	)
	fmt.Println()

	return nil
}
