package main

import (
	_ "embed"
	"fmt"
)

// ──────────────────────────────────────────────
// Embedded stub templates
// ──────────────────────────────────────────────

//go:embed stubs/model.stub
var modelStub string

//go:embed stubs/controller.stub
var controllerStub string

//go:embed stubs/request.stub
var requestStub string

//go:embed stubs/middleware.stub
var middlewareStub string

//go:embed stubs/test_main.stub
var testMainStub string

//go:embed stubs/test_spec.stub
var testSpecStub string

// ──────────────────────────────────────────────
// Scaffold output helpers
// ──────────────────────────────────────────────

// printCreated prints a green "CREATED" badge line, e.g.:
//
//	CREATED   Model User → internal/models/user.go
func printCreated(kind, name, path string) {
	fmt.Printf("  %s CREATED %s  %s %s %s\n",
		colorBgGreen,
		colorReset,
		colorGray+kind+colorReset,
		bold(name),
		gray("→ "+path),
	)
}

// printSkipped prints a yellow "SKIPPED" badge line when a file already exists.
func printSkipped(kind, name, path string) {
	fmt.Printf("  %s SKIPPED %s  %s %s %s\n",
		colorBgYellow,
		colorReset,
		colorGray+kind+colorReset,
		bold(name),
		gray("→ "+path+" (already exists)"),
	)
}
