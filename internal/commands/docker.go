package commands

import (
	"os"
	"os/exec"
	"path/filepath"
)

func getDockerComposeFilePath() (string, error) {
	configDir, err := os.UserConfigDir()

	if err != nil {
		return "", err
	}

	dockerComposeFp := filepath.Join(configDir, "fluxstream", "docker-compose.yml")

	return dockerComposeFp, nil
}

func IsDockerInstalled() bool {
	cmd := exec.Command("docker", "--version")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func DockerInfo() error {
	cmd := exec.Command("docker", "info")

	err := cmd.Run()

	return err
}
