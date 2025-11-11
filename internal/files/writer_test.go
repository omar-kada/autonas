package files

import (
	"os"
	"testing"
)

func TestWriteToFile_Success(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	content := "KEY1=VALUE1\nKEY2=VALUE2"
	err = NewWriter().WriteToFile(tmpFile.Name(), content)
	if err != nil {
		t.Fatalf("WriteToFile failed: %v", err)
	}

	data, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read temp file: %v", err)
	}

	if string(data) != content {
		t.Errorf("Expected file content %q, got %q", content, string(data))
	}
}

func TestWriteToFile_Error(t *testing.T) {
	// Attempt to write to an invalid path
	invalidPath := "/invalid_path/testfile.env"
	err := NewWriter().WriteToFile(invalidPath, "KEY=VALUE")
	if err == nil {
		t.Fatalf("Expected error when writing to invalid path, got nil")
	}
}
