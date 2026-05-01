package main

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop docker compose services",
	Long: bold(
		"down",
	) + ` stops infrastructure via ` + colorCyan + `docker compose` + colorReset + `.

` + colorGray + `Examples:` + colorReset + `
  grove down`,
	RunE: runDown,
}

func runDown(_ *cobra.Command, _ []string) error {
	fmt.Println()

	if err := runDockerComposeDown(); err != nil {
		return err
	}

	return nil
}

func runDockerComposeDown() error {
	composePath, err := resolveComposePath(".")
	if err != nil {
		return err
	}

	envFile := resolveEnvFilePath()

	command := "docker compose -f " + composePath
	if envFile != "" {
		command += " --env-file " + envFile
	}
	command += " down"

	fmt.Printf(
		"  %s  %s\n",
		badge(colorBgBlue, "DOCKER"),
		gray(command),
	)
	fmt.Println()

	if _, err := exec.LookPath("docker"); err == nil {
		args := []string{"compose", "-f", composePath}
		if envFile != "" {
			args = append(args, "--env-file", envFile)
		}
		args = append(args, "down")
		if err := runComposeCommand("docker", args...); err == nil {
			fmt.Println(done("Docker compose is down"))
			fmt.Println()
			return nil
		} else {
			composeErr := err
			if _, err := exec.LookPath("docker-compose"); err == nil {
				args := []string{"-f", composePath}
				if envFile != "" {
					args = append(args, "--env-file", envFile)
				}
				args = append(args, "down")
				if err := runComposeCommand(
					"docker-compose",
					args...); err == nil {
					fmt.Println(done("Docker compose is down"))
					fmt.Println()
					return nil
				}
			}
			return fmt.Errorf("docker compose failed: %w", composeErr)
		}
	}

	if _, err := exec.LookPath("docker-compose"); err == nil {
		args := []string{"-f", composePath}
		if envFile != "" {
			args = append(args, "--env-file", envFile)
		}
		args = append(args, "down")
		if err := runComposeCommand("docker-compose", args...); err == nil {
			fmt.Println(done("Docker compose is down"))
			fmt.Println()
			return nil
		}
		return fmt.Errorf("docker-compose failed")
	}

	return fmt.Errorf("docker is not installed or not in PATH")
}
