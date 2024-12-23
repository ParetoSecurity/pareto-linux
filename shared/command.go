package shared

import (
	"errors"
	"os/exec"
	"strings"
	"testing"
)

// RunCommandMocks is a map that stores mock command outputs.
// The key is the command string, and the value is the corresponding mock output.
var RunCommandMocks map[string]string

// RunCommand executes a command with the given name and arguments, and returns
// the combined standard output and standard error as a string. If testing is
// enabled, it returns a predefined fixture instead of executing the command.
func RunCommand(name string, arg ...string) (string, error) {

	// Check if testing is enabled and enable harnessing
	if testing.Testing() {
		fx := name + " " + strings.Join(arg, " ")
		fixtureFile, ok := RunCommandMocks[fx]
		if !ok {
			return "", errors.New("RunCommand fixture not found: " + fx)
		}
		return fixtureFile, nil
	}

	cmd := exec.Command(name, arg...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}
