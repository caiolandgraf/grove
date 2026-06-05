package watcher

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

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

	// PortGuard optionally frees the port before starting the app.
	// Set to 0 to disable.
	PortGuard int `toml:"port_guard"`
}

// DefaultConfig returns a Config populated with sensible out-of-the-box
// values so that grove dev works with zero configuration.
func DefaultConfig() Config {
	watchDirs := defaultWatchDirs()
	return Config{
		Root:      ".",
		TmpDir:    ".grove/tmp",
		Bin:       ".grove/tmp/app",
		BuildCmd:  "go build -o .grove/tmp/app ./cmd/api/",
		WatchDirs: watchDirs,
		Exclude: []string{
			".grove",
			"vendor",
			"node_modules",
			".git",
			"infra",
			"migrations",
			"bin",
			"internal/tests",
		},
		Extensions: []string{".go"},
		DebounceMs: 50,
		PortGuard:  0,
	}
}

func defaultWatchDirs() []string {
	var watch []string
	if dirExists("cmd") {
		watch = append(watch, "cmd")
	}
	if dirExists("internal") {
		watch = append(watch, "internal")
	}
	if len(watch) == 0 {
		return []string{"."}
	}
	return watch
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
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
	PortGuard  *int     `toml:"port_guard"`
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
	portGuardSet := dev.PortGuard != nil
	if portGuardSet {
		cfg.PortGuard = *dev.PortGuard
	}

	if cfg.PortGuard == 0 && !portGuardSet {
		if inferred := inferPortGuard(cfg.Root); inferred > 0 {
			cfg.PortGuard = inferred
		}
	}

	return cfg, nil
}

func inferPortGuard(root string) int {
	ports := readEnvPorts(root)
	for _, key := range []string{
		"APP_PORT",
		"PORT",
		"HTTP_PORT",
		"SERVER_PORT",
	} {
		if port, ok := ports[key]; ok && port > 0 {
			return port
		}
	}

	if baseURL, ok := readEnvString(root, "BASE_URL"); ok {
		if port := parsePortFromURL(baseURL); port > 0 {
			return port
		}
	}

	if addr, ok := readEnvString(root, "APP_ADDR"); ok {
		if port := parsePortFromAddr(addr); port > 0 {
			return port
		}
	}

	return 0
}

func readEnvPorts(root string) map[string]int {
	ports := make(map[string]int)
	if envs := readEnvFile(root); len(envs) > 0 {
		for k, v := range envs {
			if port, err := strconv.Atoi(strings.TrimSpace(v)); err == nil {
				ports[k] = port
			}
		}
	}

	for _, key := range []string{"APP_PORT", "PORT", "HTTP_PORT", "SERVER_PORT"} {
		if v, ok := os.LookupEnv(key); ok {
			if port, err := strconv.Atoi(strings.TrimSpace(v)); err == nil {
				ports[key] = port
			}
		}
	}

	return ports
}

func readEnvString(root, key string) (string, bool) {
	if envs := readEnvFile(root); len(envs) > 0 {
		if v, ok := envs[key]; ok {
			return v, true
		}
	}
	if v, ok := os.LookupEnv(key); ok {
		return v, true
	}
	return "", false
}

func readEnvFile(root string) map[string]string {
	path := filepath.Join(root, ".env")
	content, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	lines := strings.Split(string(content), "\n")
	out := make(map[string]string)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		}
		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		val = strings.TrimSpace(val)
		if len(val) >= 2 {
			if (val[0] == '"' && val[len(val)-1] == '"') ||
				(val[0] == '\'' && val[len(val)-1] == '\'') {
				val = val[1 : len(val)-1]
			}
		}
		out[key] = val
	}

	return out
}

func parsePortFromURL(raw string) int {
	parsed, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return 0
	}
	if portStr := parsed.Port(); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			return port
		}
	}
	return 0
}

func parsePortFromAddr(raw string) int {
	addr := strings.TrimSpace(raw)
	addr = strings.TrimPrefix(addr, "http://")
	addr = strings.TrimPrefix(addr, "https://")
	if strings.HasPrefix(addr, ":") {
		if port, err := strconv.Atoi(
			strings.TrimPrefix(addr, ":"),
		); err == nil {
			return port
		}
	}
	if host, portStr, err := net.SplitHostPort(addr); err == nil {
		if host != "" || portStr != "" {
			if port, err := strconv.Atoi(portStr); err == nil {
				return port
			}
		}
	}
	return 0
}
