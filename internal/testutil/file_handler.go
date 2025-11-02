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

// ExtractTxtar extracts config files & expected result from a txtar test archive.
func ExtractTxtar(t *testing.T, testDataFS embed.FS, archivePath string) map[string]string {
	t.Helper()

	ar, err := txtar.ParseFile(GetTempTestFile(t, testDataFS, archivePath))
	if err != nil {
		t.Fatalf("failed to parse txtar file: %v", err)
	}
	result := make(map[string]string)
	tempDir := t.TempDir()
	for _, file := range ar.Files {
		outPath := filepath.Join(tempDir, file.Name)
		if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
			t.Fatalf("failed to create tmp directory: %v", err)
		}
		if err := os.WriteFile(outPath, file.Data, 0o600); err != nil {
			t.Fatalf("failed to create tmp file: %v", err)
		}
		result[file.Name] = outPath
	}
	return result
}

// ReadYamlFile reads a YAML file from the given path and unmarshals it into a map.
func ReadYamlFile(t *testing.T, path string) map[string]any {
	want := make(map[string]any)
	bs, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read want file: %v", err)
	}
	err = yaml.Unmarshal(bs, &want)
	if err != nil {
		t.Fatalf("failed to parse want yaml file: %v", err)
	}
	return want
}
