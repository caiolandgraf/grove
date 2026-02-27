package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

var completionPrint bool

var completionCmd = &cobra.Command{
	Use:       "completion [bash|zsh|fish|powershell]",
	Short:     "Install shell completion or print the script",
	ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Long: bold(
		"completion",
	) + ` installs the grove completion script for your shell.

By default the script is written directly to the correct location for each
shell so tab-completion works in every new session without any extra steps.

Pass ` + colorGreen + `--print` + colorReset + ` (` + colorGreen + `-p` + colorReset + `) to print the script to stdout instead,
useful for piping or manual placement.

` + colorBold + colorGray + `Fish` + colorReset + `
  grove completion fish
  ` + colorDim + `→ $XDG_CONFIG_HOME/fish/completions/grove.fish` + colorReset + `

` + colorBold + colorGray + `Zsh` + colorReset + `
  grove completion zsh
  ` + colorDim + `→ ~/.zsh/completions/_grove` + colorReset + `

` + colorBold + colorGray + `Bash` + colorReset + `
  grove completion bash
  ` + colorDim + `→ ~/.local/share/bash-completion/completions/grove` + colorReset + `

` + colorBold + colorGray + `PowerShell` + colorReset + `
  grove completion powershell
  ` + colorDim + `→ printed to stdout (add to your $PROFILE manually)` + colorReset,
	RunE: runCompletion,
}

func init() {
	completionCmd.Flags().BoolVarP(
		&completionPrint,
		"print", "p", false,
		"Print the script to stdout instead of installing it",
	)
}

func runCompletion(_ *cobra.Command, args []string) error {
	shell := args[0]

	// ── --print: old behaviour, just write to stdout ──────────────────────
	if completionPrint || shell == "powershell" {
		return printCompletion(shell)
	}

	// ── auto-install ───────────────────────────────────────────────────────
	destPath, note, err := completionDestPath(shell)
	if err != nil {
		return err
	}

	if err := ensureDir(filepath.Dir(destPath)); err != nil {
		return fmt.Errorf(
			"cannot create directory %s: %w",
			filepath.Dir(destPath),
			err,
		)
	}

	f, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("cannot write completion file: %w", err)
	}
	defer f.Close()

	if err := writeCompletion(shell, f); err != nil {
		return err
	}

	fmt.Println()
	fmt.Println(done(
		"Completion installed " + gray("→ "+destPath),
	))

	if note != "" {
		fmt.Println()
		fmt.Println(info(note))
	}

	fmt.Println()

	return nil
}

// completionDestPath returns the install path and an optional post-install note
// for the given shell.
func completionDestPath(shell string) (path, note string, err error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", "", fmt.Errorf("cannot determine home directory: %w", err)
	}

	switch shell {
	case "fish":
		configHome := os.Getenv("XDG_CONFIG_HOME")
		if configHome == "" {
			configHome = filepath.Join(home, ".config")
		}
		return filepath.Join(
			configHome,
			"fish",
			"completions",
			"grove.fish",
		), "", nil

	case "zsh":
		dest := filepath.Join(home, ".zsh", "completions", "_grove")
		n := colorCyan + "~/.zsh/completions" + colorReset +
			" must be in your fpath. Add to " +
			colorCyan + "~/.zshrc" + colorReset + " if missing:\n\n" +
			"    " + colorGray + "fpath=(~/.zsh/completions $fpath)" + colorReset + "\n" +
			"    " + colorGray + "autoload -U compinit && compinit" + colorReset
		return dest, n, nil

	case "bash":
		// XDG-compliant location picked up automatically by bash-completion ≥ 2.0
		xdgData := os.Getenv("XDG_DATA_HOME")
		if xdgData == "" {
			xdgData = filepath.Join(home, ".local", "share")
		}
		dest := filepath.Join(
			xdgData,
			"bash-completion",
			"completions",
			"grove",
		)
		var n string
		if runtime.GOOS == "darwin" {
			n = "On macOS, make sure " + colorCyan + "bash-completion@2" + colorReset +
				" is installed:\n\n" +
				"    " + colorGray + "brew install bash-completion@2" + colorReset + "\n" +
				"    " + colorGray + `echo 'export BASH_COMPLETION_COMPAT_DIR="$(brew --prefix)/etc/bash_completion.d"' >> ~/.bashrc` + colorReset + "\n" +
				"    " + colorGray + `echo '[[ -r "$(brew --prefix)/etc/profile.d/bash_completion.sh" ]] && source "$(brew --prefix)/etc/profile.d/bash_completion.sh"' >> ~/.bashrc` + colorReset
		}
		return dest, n, nil
	}

	return "", "", fmt.Errorf("unsupported shell %q", shell)
}

// writeCompletion generates the completion script for shell into w.
func writeCompletion(shell string, w *os.File) error {
	switch shell {
	case "bash":
		return rootCmd.GenBashCompletion(w)
	case "zsh":
		return rootCmd.GenZshCompletion(w)
	case "fish":
		return rootCmd.GenFishCompletion(w, true)
	case "powershell":
		return rootCmd.GenPowerShellCompletionWithDesc(w)
	}
	return fmt.Errorf("unsupported shell %q", shell)
}

// printCompletion writes the completion script to stdout.
func printCompletion(shell string) error {
	switch shell {
	case "bash":
		return rootCmd.GenBashCompletion(os.Stdout)
	case "zsh":
		return rootCmd.GenZshCompletion(os.Stdout)
	case "fish":
		return rootCmd.GenFishCompletion(os.Stdout, true)
	case "powershell":
		fmt.Println()
		fmt.Println(
			info(
				"Add the output below to your PowerShell " + colorCyan + "$PROFILE" + colorReset + ":",
			),
		)
		fmt.Println()
		return rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
	}
	return fmt.Errorf(
		"unsupported shell %q — choose: bash, zsh, fish, powershell",
		shell,
	)
}
