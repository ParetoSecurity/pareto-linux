package shared

import (
	"testing"
)

func TestReadFile(t *testing.T) {
	// Mock data for testing
	ReadFileMocks = map[string]string{
		"testfile1.txt": "This is a test file content",
	}

	t.Run("ReadFile from mock", func(t *testing.T) {
		content, err := ReadFile("testfile1.txt")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		expected := "This is a test file content"
		if string(content) != expected {
			t.Fatalf("expected %s, got %s", expected, string(content))
		}
	})

	t.Run("ReadFile mock not found", func(t *testing.T) {
		_, err := ReadFile("nonexistent.txt")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		expectedErr := "ReadFile fixture not found: nonexistent.txt"
		if err.Error() != expectedErr {
			t.Fatalf("expected %s, got %s", expectedErr, err.Error())
		}
	})

}
