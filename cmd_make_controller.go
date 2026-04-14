package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var makeControllerNoAuth bool

var makeControllerCmd = &cobra.Command{
	Use:   "make:controller <Name>",
	Short: "Scaffold a new fuego controller",
	Long: bold(
		"make:controller",
	) + ` scaffolds a new fuego controller in ` + colorCyan + `internal/controllers/` + colorReset + `.

By default, the generated controller uses a struct + constructor header (ready for session/auth wiring).
Pass ` + colorGreen + `--no-auth` + colorReset + ` to generate the legacy function-based controller stub instead.

The entity name is ` + colorBold + `automatically singularized` + colorReset + ` before generating files,
so ` + colorCyan + `Posts` + colorReset + ` and ` + colorCyan + `Post` + colorReset + ` both produce the same ` + colorCyan + `PostController` + colorReset + `.

` + colorGray + `Examples:` + colorReset + `
  grove make:controller Post
  grove make:controller Posts        # same as Post (singularized)
  grove make:controller BlogPost
  grove make:controller user_profile
  grove make:controller Bill         # new header (default)
  grove make:controller Bill --no-auth   # legacy stub`,
	Args: cobra.ExactArgs(1),
	RunE: runMakeController,
}

func init() {
	makeControllerCmd.Flags().BoolVar(
		&makeControllerNoAuth,
		"no-auth",
		false,
		"Generate the legacy function-based controller stub (no session/auth header)",
	)
}

func runMakeController(_ *cobra.Command, args []string) error {
	name := toPascalCase(toSingular(args[0]))
	snake := toSnakeCase(name)

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

	if makeControllerNoAuth {
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

	fmt.Printf(
		"             %s%s\n",
		colorGray,
		"api := fuego.Group(s, \"/api/v1\")"+colorReset,
	)
	fmt.Printf(
		"             %s%s\n",
		colorGray,
		"group := fuego.Group(api, \"/"+toPlural(snake)+"\")"+colorReset,
	)
	fmt.Printf(
		"             %s%s\n",
		colorGray,
		toLowerFirst(
			name,
		)+"Controller := controllers.New"+name+"Controller(app.Session)"+colorReset,
	)
	fmt.Printf(
		"             %sfuego.Get(group, \"/\", %sController.List)\n",
		colorGray,
		toSnakeCase(toLowerFirst(name)),
	)
	fmt.Printf(
		"             %sfuego.Post(group, \"/\", %sController.Create)\n",
		colorGray,
		toSnakeCase(toLowerFirst(name)),
	)
	fmt.Printf(
		"             %sfuego.Get(group, \"/{%s_id}\", %sController.Get)\n",
		colorGray,
		snake,
		toSnakeCase(toLowerFirst(name)),
	)
	fmt.Printf(
		"             %sfuego.Put(group, \"/{%s_id}\", %sController.Update)\n",
		colorGray,
		snake,
		toSnakeCase(toLowerFirst(name)),
	)
	fmt.Printf(
		"             %sfuego.Delete(group, \"/{%s_id}\", %sController.Delete)\n",
		colorGray,
		snake,
		toSnakeCase(toLowerFirst(name)),
	)
	fmt.Println()

	return nil
}
