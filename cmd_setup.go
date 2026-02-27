package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// ──────────────────────────────────────────────
// Constants
// ──────────────────────────────────────────────

const (
	setupTemplateRepo   = "caiolandgraf/grove-base"
	setupTemplateBranch = "main"
)

// ──────────────────────────────────────────────
// Command definition
// ──────────────────────────────────────────────

var setupModuleFlag string

var setupCmd = &cobra.Command{
	Use:   "setup <project-name>",
	Short: "Scaffold a new Grove project from the official template",
	Long: bold("setup") + ` downloads and scaffolds a complete Grove project
from the official template repository on GitHub.

` + colorGray + `Examples:` + colorReset + `
  grove setup my-api
  grove setup my-api --module github.com/acme/my-api`,
	Args: cobra.ExactArgs(1),
	RunE: runSetup,
}

func init() {
	setupCmd.Flags().StringVar(
		&setupModuleFlag,
		"module", "",
		"Go module path (defaults to project name)",
	)
}

// ──────────────────────────────────────────────
// Spinner
// ──────────────────────────────────────────────

var spinFrames = []rune{'⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'}

type step struct {
	label  string
	stopCh chan struct{}
	doneCh chan struct{}
}

func startStep(label string) *step {
	s := &step{
		label:  label,
		stopCh: make(chan struct{}),
		doneCh: make(chan struct{}),
	}
	go func() {
		defer close(s.doneCh)
		tick := time.NewTicker(80 * time.Millisecond)
		defer tick.Stop()
		i := 0
		for {
			select {
			case <-s.stopCh:
				return
			case <-tick.C:
				fmt.Printf(
					"\r\033[2K    %s%c%s  %s",
					colorCyan, spinFrames[i%len(spinFrames)], colorReset,
					s.label,
				)
				i++
			}
		}
	}()
	return s
}

func (s *step) succeed(extra string) {
	close(s.stopCh)
	<-s.doneCh
	if extra != "" {
		fmt.Printf(
			"\r\033[2K    %s✓%s  %-36s %s\n",
			colorGreen, colorReset,
			s.label,
			colorDim+extra+colorReset,
		)
	} else {
		fmt.Printf("\r\033[2K    %s✓%s  %s\n", colorGreen, colorReset, s.label)
	}
}

func (s *step) fail(extra string) {
	close(s.stopCh)
	<-s.doneCh
	line := s.label
	if extra != "" {
		line += "  " + colorDim + extra + colorReset
	}
	fmt.Printf("\r\033[2K    %s✕%s  %s\n", colorRed, colorReset, line)
}

// ──────────────────────────────────────────────
// Main runner
// ──────────────────────────────────────────────

func runSetup(_ *cobra.Command, args []string) error {
	projectName := args[0]

	modulePath := setupModuleFlag
	if modulePath == "" {
		modulePath = projectName
	}

	// Validate target directory
	if _, err := os.Stat(projectName); err == nil {
		return fmt.Errorf("directory %q already exists", projectName)
	}

	// Hide cursor for the entire setup flow
	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h")

	// Cleanup on failure
	succeeded := false
	defer func() {
		if !succeeded {
			_ = os.RemoveAll(projectName)
		}
	}()

	printSetupHeader(projectName, modulePath)

	// ── Step 1: Download ───────────────────────────────────────────────────
	s := startStep("Downloading template")
	zipPath, size, err := downloadTemplate()
	if err != nil {
		s.fail(err.Error())
		return fmt.Errorf("download failed: %w", err)
	}
	s.succeed(fmtBytes(size))
	defer os.Remove(zipPath)

	// ── Step 2: Extract ────────────────────────────────────────────────────
	s = startStep("Extracting files")
	fileCount, err := extractTemplate(zipPath, projectName)
	if err != nil {
		s.fail(err.Error())
		return fmt.Errorf("extraction failed: %w", err)
	}
	s.succeed(fmt.Sprintf("%d files", fileCount))

	// ── Step 3: Configure module ───────────────────────────────────────────
	s = startStep("Configuring module")
	if err := configureModule(projectName, modulePath); err != nil {
		s.fail(err.Error())
		return fmt.Errorf("configuration failed: %w", err)
	}
	s.succeed(modulePath)

	// ── Step 4: Install dependencies ──────────────────────────────────────
	s = startStep("Installing dependencies")
	start := time.Now()
	if err := runGoModTidy(projectName); err != nil {
		s.fail(err.Error())
		return fmt.Errorf("go mod tidy failed: %w", err)
	}
	s.succeed(fmtDuration(time.Since(start)))

	succeeded = true
	printSetupSuccess(projectName)
	return nil
}

// ──────────────────────────────────────────────
// Download
// ──────────────────────────────────────────────

