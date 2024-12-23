package shared

import (
	"os"
	"testing"
)

func ReadFile(name string) ([]byte, error) {
	if testing.Testing() {
		fixturePath := "fixtures/" + name
		if _, err := os.Stat(fixturePath); os.IsNotExist(err) {
			content, err := os.ReadFile(name)
			if err != nil {
				return nil, err
			}
			if err := os.WriteFile(fixturePath, content, 0644); err != nil {
				return nil, err
			}
			return content, nil
		}
		return os.ReadFile(fixturePath)
	}
	return os.ReadFile(name)
}
