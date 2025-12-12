package files

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopy_Success_SingleFile(t *testing.T) {
	// Create temporary source directory
	srcDir := t.TempDir()
	dstDir := t.TempDir()

	// Create a file in source directory
	srcFile := filepath.Join(srcDir, "test.txt")
	testContent := "Hello, World!"

	err := os.WriteFile(srcFile, []byte(testContent), 0644)
	assert.NoError(t, err, "Failed to create source file")

	copier := NewCopier()
	err = copier.Copy(srcDir, dstDir)
	assert.NoError(t, err, "Copy failed")

	// Verify the file was copied
	dstFile := filepath.Join(dstDir, "test.txt")
	data, err := os.ReadFile(dstFile)
	assert.NoError(t, err, "Failed to read copied file")
	assert.Equal(t, testContent, string(data), "file content should be the same")
}

func TestCopy_Success_NestedDirectories(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := t.TempDir()

	// Create nested directory structure
	nestedDir := filepath.Join(srcDir, "subdir", "nested")
	err := os.MkdirAll(nestedDir, 0750)
	assert.NoError(t, err, "Failed to create nested directory")

	// Create files in nested directories
	rootFile := filepath.Join(srcDir, "root.txt")
	nestedFile := filepath.Join(nestedDir, "nested.txt")

	err = os.WriteFile(rootFile, []byte("root content"), 0644)
	assert.NoError(t, err, "Failed to create root file")
	err = os.WriteFile(nestedFile, []byte("nested content"), 0644)
	assert.NoError(t, err, "Failed to create nested file")

	copier := NewCopier()
	err = copier.Copy(srcDir, dstDir)
	assert.NoError(t, err, "Copy failed")

	// Verify root file
	dstRootFile := filepath.Join(dstDir, "root.txt")
	data, err := os.ReadFile(dstRootFile)
	assert.NoError(t, err, "Failed to read root file")
	assert.Equal(t, "root content", string(data), "root file content should be the same")

	// Verify nested file
	dstNestedFile := filepath.Join(dstDir, "subdir", "nested", "nested.txt")
	data, err = os.ReadFile(dstNestedFile)
	assert.NoError(t, err, "Failed to read nested file")
	assert.Equal(t, "nested content", string(data), "nested file content should be the same")
}

func TestCopy_Error_SourceNotExists(t *testing.T) {
	dstDir := t.TempDir()
	nonExistentSrc := filepath.Join(dstDir, "nonexistent")

	copier := NewCopier()
	err := copier.Copy(nonExistentSrc, dstDir)
	assert.Error(t, err, "should return error when non existent source")
}

func TestCopy_Error_InvalidDestination(t *testing.T) {
	srcDir := t.TempDir()

	// Create a file in source
	srcFile := filepath.Join(srcDir, "test.txt")
	err := os.WriteFile(srcFile, []byte("content"), 0644)
	assert.NoError(t, err, "Failed to create source file")

	invalidDst := "/nonexistent/path/that/does/not/exist/destination"

	copier := NewCopier()
	err = copier.Copy(srcDir, invalidDst)
	assert.Error(t, err, "should return error when non existent destination")
}

func TestCopy_Success_EmptyDirectory(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := t.TempDir()

	// Copy empty directory
	copier := NewCopier()
	err := copier.Copy(srcDir, dstDir)
	assert.NoError(t, err, "Copy of empty directory shouldn't failed")

	// Verify destination exists
	_, err = os.Stat(dstDir)
	assert.NoError(t, err, "Destination directory should exist")
}
