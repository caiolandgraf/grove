package main

import (
	"fmt"

	"github.com/caiolandgraf/grove/internal/watcher"
	"github.com/spf13/cobra"
)

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Start the development server with built-in hot reload",
	Long: bold(
		"dev",
	) + ` compiles and runs your application, then watches for file
changes and automatically recompiles + restarts the binary on every save.

No external tools required — hot reload is built right into Grove.

Configure behaviour via the ` + colorCyan + `[dev]` + colorReset + ` section in ` + colorCyan + `grove.toml` + colorReset + `:

  ` + colorGray + `[dev]` + colorReset + `
  ` + colorGray + `root        = "."` + colorReset + `
  ` + colorGray + `bin         = ".grove/tmp/app"` + colorReset + `
  ` + colorGray + `build_cmd   = "go build -o .grove/tmp/app ./cmd/api/"` + colorReset + `
  ` + colorGray + `watch_dirs  = ["."]` + colorReset + `
  ` + colorGray + `exclude     = [".grove", "vendor", "node_modules", "tests"]` + colorReset + `
  ` + colorGray + `extensions  = [".go"]` + colorReset + `
  ` + colorGray + `debounce_ms = 300` + colorReset + `

All fields are optional — sensible defaults are used when ` + colorCyan + `grove.toml` + colorReset + `
is absent or the ` + colorCyan + `[dev]` + colorReset + ` section is omitted.

` + colorGray + `Examples:` + colorReset + `
  grove dev`,
	RunE: runDev,
}

// DevCmd exposes the cobra command so it can be wired from main.go.
// It is also the entry-point called by tests or external tooling.
func DevCmd() *cobra.Command { return devCmd }

func runDev(_ *cobra.Command, _ []string) error {
	fmt.Println()

	cfg, err := watcher.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load grove.toml: %w", err)
	}

	return watcher.New(cfg).Start()
}