func downloadTemplate() (tmpFile string, size int64, err error) {
	url := fmt.Sprintf(
		"https://github.com/%s/archive/refs/heads/%s.zip",
		setupTemplateRepo, setupTemplateBranch,
	)

	resp, err := http.Get(url) //nolint:noctx
	if err != nil {
		return "", 0, fmt.Errorf("network error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("HTTP %s", resp.Status)
	}

	f, err := os.CreateTemp("", "grove-setup-*.zip")
	if err != nil {
		return "", 0, err
	}
	defer f.Close()

	size, err = io.Copy(f, resp.Body)
	if err != nil {
		_ = os.Remove(f.Name())
		return "", 0, fmt.Errorf("download error: %w", err)
	}

	return f.Name(), size, nil
}

// ──────────────────────────────────────────────
// Extraction
// ──────────────────────────────────────────────

// setupSkipPaths are paths (relative to the repo root) that should never be
// copied into the new project.
var setupSkipPaths = []string{
	"bin/",
	"tmp/",
	"grove", // compiled binary at repo root
}

func extractTemplate(zipPath, destDir string) (int, error) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return 0, err
	}
	defer r.Close()

	// GitHub ZIPs always have a single top-level directory named
	// "{repo}-{branch}/".  Find it so we can strip it.
	prefix := ""
	for _, f := range r.File {
		if idx := strings.Index(f.Name, "/"); idx >= 0 {
			prefix = f.Name[:idx+1]
			break
		}
	}

	count := 0
	for _, f := range r.File {
		rel := strings.TrimPrefix(f.Name, prefix)
		if rel == "" {
			continue
		}
		if setupShouldSkip(rel) {
			continue
		}

		dest := filepath.Join(destDir, filepath.FromSlash(rel))

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(dest, 0o755); err != nil {
				return count, err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return count, err
		}

		if err := extractFile(f, dest); err != nil {
			return count, err
		}
		count++
	}

	return count, nil
}

func extractFile(f *zip.File, dest string) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	out, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, rc)
	return err
}

func setupShouldSkip(rel string) bool {
	for _, skip := range setupSkipPaths {
		if rel == skip || strings.HasPrefix(rel, skip) {
			return true
		}
	}
	return false
}

// ──────────────────────────────────────────────
// Module configuration
// ──────────────────────────────────────────────

func configureModule(projectDir, newModule string) error {
	goModPath := filepath.Join(projectDir, "go.mod")

	raw, err := os.ReadFile(goModPath)
	if err != nil {
		return fmt.Errorf("cannot read go.mod: %w", err)
	}

	// Extract the original module name declared in go.mod
	oldModule := ""
	for _, line := range strings.Split(string(raw), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "module ") {
			oldModule = strings.TrimSpace(
				strings.TrimPrefix(trimmed, "module "),
			)
			break
		}
	}
	if oldModule == "" {
		return fmt.Errorf("module directive not found in go.mod")
	}
	if oldModule == newModule {
		return nil // nothing to do
	}

	return filepath.WalkDir(
		projectDir,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				if d.Name() == ".git" {
					return filepath.SkipDir
				}
				return nil
			}

			name := d.Name()
			if name != "go.mod" && !strings.HasSuffix(name, ".go") {
				return nil
			}

			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			replaced := strings.ReplaceAll(
				string(content),
				oldModule,
				newModule,
			)
			if replaced == string(content) {
				return nil // nothing changed
			}

			return os.WriteFile(path, []byte(replaced), 0o644)
		},
	)
}

// ──────────────────────────────────────────────
// go mod tidy
// ──────────────────────────────────────────────

func runGoModTidy(projectDir string) error {
	if _, err := exec.LookPath("go"); err != nil {
		return fmt.Errorf("go binary not found in PATH")
	}

	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = projectDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		msg := strings.TrimSpace(string(out))
		if msg == "" {
			msg = err.Error()
		}
		return fmt.Errorf("%s", msg)
	}
	return nil
}

// ──────────────────────────────────────────────
// UI helpers
// ──────────────────────────────────────────────

func printSetupHeader(projectName, modulePath string) {
	sep := "  " + colorDim + strings.Repeat("─", 54) + colorReset

	logo := "\n" +
		logoG1 + `  █▀▀ █▀█ █▀█ █░█ █▀▀  ` + colorReset + "\n" +
		logoG3 + `  █▄█ █▀▄ █▄█ ▀▄▀ ██▄  ` + colorReset + "\n"

	fmt.Println(logo)
	fmt.Printf(
		"  %sProject%s   %s\n",
		colorBold+colorGray, colorReset,
		bold(projectName),
	)
	fmt.Printf(
		"  %sModule%s    %s\n",
		colorBold+colorGray, colorReset,
		colorCyan+modulePath+colorReset,
	)
	fmt.Printf(
		"  %sTemplate%s  %s\n",
		colorBold+colorGray, colorReset,
		colorDim+setupTemplateRepo+colorReset,
	)
	fmt.Println()
	fmt.Println(sep)
	fmt.Println()
}

func printSetupSuccess(projectName string) {
	sep := "  " + colorDim + strings.Repeat("─", 54) + colorReset

	fmt.Println()
	fmt.Println(sep)
	fmt.Println()
	fmt.Println(done("Project created in " + bold("./"+projectName)))
	fmt.Println()
	fmt.Println(nextSteps())
	fmt.Printf("    %scd %s%s\n", colorGreen, projectName, colorReset)
	fmt.Printf("    %scp .env.example .env%s\n", colorGreen, colorReset)
	fmt.Printf("    %sgrove serve%s\n", colorGreen, colorReset)
	fmt.Println()
}

// ──────────────────────────────────────────────
// Formatting helpers
// ──────────────────────────────────────────────

func fmtBytes(n int64) string {
	switch {
	case n >= 1<<20:
		return fmt.Sprintf("%.1f MB", float64(n)/float64(1<<20))
	case n >= 1<<10:
		return fmt.Sprintf("%.0f KB", float64(n)/float64(1<<10))
	default:
		return fmt.Sprintf("%d B", n)
	}
}

func fmtDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	return fmt.Sprintf("%.1fs", d.Seconds())
}
