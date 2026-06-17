package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	makeSeederRegister bool
	makeSeederPath     string
)

var makeSeederCmd = &cobra.Command{
	Use:   "make:seeder <Name>",
	Short: "Scaffold a new database seeder",
	Long: bold(
		"make:seeder",
	) + ` scaffolds a new database seeder in ` + colorCyan + `internal/app/database/seeders/` + colorReset + `.

Seeders follow a simple interface:
  - Name() string
  - Seed(db *gorm.DB) error

The generated file will be named ` + colorCyan + `<name>_seeder.go` + colorReset + ` and will contain:
  - type <Name>Seeder struct{}
  - func (<Name>Seeder) Name() string
  - func (<Name>Seeder) Seed(db *gorm.DB) error

` + colorGray + `Examples:` + colorReset + `
  grove make:seeder Users
  grove make:seeder BillCategories
  grove make:seeder order_items

` + colorYellow + `Tip:` + colorReset + ` use ` + colorGreen + `--register` + colorReset + ` to add the seeder to ` + colorCyan + `internal/app/database/seeders/seeder.go` + colorReset + ` automatically when possible.`,
	Args: cobra.ExactArgs(1),
	RunE: runMakeSeeder,
}

func init() {
	makeSeederCmd.Flags().BoolVar(
		&makeSeederRegister,
		"register",
		false,
		"Also register the new seeder in internal/app/database/seeders/seeder.go (if present)",
	)
	makeSeederCmd.Flags().StringVar(
		&makeSeederPath,
		"path",
		filepath.Join("internal", "app", "database", "seeders"),
		"Target seeders directory",
	)
}

func runMakeSeeder(_ *cobra.Command, args []string) error {
	// We follow model/controller conventions: accept plural/mixed-case input and normalize.
	name := toPascalCase(toSingular(args[0]))
	snake := toSnakeCase(name)
	destPath := filepath.Join(makeSeederPath, snake+"_seeder.go")

	fmt.Println()
	fmt.Printf(
		"  %sCreating seeder%s %s\n",
		colorGray, colorReset,
		bold(name),
	)
	fmt.Println()

	if fileExists(destPath) {
		printSkipped("Seeder", name, destPath)
		fmt.Println()
		return nil
	}

	module := getModuleName()

	data := struct {
		Name   string
		Module string
	}{
		Name:   name,
		Module: module,
	}

	content, err := renderStub(seederStub, "seeder", data)
	if err != nil {
		return err
	}

	if err := writeFile(destPath, content); err != nil {
		return err
	}

	printCreated("Seeder", name, destPath)

	seedersOrchestratorPath := filepath.Join(makeSeederPath, "seeder.go")
	registered := false
	if makeSeederRegister {
		ok, regErr := registerSeederInOrchestrator(
			seedersOrchestratorPath,
			name,
		)
		if regErr != nil {
			return regErr
		}
		registered = ok
		if registered {
			fmt.Println(
				success("Registered seeder in " + seedersOrchestratorPath),
			)
		} else {
			fmt.Println(
				warn(
					"Could not auto-register seeder (no seeders list found) — register it manually in " + seedersOrchestratorPath,
				),
			)
		}
	}

	fmt.Println()
	fmt.Println(nextSteps())
	fmt.Printf(
		"    %s1.%s Implement seeding logic in %s\n",
		colorGray, colorReset,
		colorCyan+destPath+colorReset,
	)

	if !registered {
		fmt.Printf(
			"    %s2.%s Register it in %s (seeders list)\n",
			colorGray, colorReset,
			colorCyan+seedersOrchestratorPath+colorReset,
		)
		fmt.Printf(
			"    %s3.%s Run %s to execute all seeders\n",
			colorGray, colorReset,
			colorGreen+"grove db:seed"+colorReset,
		)
	} else {
		fmt.Printf(
			"    %s2.%s Run %s to execute all seeders\n",
			colorGray, colorReset,
			colorGreen+"grove db:seed"+colorReset,
		)
	}
	fmt.Println()

	return nil
}

func registerSeederInOrchestrator(
	orchestratorPath string,
	name string,
) (bool, error) {
	b, err := os.ReadFile(orchestratorPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Orchestrator doesn't exist; nothing to register into.
			return false, nil
		}
		return false, fmt.Errorf("read %s: %w", orchestratorPath, err)
	}

	src := string(b)
	seederInstance := name + "Seeder{}"

	// Already registered
	if strings.Contains(src, seederInstance) {
		return true, nil
	}

	// Find the seeders slice declaration and insert before the closing brace.
	// We look for:
	//   seeders := []Seeder{
	// and then the first matching closing "}" for that composite literal.
	startNeedle := "seeders := []Seeder{"
	start := strings.Index(src, startNeedle)
	if start < 0 {
		return false, nil
	}

	openPos := start + strings.Index(src[start:], "{")
	if openPos < start {
		return false, nil
	}

	closePos, findErr := findMatchingBrace(src, openPos)
	if findErr != nil {
		return false, fmt.Errorf(
			"parse seeders list in %s: %w",
			orchestratorPath,
			findErr,
		)
	}

	before := src[:closePos]
	after := src[closePos:]

	// Ensure there's a trailing newline before we insert.
	if !strings.HasSuffix(before, "\n") {
		before += "\n"
	}

	indent := detectIndentForSeederList(src[openPos+1 : closePos])
	insertion := indent + seederInstance + ",\n"

	newSrc := before + insertion + after

	// Write back
	if err := os.WriteFile(
		orchestratorPath,
		[]byte(newSrc),
		0o644,
	); err != nil {
		return false, fmt.Errorf("write %s: %w", orchestratorPath, err)
	}

	// Best-effort gofmt: only if this is a .go file and gofmt exists in PATH,
	// but we avoid hard-failing if gofmt is unavailable.
	_ = gofmtFile(orchestratorPath)

	return true, nil
}

func detectIndentForSeederList(body string) string {
	// Default: tab, matching most Go code.
	indent := "\t"

	// If we can find an existing line inside the list, reuse its indentation.
	scanner := bufio.NewScanner(strings.NewReader(body))
	for scanner.Scan() {
		line := scanner.Text()
		trim := strings.TrimSpace(line)
		if trim == "" || strings.HasPrefix(trim, "//") {
			continue
		}
		i := 0
		for i < len(line) && (line[i] == ' ' || line[i] == '\t') {
			i++
		}
		if i > 0 {
			return line[:i]
		}
		break
	}

	return indent
}

func gofmtFile(path string) error {
	// We can't import or depend on go/format here because we want gofmt's standard behavior,
	// and this repo already shells out for other tasks elsewhere.
	c := execCommand("gofmt", "-w", path)
	return c.Run()
}

// execCommand is a tiny indirection to allow build systems to stub it if needed.
var execCommand = func(name string, args ...string) *exec.Cmd {
	return exec.Command(name, args...)
}
