package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "grove",
	Short:         "An opinionated Go foundation for production-ready applications",
	Long:          buildBanner(),
	SilenceUsage:  true,
	SilenceErrors: true,
}

func buildBanner() string {
	sep := "  " + colorDim + strings.Repeat("─", 54) + colorReset

	logo := "\n" +
		logoG1 + `  █▀▀ █▀█ █▀█ █░█ █▀▀  ` + colorReset + "\n" +
		logoG3 + `  █▄█ █▀▄ █▄█ ▀▄▀ ██▄  ` + colorReset + "\n"

	tagline := "\n  " + colorDim + "An opinionated Go foundation for production-ready applications" + colorReset + "\n"

	generators := "\n" +
		"  " + colorBold + colorGray + "GENERATORS" + colorReset + "\n" +
		"    grove " + colorGreen + "make:model" + colorReset + "       <Name>   Scaffold a GORM model\n" +
		"    grove " + colorGreen + "make:controller" + colorReset + "  <Name>   Scaffold a fuego controller\n" +
		"    grove " + colorGreen + "make:dto" + colorReset + "         <Name>   Scaffold a DTO request/response file\n" +
		"    grove " + colorGreen + "make:middleware" + colorReset + "  <Name>   Scaffold an HTTP middleware\n" +
		"    grove " + colorGreen + "make:migration" + colorReset + "   <name>   Generate a migration via atlas migrate diff\n" +
		"    grove " + colorGreen + "make:resource" + colorReset + "    <Name>   Scaffold model + controller + DTO at once\n"

	server := "\n" +
		"  " + colorBold + colorGray + "SERVER" + colorReset + "\n" +
		"    grove " + colorBlue + "serve" + colorReset + "             Start the development server\n" +
		"    grove " + colorBlue + "build" + colorReset + "             Compile the application to a binary\n"

	database := "\n" +
		"  " + colorBold + colorGray + "DATABASE" + colorReset + "\n" +
		"    grove " + colorBlue + "migrate" + colorReset + "                 Apply all pending migrations\n" +
		"    grove " + colorBlue + "migrate:rollback" + colorReset + "        Rollback the last migration\n" +
		"    grove " + colorBlue + "migrate:status" + colorReset + "          Show migration status\n" +
		"    grove " + colorBlue + "migrate:fresh" + colorReset + "           Drop + re-apply all migrations\n" +
		"    grove " + colorBlue + "migrate:hash" + colorReset + "            Rehash the migrations directory\n"

	testing := "\n" +
		"  " + colorBold + colorGray + "TESTING" + colorReset + "\n" +
		"    grove " + colorGreen + "make:test" + colorReset + "        <Name>   Scaffold a new gest spec file\n" +
		"    grove " + colorBlue + "test" + colorReset + "              Run all gest specs in internal/tests\n" +
		"    grove " + colorBlue + "test -c" + colorReset + "           Run specs + display coverage report\n"

	shell := "\n" +
		"  " + colorBold + colorGray + "SHELL" + colorReset + "\n" +
		"    grove " + colorGray + "completion" + colorReset + "  [bash|zsh|fish|powershell]   Generate completion script\n"

	return logo + tagline + "\n" + sep + generators + server + database + testing + shell + "\n" + sep + "\n"
}

func init() {
	// completion
	rootCmd.AddCommand(completionCmd)

	// setup
	rootCmd.AddCommand(setupCmd)

	// make:* commands
	rootCmd.AddCommand(makeModelCmd)
	rootCmd.AddCommand(makeControllerCmd)
	rootCmd.AddCommand(makeDtoCmd)
	rootCmd.AddCommand(makeRequestCmd)
	rootCmd.AddCommand(makeMiddlewareCmd)
	rootCmd.AddCommand(makeMigrationCmd)
	rootCmd.AddCommand(makeResourceCmd)
	rootCmd.AddCommand(makeTestCmd)

	// server
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(buildCmd)

	// migrate commands
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(migrateRollbackCmd)
	rootCmd.AddCommand(migrateStatusCmd)
	rootCmd.AddCommand(migrateFreshCmd)
	rootCmd.AddCommand(migrateHashCmd)

	// testing
	rootCmd.AddCommand(testCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, fail(err.Error()))
		os.Exit(1)
	}
}
