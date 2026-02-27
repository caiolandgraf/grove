package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// ──────────────────────────────────────────────
// Model
// ──────────────────────────────────────────────

func scaffoldModel(name string) error {
	snake := toSnakeCase(name)
	tableName := toPlural(snake)
	destPath := filepath.Join("internal", "models", snake+".go")

	if fileExists(destPath) {
		printSkipped("Model", name, destPath)
		return nil
	}

	module := getModuleName()

	data := struct {
		Name      string
		TableName string
		Module    string
	}{
		Name:      name,
		TableName: tableName,
		Module:    module,
	}

	content, err := renderStub(modelStub, "model", data)
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
// Controller
// ──────────────────────────────────────────────

func scaffoldController(name string) error {
	snake := toSnakeCase(name)
	kebab := toKebabCase(name)
	destPath := filepath.Join("internal", "controllers", kebab+"-controller.go")

	if fileExists(destPath) {
		printSkipped("Controller", name, destPath)
		return nil
	}

	module := getModuleName()

	data := struct {
		Name      string
		ParamName string
		Module    string
	}{
		Name:      name,
		ParamName: snake,
		Module:    module,
	}

	content, err := renderStub(controllerStub, "controller", data)
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
// Request / DTO
// ──────────────────────────────────────────────

func scaffoldRequest(name string) error {
	kebab := toKebabCase(name)
	snake := toSnakeCase(name)
	destPath := filepath.Join("internal", "dto", kebab+"-dto.go")

	if fileExists(destPath) {
		printSkipped("DTO", name, destPath)
		return nil
	}

	data := struct {
		Name      string
		SnakeName string
	}{
		Name:      name,
		SnakeName: snake,
	}

	content, err := renderStub(requestStub, "request", data)
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
	destPath := filepath.Join("internal", "middleware", kebab+"-middleware.go")

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
// Test main entrypoint
// ──────────────────────────────────────────────

// scaffoldTestMain creates internal/tests/main.go (the gest entrypoint) if it
// does not already exist. It is called automatically by make:test.
func scaffoldTestMain() error {
	destPath := filepath.Join("internal", "tests", "main.go")

	if fileExists(destPath) {
		return nil // silently skip — user already has a main
	}

	module := getModuleName()

	data := struct {
		Module string
	}{
		Module: module,
	}

	content, err := renderStub(testMainStub, "test_main", data)
	if err != nil {
		return err
	}

	if err := writeFile(destPath, content); err != nil {
		return err
	}

	printCreated("Test entrypoint", "main", destPath)
	return nil
}

// ──────────────────────────────────────────────
// Test spec
// ──────────────────────────────────────────────

func scaffoldTestSpec(name string) error {
	snake := toSnakeCase(name)
	destPath := filepath.Join("internal", "tests", snake+"_spec.go")

	if fileExists(destPath) {
		printSkipped("Spec", name, destPath)
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

	content, err := renderStub(testSpecStub, "test_spec", data)
	if err != nil {
		return err
	}

	if err := writeFile(destPath, content); err != nil {
		return err
	}

	printCreated("Spec", name, destPath)
	return nil
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
