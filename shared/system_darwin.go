package shared

import (
	"fmt"

	"os/exec"
	"strings"
)

func OSVersion() (string, error) {
	cmd := exec.Command("sw_vers", "-productVersion")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	version := strings.TrimSpace(string(output))
	if version == "" {
		return "", fmt.Errorf("unable to retrieve macOS version")
	}

	return version, nil
}
