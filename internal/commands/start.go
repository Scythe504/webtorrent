package commands

import (
	"os"
	"os/exec"
)


func Start() error {
	dockerComposeFilePath, err := getDockerComposeFilePath()

	if err != nil {
		return err
	}

	cmd := exec.Command("docker-compose", "-f", dockerComposeFilePath, "up")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
