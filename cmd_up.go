package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Start docker compose and run the dev server",
	Long: bold(
		"up",
	) + ` starts infrastructure via ` + colorCyan + `docker compose` + colorReset + `
then launches ` + colorGreen + `grove dev` + colorReset + `.

` + colorGray + `Examples:` + colorReset + `
  grove up`,
	RunE: runUp,
}

func runUp(_ *cobra.Command, _ []string) error {
	fmt.Println()

	if err := runDockerComposeUp(); err != nil {
		return err
	}

	return runDev(nil, nil)
}

func runDockerComposeUp() error {
	composePath, err := resolveComposePath(".")
	if err != nil {
		return err
	}

	envFile := resolveEnvFilePath()

	command := "docker compose -f " + composePath
	if envFile != "" {
		command += " --env-file " + envFile
	}
	command += " up -d"

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
		args = append(args, "up", "-d")
		if err := runComposeCommand("docker", args...); err == nil {
			fmt.Println(done("Docker compose is up"))
			fmt.Println()
			return nil
		} else {
			composeErr := err
			if _, err := exec.LookPath("docker-compose"); err == nil {
				args := []string{"-f", composePath}
				if envFile != "" {
					args = append(args, "--env-file", envFile)
				}
				args = append(args, "up", "-d")
				if err := runComposeCommand(
					"docker-compose",
					args...); err == nil {
					fmt.Println(done("Docker compose is up"))
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
		args = append(args, "up", "-d")
		if err := runComposeCommand("docker-compose", args...); err == nil {
			fmt.Println(done("Docker compose is up"))
			fmt.Println()
			return nil
		}
		return fmt.Errorf("docker-compose failed")
	}

	return fmt.Errorf("docker is not installed or not in PATH")
}

func runComposeCommand(name string, args ...string) error {
	c := exec.Command(name, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	return c.Run()
}

func resolveEnvFilePath() string {
	if _, err := os.Stat(".env"); err == nil {
		return ".env"
	}
	return ""
}
