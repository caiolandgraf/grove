package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

// gestModule is the fully-qualified module path used by go get to install or
// update gest to the latest published version.
const gestModule = "github.com/caiolandgraf/gest/v2@latest"

// gestCLIModule is the module path for the gest CLI binary.
const gestCLIModule = "github.com/caiolandgraf/gest/v2/cmd/gest@latest"

// ensureGest runs "go get github.com/caiolandgraf/gest/v2@latest" in the current
// working directory, adding or upgrading gest in the project's go.mod.
// Output from the command is forwarded with grove's standard indentation.
func ensureGest() error {
	c := exec.Command("go", "get", gestModule)
	c.Stdout = newIndentWriter(os.Stdout, "  ")
	c.Stderr = newIndentWriter(os.Stderr, "  ")
	return c.Run()
}

// modulePackage returns the Go package name for a domain module (e.g. "posts").
func modulePackage(name string) string {
	return toPlural(toSnakeCase(name))
}

// moduleDir returns the filesystem path for a domain module.
func moduleDir(name string) string {
	return filepath.Join("internal", "modules", modulePackage(name))
}

type moduleScaffoldData struct {
	Name        string
	NameLower   string
	Module      string
	Package     string
	TableName   string
	PluralName  string
	RoutePrefix string
	PathParam   string
}

func newModuleScaffoldData(name string) moduleScaffoldData {
	snake := toSnakeCase(name)
	return moduleScaffoldData{
		Name:        name,
		NameLower:   toLowerFirst(name),
		Module:      getModuleName(),
		Package:     modulePackage(name),
		TableName:   toPlural(snake),
		PluralName:  toPascalCase(toPlural(name)),
		RoutePrefix: toPlural(snake),
		PathParam:   "{" + snake + "_id}",
	}
}

// ──────────────────────────────────────────────
// Model
// ──────────────────────────────────────────────

func scaffoldModel(name string) error {
	destPath := filepath.Join(moduleDir(name), "model.go")

	if fileExists(destPath) {
		printSkipped("Model", name, destPath)
		return nil
	}

	content, err := renderStub(modelStub, "model", newModuleScaffoldData(name))
	if err != nil {
		return err
	}

	if err := writeFile(destPath, content); err != nil {
		return err
	}

	printCreated("Model", name, destPath)
	return nil
}

// ──────────────────────────────────────────────
// Service
// ──────────────────────────────────────────────

func scaffoldService(name string) error {
	destPath := filepath.Join(moduleDir(name), "service.go")

	if fileExists(destPath) {
		printSkipped("Service", name, destPath)
		return nil
	}

	content, err := renderStub(serviceStub, "service", newModuleScaffoldData(name))
	if err != nil {
		return err
	}

	if err := writeFile(destPath, content); err != nil {
		return err
	}

	printCreated("Service", name, destPath)
	return nil
}

// ──────────────────────────────────────────────
// Controller
// ──────────────────────────────────────────────

func scaffoldController(name string, noAuth bool) error {
	destPath := filepath.Join(moduleDir(name), "controller.go")

	if fileExists(destPath) {
		printSkipped("Controller", name, destPath)
		return nil
	}

	stub := controllerStub
	if noAuth {
		stub = controllerLegacyStub
	}

	content, err := renderStub(stub, "controller", newModuleScaffoldData(name))
	if err != nil {
		return err
	}

	if err := writeFile(destPath, content); err != nil {
		return err
	}

	printCreated("Controller", name, destPath)
	return nil
}

// ──────────────────────────────────────────────
// Docs
// ──────────────────────────────────────────────

func scaffoldDocs(name string) error {
	destPath := filepath.Join(moduleDir(name), "docs.go")

	if fileExists(destPath) {
		printSkipped("Docs", name, destPath)
		return nil
	}

	content, err := renderStub(docsStub, "docs", newModuleScaffoldData(name))
	if err != nil {
		return err
	}

	if err := writeFile(destPath, content); err != nil {
		return err
	}

	printCreated("Docs", name, destPath)
	return nil
}

// ──────────────────────────────────────────────
// DTO
// ──────────────────────────────────────────────

func scaffoldRequest(name string) error {
	destPath := filepath.Join(moduleDir(name), "dto.go")

	if fileExists(destPath) {
		printSkipped("DTO", name, destPath)
		return nil
	}

	content, err := renderStub(dtoStub, "dto", newModuleScaffoldData(name))
	if err != nil {
		return err
	}

	if err := writeFile(destPath, content); err != nil {
		return err
	}

	printCreated("DTO", name, destPath)
	return nil
}

// ──────────────────────────────────────────────
// Middleware
// ──────────────────────────────────────────────

func scaffoldMiddleware(name string) error {
	kebab := toKebabCase(name)
	destPath := filepath.Join("internal", "app", "middleware", kebab+"-middleware.go")

	if fileExists(destPath) {
		printSkipped("Middleware", name, destPath)
		return nil
	}

	data := struct {
		Name string
	}{
		Name: name,
	}

	content, err := renderStub(middlewareStub, "middleware", data)
	if err != nil {
		return err
	}

	if err := writeFile(destPath, content); err != nil {
		return err
	}

	printCreated("Middleware", name, destPath)
	return nil
}

