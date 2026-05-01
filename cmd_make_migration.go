package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var makeMigrationEnv string

var makeMigrationCmd = &cobra.Command{
	Use:   "make:migration <name>",
	Short: "Generate a new migration via atlas migrate diff",
	Long: bold(
		"make:migration",
	) + ` generates a new SQL migration file by diffing your GORM
	models against the current database schema using Atlas.

	` + colorYellow + `Important:` + colorReset + ` always edit your model in ` + colorCyan + `internal/models/` + colorReset + ` first,
	then run this command — Atlas will produce the exact SQL diff between your
	updated struct and the current database schema.

` + colorGray + `Examples:` + colorReset + `
  grove make:migration add_posts_table
  grove make:migration add_email_to_users
  grove make:migration create_orders_table --env dev`,
	Args: cobra.ExactArgs(1),
	RunE: runMakeMigration,
}

func init() {
	makeMigrationCmd.Flags().StringVar(
		&makeMigrationEnv,
		"env", "local",
		"Atlas environment to use (local, dev, production)",
	)
}

func runMakeMigration(cmd *cobra.Command, args []string) error {
	name := args[0]
	// Normalize: spaces → underscores, lowercase
	name = strings.ToLower(strings.ReplaceAll(name, " ", "_"))
	name = strings.ReplaceAll(name, "-", "_")

	fmt.Println()
	fmt.Printf(
		"  %sGenerating migration%s %s %s\n",
		colorGray, colorReset,
		bold(name),
		gray("(atlas migrate diff --env "+makeMigrationEnv+")"),
	)
	fmt.Println()

	// Check atlas is available
	if _, err := exec.LookPath("atlas"); err != nil {
		return fmt.Errorf(
			"atlas CLI not found in PATH\n\n  Install it from: %s",
			colorCyan+"https://atlasgo.io/docs"+colorReset,
		)
	}

	if err := ensureAtlasGormProvider(); err != nil {
		return err
	}

	atlasArgs := []string{
		"migrate", "diff", name,
		"--env", makeMigrationEnv,
	}

	aw := newAtlasOutputWriter(os.Stdout)

	c := exec.Command("atlas", atlasArgs...)
	c.Stdout = aw
	c.Stderr = aw
	c.Stdin = os.Stdin

	err := c.Run()
	aw.Flush()

	if err != nil {
		return fmt.Errorf("atlas migrate diff failed: %w", err)
	}

	fmt.Println()
	fmt.Println(done(
		"Migration " + bold(
			name,
		) + " created in " + colorCyan + "migrations/" + colorReset,
	))
	fmt.Println()
	fmt.Println(nextSteps())
	fmt.Printf(
		"    %s1.%s Review the generated SQL in %s\n",
		colorGray, colorReset,
		colorCyan+"migrations/"+colorReset,
	)
	fmt.Printf(
		"    %s2.%s Run %s to apply it\n",
		colorGray, colorReset,
		colorGreen+"grove migrate"+colorReset,
	)
	fmt.Println()

	return nil
}

func ensureAtlasGormProvider() error {
	const providerExec = "./.grove/bin/atlas-gorm"
	const providerPath = ".grove/bin/atlas-gorm"

	if _, err := os.Stat("atlas.hcl"); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read atlas.hcl: %w", err)
	}

	if _, err := os.Stat(providerPath); err != nil {
		if os.IsNotExist(err) {
			fmt.Printf(
				"  %sPreparing Atlas provider%s %s\n",
				colorGray, colorReset,
				gray("(ariga.io/atlas-provider-gorm)"),
			)
			if err := ensureDir(filepath.Dir(providerPath)); err != nil {
				return fmt.Errorf("failed to create provider dir: %w", err)
			}
			if _, err := exec.LookPath("go"); err != nil {
				return fmt.Errorf("go binary not found in PATH")
			}

			bw := newBuildOutputWriter(os.Stderr)
			c := exec.Command(
				"go",
				"build",
				"-o",
				providerPath,
				"ariga.io/atlas-provider-gorm",
			)
			c.Stdout = bw
			c.Stderr = bw

			if err := c.Run(); err != nil {
				fmt.Println()
				return fmt.Errorf("failed to build atlas provider: %w", err)
			}
		} else {
			return fmt.Errorf("failed to stat %s: %w", providerPath, err)
		}
	}

	updated, err := updateAtlasHCLProgram(providerExec)
	if err != nil {
		return err
	}
	if updated {
		fmt.Printf(
			"  %sUsing cached Atlas provider%s %s\n",
			colorGray, colorReset,
			gray("("+providerExec+")"),
		)
	}

	return nil
}

func updateAtlasHCLProgram(providerExec string) (bool, error) {
	path := "atlas.hcl"
	content, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}

	src := string(content)
	if strings.Contains(src, providerExec) {
		return false, nil
	}

	lines := strings.Split(src, "\n")
	inGorm := false

	for i := 0; i < len(lines); i++ {
		trimmed := strings.TrimSpace(lines[i])
		if strings.HasPrefix(trimmed, "data \"external_schema\" \"gorm\"") {
			inGorm = true
			continue
		}

		if inGorm && strings.HasPrefix(trimmed, "program = [") {
			idx := strings.Index(lines[i], "program")
			if idx < 0 {
				return false, fmt.Errorf("malformed program line in %s", path)
			}
			indent := lines[i][:idx]
			j := i + 1
			var tokens []string
			for j < len(lines) {
				line := strings.TrimSpace(lines[j])
				if line == "]" {
					break
				}
				if line != "" {
					token := strings.TrimSuffix(line, ",")
					token = strings.Trim(token, "\"")
					if token != "" {
						tokens = append(tokens, token)
					}
				}
				j++
			}

			if j >= len(lines) {
				return false, fmt.Errorf("malformed program block in %s", path)
			}

			loadIndex := -1
			for idx, token := range tokens {
				if token == "load" {
					loadIndex = idx
					break
				}
			}

			var tail []string
			if loadIndex >= 0 && loadIndex+1 < len(tokens) {
				tail = tokens[loadIndex+1:]
			}

			newTokens := []string{providerExec, "load"}
			newTokens = append(newTokens, tail...)

			var newLines []string
			newLines = append(newLines, indent+"program = [")
			for _, token := range newTokens {
				newLines = append(newLines, indent+"  \""+token+"\",")
			}
			newLines = append(newLines, indent+"]")

			lines = append(lines[:i], append(newLines, lines[j+1:]...)...)
			if err := os.WriteFile(
				path,
				[]byte(strings.Join(lines, "\n")),
				0o644,
			); err != nil {
				return false, err
			}
			return true, nil
		}

		if inGorm && trimmed == "}" {
			inGorm = false
		}
	}

	fmt.Println(warn("Could not update atlas.hcl to use cached provider."))
	return false, nil
}
