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
		return fmt.Errorf("docker not installed")
	}

	if err := DockerInfo(); err != nil {
		return fmt.Errorf("docker daemon not running — please ensure Docker is started: %v", err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configdir := filepath.Join(homeDir, ".config", "fluxstream")
	if err := os.MkdirAll(configdir, 0o700); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	datadir := filepath.Join(homeDir, ".local", "share", "fluxstream", "downloads")
	if err := os.MkdirAll(datadir, 0o750); err != nil {
		return fmt.Errorf("failed to create data directories: %v", err)
	}

	replaced := strings.ReplaceAll(dockerComposeTemplate, "{{DOWNLOAD_PATH}}", datadir)
	if err := os.WriteFile(filepath.Join(configdir, "docker-compose.yml"), []byte(replaced), 0o600); err != nil {
		return fmt.Errorf("failed to write docker-compose.yml: %v", err)
	}

	fmt.Println("Fluxstream setup complete!")
	fmt.Printf("Config: %s\nData: %s\n", configdir, datadir)
	return nil
}
