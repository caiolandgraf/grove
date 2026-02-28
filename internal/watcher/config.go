package watcher

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

// Config holds all settings for the dev watcher.
// Every field maps 1-to-1 with the [dev] section in grove.toml.
type Config struct {
	// Root is the working directory from which build commands are run.
	Root string `toml:"root"`

	// TmpDir is the directory used for the compiled binary and other artifacts.
	TmpDir string `toml:"tmp_dir"`

	// Bin is the path to the compiled binary that will be executed.
	Bin string `toml:"bin"`

	// BuildCmd is the shell command used to compile the project.
	BuildCmd string `toml:"build_cmd"`

	// WatchDirs is the list of directories to watch for file changes.
	WatchDirs []string `toml:"watch_dirs"`

	// Exclude is a list of directory/file names that must never be watched.
	Exclude []string `toml:"exclude"`

	// Extensions is the list of file extensions that trigger a rebuild.
	Extensions []string `toml:"extensions"`

	// DebounceMs is the debounce window in milliseconds. Burst saves within
	// this window are collapsed into a single rebuild.
	DebounceMs int `toml:"debounce_ms"`
}

// DefaultConfig returns a Config populated with sensible out-of-the-box
// values so that grove dev works with zero configuration.
func DefaultConfig() Config {
	return Config{
		Root:      ".",
		TmpDir:    ".grove/tmp",
		Bin:       ".grove/tmp/app",
		BuildCmd:  "go build -o .grove/tmp/app ./cmd/api/",
		WatchDirs: []string{"."},
		Exclude: []string{
			".grove",
			"vendor",
			"node_modules",
			".git",
			"tests",
		},
		Extensions: []string{".go"},
		DebounceMs: 50,
	}
}

// groveFile mirrors the top-level structure of grove.toml so that the TOML
// decoder can navigate directly to the [dev] table.
type groveFile struct {
	Dev devSection `toml:"dev"`
}

// devSection mirrors Config but with pointer fields so we can distinguish
// "field was set in grove.toml" from "field was left at the zero value".
type devSection struct {
	Root       string   `toml:"root"`
	TmpDir     string   `toml:"tmp_dir"`
	Bin        string   `toml:"bin"`
	BuildCmd   string   `toml:"build_cmd"`
	WatchDirs  []string `toml:"watch_dirs"`
	Exclude    []string `toml:"exclude"`
	Extensions []string `toml:"extensions"`
	DebounceMs int      `toml:"debounce_ms"`
}

// LoadConfig reads the [dev] section from grove.toml in the current working
// directory and merges its values on top of DefaultConfig.
//
// If grove.toml does not exist the defaults are returned without error, making
// grove dev work out of the box with zero configuration.
func LoadConfig() (Config, error) {
	cfg := DefaultConfig()

	raw, err := os.ReadFile("grove.toml")
	if err != nil {
		if os.IsNotExist(err) {
			// No grove.toml — use defaults silently.
			return cfg, nil
		}
		return cfg, fmt.Errorf("grove.toml: %w", err)
	}

	var file groveFile
	if _, err := toml.Decode(string(raw), &file); err != nil {
		return cfg, fmt.Errorf("grove.toml parse error: %w", err)
	}

	dev := file.Dev

	// Merge: only override a default when the grove.toml value is non-zero,
	// so a partial [dev] section still benefits from the remaining defaults.
	if dev.Root != "" {
		cfg.Root = dev.Root
	}
	if dev.TmpDir != "" {
		cfg.TmpDir = dev.TmpDir
	}
	if dev.Bin != "" {
		cfg.Bin = dev.Bin
	}
	if dev.BuildCmd != "" {
		cfg.BuildCmd = dev.BuildCmd
	}
	if len(dev.WatchDirs) > 0 {
		cfg.WatchDirs = dev.WatchDirs
	}
	if len(dev.Exclude) > 0 {
		cfg.Exclude = dev.Exclude
	}
	if len(dev.Extensions) > 0 {
		cfg.Extensions = dev.Extensions
	}
	if dev.DebounceMs > 0 {
		cfg.DebounceMs = dev.DebounceMs
	}

	return cfg, nil
}
