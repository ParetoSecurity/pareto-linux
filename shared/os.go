package shared

import (
	"errors"
	"os"
	"testing"
)

// ReadFileMocks is a map that simulates file reading operations by mapping
// file paths (as keys) to their corresponding file contents (as values).
// This can be used for testing purposes to mock the behavior of reading files
// without actually accessing the file system.
var ReadFileMocks map[string]string

// ReadFile reads the content of the file specified by the given name.
// If the code is running in a testing environment, it will return the content
// from the ReadFileMocks map instead of reading from the actual file system.
// If the file name is not found in the ReadFileMocks map, it returns an error.
// Otherwise, it reads the file content from the file system.
func ReadFile(name string) ([]byte, error) {
	if testing.Testing() {
		fixtureFile, ok := ReadFileMocks[name]
		if !ok {
			return []byte(""), errors.New("ReadFile fixture not found: " + name)
		}
		return []byte(fixtureFile), nil

	}
	return os.ReadFile(name)
}