// ──────────────────────────────────────────────
// Test spec
// ──────────────────────────────────────────────

// scaffoldTestSpec creates internal/tests/<snake>_test.go for gest v2.
// On the first call it also runs "go get" to add gest to the project's go.mod.
func scaffoldTestSpec(name string) error {
	snake := toSnakeCase(name)
	destPath := filepath.Join("internal", "tests", snake+"_test.go")

	isFirstSpec := !dirHasTestFiles(filepath.Join("internal", "tests"))

	if fileExists(destPath) {
		printSkipped("Test", name, destPath)
		return nil
	}

	pkg := getPackageName()

	data := struct {
		Name    string
		Package string
		Label   string
	}{
		Name:    name,
		Package: pkg,
		Label:   toWords(name),
	}

	content, err := renderStub(testSpecStub, "test_spec", data)
	if err != nil {
		return err
	}

	if err := writeFile(destPath, content); err != nil {
		return err
	}

	printCreated("Test", name, destPath)

	if isFirstSpec {
		fmt.Println()
		fmt.Printf(
			"  %sInstalling gest%s %s\n",
			colorGray, colorReset,
			gray("(go get "+gestModule+")"),
		)
		fmt.Println()

		if err := ensureGest(); err != nil {
			fmt.Println(warn("Failed to install gest automatically."))
			fmt.Printf(
				"  %sRun manually: %s\n",
				colorGray,
				colorGreen+"go get "+gestModule+colorReset,
			)
		}
	}

	return nil
}

// dirHasTestFiles reports whether dir contains at least one *_test.go file.
func dirHasTestFiles(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, e := range entries {
		if !e.IsDir() && len(e.Name()) > 8 &&
			e.Name()[len(e.Name())-8:] == "_test.go" {
			return true
		}
	}
	return false
}

// ──────────────────────────────────────────────
// Low-level helpers
// ──────────────────────────────────────────────

// renderStub parses and executes a text/template stub with the given data.
func renderStub(stub, name string, data any) ([]byte, error) {
	tmpl, err := template.New(name).Parse(stub)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s stub: %w", name, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("failed to render %s stub: %w", name, err)
	}

	return buf.Bytes(), nil
}

// writeFile ensures the parent directory exists and writes content to path.
func writeFile(path string, content []byte) error {
	if err := ensureDir(filepath.Dir(path)); err != nil {
		return fmt.Errorf("failed to create directory for %s: %w", path, err)
	}

	if err := os.WriteFile(path, content, 0o644); err != nil {
		return fmt.Errorf("failed to write %s: %w", path, err)
	}

	return nil
}

// getPackageName returns the Go package name to use for generated test files.
func getPackageName() string {
	return "tests"
}

// wireModule registers a domain module in internal/modules/register.go.
func wireModule(name string) error {
	registerPath := filepath.Join("internal", "modules", "register.go")
	if !fileExists(registerPath) {
		fmt.Printf("  %s WARNING %s  register file not found at %s. Skipping auto-wiring.\n",
			colorBgYellow, colorReset, registerPath)
		return nil
	}

	contentBytes, err := os.ReadFile(registerPath)
	if err != nil {
		return fmt.Errorf("failed to read register file: %w", err)
	}

	content := string(contentBytes)
	pkg := modulePackage(name)
	moduleImport := fmt.Sprintf("%s/internal/modules/%s", getModuleName(), pkg)
	wireLine := fmt.Sprintf("\tfunc(b Boot) Module { return %s.Wire(b.DB) },", pkg)

	if strings.Contains(content, wireLine) || strings.Contains(content, pkg+".Wire(") {
		return nil
	}

	if !strings.Contains(content, moduleImport) {
		importNeedle := "\t\"github.com/go-fuego/fuego\""
		importLine := fmt.Sprintf("\t\"%s\"\n%s", moduleImport, importNeedle)
		if strings.Contains(content, importNeedle) {
			content = strings.Replace(content, importNeedle, importLine, 1)
		} else {
			return fmt.Errorf("could not find import block anchor in register.go")
		}
	}

	registryNeedle := "var registry = []Factory{"
	idx := strings.Index(content, registryNeedle)
	if idx == -1 {
		return fmt.Errorf("could not find registry in register.go")
	}

	insertPos := idx + len(registryNeedle)
	if insertPos < len(content) && content[insertPos] == '\n' {
		insertPos++
	}

	content = content[:insertPos] + wireLine + "\n" + content[insertPos:]

	if err := os.WriteFile(registerPath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("failed to write updated register.go: %w", err)
	}

	fmt.Printf("  %s WIRED   %s  %s %s %s\n",
		colorBgGreen,
		colorReset,
		colorGray+"Module"+colorReset,
		bold(name),
		gray("→ "+registerPath),
	)
	return nil
}
