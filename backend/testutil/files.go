// Package testutil contains utility functions to be used for testing
package testutil

import (
	"embed"
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/tools/txtar"
	"gopkg.in/yaml.v3"
)

// GetTempTestFile writes an embedded fixture (from testDataFS) to a temp file and returns its path.
func GetTempTestFile(t *testing.T, testDataFS embed.FS, name string) string {
	t.Helper()
	bs, err := testDataFS.ReadFile("test_data/" + name)
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	dst := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(dst, bs, 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	return dst
}

// File represents a temporary file used in tests
type File struct {
	Name string
	Path string
}

// ExtractTxtar extracts config files & expected result from a txtar test archive.
func ExtractTxtar(t *testing.T, testDataFS embed.FS, archivePath string) []File {
	t.Helper()
	ar, err := txtar.ParseFile(GetTempTestFile(t, testDataFS, archivePath))
	if err != nil {
		t.Fatalf("failed to parse txtar file: %v", err)
	}
	result := make([]File, 0, len(ar.Files))
	tempDir := t.TempDir()
	for _, file := range ar.Files {
		outPath := filepath.Join(tempDir, file.Name)
		if err := os.MkdirAll(filepath.Dir(outPath), 0750); err != nil {
			t.Fatalf("failed to create tmp directory: %v", err)
		}
		if err := os.WriteFile(outPath, file.Data, 0600); err != nil {
			t.Fatalf("failed to create tmp file: %v", err)
		}
		result = append(result, File{
			Name: file.Name,
			Path: outPath,
		})
	}
	return result
}

// ReadYamlFile reads a YAML file from the given path and unmarshals it into a map.
func ReadYamlFile(t *testing.T, path string) map[string]any {
	t.Helper()

	result := make(map[string]any)
	bs := ReadFile(t, path)
	err := yaml.Unmarshal(bs, &result)
	if err != nil {
		t.Fatalf("failed to parse want yaml file: %v", err)
	}
	return result
}

// ReadFile reads file byte content and returns it
func ReadFile(t *testing.T, path string) []byte {
	t.Helper()

	bs, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read want file: %v", err)
	}
	return bs
}
