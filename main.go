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
		"    grove " + colorBlue + "dev" + colorReset + "               Hot reload — watch, build & restart on save\n" +
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
		"    grove " + colorBlue + "test -c" + colorReset + "           Run specs + display coverage report\n" +
		"    grove " + colorBlue + "test -w" + colorReset + "           Watch mode — re-run specs on every save\n" +
		"    grove " + colorBlue + "test -wc" + colorReset + "          Watch mode + coverage report\n"

	setup := "\n" +
		"  " + colorBold + colorGray + "SETUP" + colorReset + "\n" +
		"    grove " + colorGray + "setup" + colorReset + "       <project-name>   Scaffold a new Grove project from template\n" +
		"    grove " + colorGray + "completion" + colorReset + "  [bash|zsh|fish|powershell]   Generate completion script\n"

	return logo + tagline + "\n" + sep + generators + server + database + testing + setup + "\n" + sep + "\n"
}

func init() {
	// ── Command groups (organises the "Available Commands" cobra help block) ──
	rootCmd.AddGroup(
		&cobra.Group{ID: "generators", Title: "Generators:"},
		&cobra.Group{ID: "testing", Title: "Testing:"},
		&cobra.Group{ID: "server", Title: "Server & Build:"},
		&cobra.Group{ID: "database", Title: "Database:"},
		&cobra.Group{ID: "setup", Title: "Setup:"},
	)

	// ── Generators ────────────────────────────────────────────────────────────
	makeModelCmd.GroupID = "generators"
	makeControllerCmd.GroupID = "generators"
	makeDtoCmd.GroupID = "generators"
	makeMiddlewareCmd.GroupID = "generators"
	makeMigrationCmd.GroupID = "generators"
	makeResourceCmd.GroupID = "generators"
	makeTestCmd.GroupID = "generators"

	rootCmd.AddCommand(makeModelCmd)
	rootCmd.AddCommand(makeControllerCmd)
	rootCmd.AddCommand(makeDtoCmd)
	rootCmd.AddCommand(makeRequestCmd) // hidden backward-compat alias
	rootCmd.AddCommand(makeMiddlewareCmd)
	rootCmd.AddCommand(makeMigrationCmd)
	rootCmd.AddCommand(makeResourceCmd)
	rootCmd.AddCommand(makeTestCmd)

	// ── Testing ───────────────────────────────────────────────────────────────
	testCmd.GroupID = "testing"

	rootCmd.AddCommand(testCmd)

	// ── Server & Build ────────────────────────────────────────────────────────
	devCmd.GroupID = "server"
	serveCmd.GroupID = "server"
	buildCmd.GroupID = "server"

	rootCmd.AddCommand(devCmd)
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(buildCmd)

	// ── Database ──────────────────────────────────────────────────────────────
	migrateCmd.GroupID = "database"
	migrateRollbackCmd.GroupID = "database"
	migrateStatusCmd.GroupID = "database"
	migrateFreshCmd.GroupID = "database"
	migrateHashCmd.GroupID = "database"

	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(migrateRollbackCmd)
	rootCmd.AddCommand(migrateStatusCmd)
	rootCmd.AddCommand(migrateFreshCmd)
	rootCmd.AddCommand(migrateHashCmd)

	// ── Setup ─────────────────────────────────────────────────────────────────
	setupCmd.GroupID = "setup"
	completionCmd.GroupID = "setup"

	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(completionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, fail(err.Error()))
		os.Exit(1)
	}
}
