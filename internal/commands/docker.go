package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Gets the docker-compose.yml file path from the clients .config/fluxstream directory
func getDockerComposeFilePath() (string, error) {
	configDir, err := os.UserConfigDir()

	if err != nil {
		return "", err
	}

	dockerComposeFp := filepath.Join(configDir, "fluxstream", "docker-compose.yml")

	return dockerComposeFp, nil
}

// Checks whether docker is installed or not with docker -v
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

// DockerCompose runs a docker-compose command with the given arguments
func DockerCompose(args ...string) error {
	dockerComposeFilePath, err := getDockerComposeFilePath()
	if err != nil {
		return fmt.Errorf("failed to get docker-compose file path: %v", err)
	}

	// Prepend the -f flag to always use our config
	cmdArgs := append([]string{"-f", dockerComposeFilePath}, args...)
	cmd := exec.Command("docker-compose", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// DockerComposeOutput runs a docker-compose command and returns its stdout as string
func DockerComposeOutput(args ...string) (string, error) {
	dockerComposeFilePath, err := getDockerComposeFilePath()
	if err != nil {
		return "", fmt.Errorf("failed to get docker-compose file path: %v", err)
	}

	cmdArgs := append([]string{"-f", dockerComposeFilePath}, args...)
	cmd := exec.Command("docker-compose", cmdArgs...)
	output, err := cmd.CombinedOutput()

	return string(output), err
}

// DockerRunning checks if Docker is installed and daemon is running
func DockerRunning() error {
	if !IsDockerInstalled() {
		return fmt.Errorf("docker is not installed")
	}
	if err := DockerInfo(); err != nil {
		return fmt.Errorf("docker daemon not running: %v", err)
	}
	return nil
}

// DockerStart starts FluxStream containers
func DockerStart(detached bool) error {
	args := []string{"up"}
	if detached {
		args = append(args, "-d")
	}
	return DockerCompose(args...)
}

// DockerStop stops all FluxStream containers
func DockerStop() error {
	return DockerCompose("down")
}

// DockerStatus shows container list
func DockerStatus() error {
	return DockerCompose("ps")
}

// DockerLogs streams logs
func DockerLogs() error {
	return DockerCompose("logs", "-f")
}
