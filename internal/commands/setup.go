package commands

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//go:embed templates/docker-compose.yml
var dockerComposeTemplate string

func Setup() error {
	if !IsDockerInstalled() {
		printError("Docker is not installed.")
		fmt.Println("\nPlease install Docker first:")
		printInfo("  https://docs.docker.com/desktop/#next-steps")
		printInfo("  https://docs.docker.com/engine/")
		return fmt.Errorf("docker not installed")
	}

	if err := DockerInfo(); err != nil {
		printError("Docker daemon is not running.")
		fmt.Println("\nPlease start Docker and try again:")
		fmt.Println("  - macOS/Windows: Open Docker Desktop")
		fmt.Println("  - Linux: sudo systemctl start docker")
		return fmt.Errorf("docker daemon not running: %v", err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}

	configdir := filepath.Join(homeDir, ".config", "fluxstream")
	if err := os.MkdirAll(configdir, 0o700); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	datadir := filepath.Join(homeDir, ".local", "share", "fluxstream", "downloads")
	if err := os.MkdirAll(datadir, 0o755); err != nil {
		return fmt.Errorf("failed to create data directory: %v", err)
	}

	replaced := strings.ReplaceAll(dockerComposeTemplate, "{{DOWNLOAD_PATH}}", datadir)
	composeFile := filepath.Join(configdir, "docker-compose.yml")
	if err := os.WriteFile(composeFile, []byte(replaced), 0o644); err != nil {
		return fmt.Errorf("failed to write docker-compose.yml: %v", err)
	}

	printSuccess("FluxStream setup complete!")
	fmt.Printf("\n%s %s\n", colorize(colorBlue, "Configuration:"), configdir)
	fmt.Printf("%s %s\n", colorize(colorBlue, "Downloads:"), datadir)
	fmt.Printf("\n%s\n", colorize(colorGreen, "Ready to start! Run: fluxstream start"))
	fmt.Printf("\n%s\n", colorize(colorCyan, "For more info: https://docs.fluxstream.app"))

	return nil
}
