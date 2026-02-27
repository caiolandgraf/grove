package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion script",
	Long: bold("completion") + ` generates a shell completion script for grove.

Load it in your current shell session or persist it so it applies to every
new shell. Pick the snippet for your shell below.

` + colorBold + colorGray + `Bash` + colorReset + `
  # Load once in the current session:
  source <(grove completion bash)

  # Persist (add to ~/.bashrc):
  echo 'source <(grove completion bash)' >> ~/.bashrc

` + colorBold + colorGray + `Zsh` + colorReset + `
  # Load once in the current session:
  source <(grove completion zsh)

  # Persist (add to ~/.zshrc):
  echo 'source <(grove completion zsh)' >> ~/.zshrc

  # Or install to your fpath:
  grove completion zsh > "${fpath[1]}/_grove"

` + colorBold + colorGray + `Fish` + colorReset + `
  # Load once in the current session:
  grove completion fish | source

  # Persist:
  grove completion fish > ~/.config/fish/completions/grove.fish

` + colorBold + colorGray + `PowerShell` + colorReset + `
  # Load once in the current session:
  grove completion powershell | Out-String | Invoke-Expression

  # Persist (add to your PowerShell profile):
  grove completion powershell >> $PROFILE`,
	ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			return rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			return rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			return rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		default:
			return fmt.Errorf(
				"unsupported shell %q — choose: bash, zsh, fish, powershell",
				args[0],
			)
		}
	},
}
