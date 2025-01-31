package checks

import (
	"os"
	"testing"
)

var osStatMock map[string]bool

// osStat checks if a file exists by attempting to get its file info.
// During testing, it uses a mock implementation via osStatMock.
// It returns the file path if the file exists, otherwise returns an empty string and error.
func osStat(file string) (string, error) {
	if testing.Testing() {
		if found := osStatMock[file]; found {
			return file, nil
		}
		return "", os.ErrNotExist
	}
	_, err := os.Stat(file)
	if err != nil {
		return "", err
	}
	return file, nil
}
