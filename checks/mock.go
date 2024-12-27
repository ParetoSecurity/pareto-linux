package checks

import (
	"os/exec"
	"testing"
)

// checkPortMock is a mock function used for testing purposes. It simulates
// checking the availability of a port for a given protocol. The function
// takes an integer port number and a string representing the protocol
// (e.g., "tcp", "udp") as arguments, and returns a boolean indicating
// whether the port is available (true) or not (false).
var checkPortMock func(port int, proto string) bool

// lookPathMock is a mock function that simulates the behavior of
// the os/exec.LookPath function. It takes a file name as input
// and returns the path to the executable file along with an error
// if the file is not found or any other issue occurs.
var lookPathMock func(file string) (string, error)

func lookPath(file string) (string, error) {
	if testing.Testing() {
		return lookPathMock(file)
	}
	return exec.LookPath(file)
}
