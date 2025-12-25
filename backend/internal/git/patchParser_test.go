package git

import (
	"strings"
	"testing"

	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"
	"github.com/stretchr/testify/assert"
)

// Create a mock patch and commit
var diffFile1 = strings.Join(
	[]string{
		"diff --git a/file1.txt b/file2.txt",
		"",
		"index 000000..111111 100644",
		"--- a/file1.txt\n+++ b/file2.txt",
		"@@ -1 +1 @@",
		"-hello",
		"+world",
	},
	"\n",
)

var diffFile2 = strings.Join(
	[]string{
		"diff --git a/another_file.txt b/another_file.txt",
		"",
		"index 000000..111111 100644",
		"--- a/file1.txt\n+++ b/file2.txt",
		"@@ -1 +1 @@",
		"-hello",
		"+world",
	},
	"\n",
)

func TestParse(t *testing.T) {
	mockPatch := diffFile1 + "\n" + diffFile2
	mockCommit := &object.Commit{
		Message: "Test commit message",
		Author: object.Signature{
			Name: "Test Author",
		},
		Hash: plumbing.NewHash("a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3"),
	}

	// Create a new patch parser
	parser := NewPatchParser()

	// Call the Parse method
	patch, err := parser.Parse(mockPatch, mockCommit)

	// Assert no error occurred
	assert.NoError(t, err)

	// Assert the patch fields are correctly set
	assert.Equal(t, mockCommit.Message, patch.Title)
	assert.Equal(t, mockPatch, patch.Diff)
	assert.Equal(t, mockCommit.Author.Name, patch.Author)
	assert.Equal(t, mockCommit.Hash.String(), patch.CommitHash)
	assert.Equal(t, diffFile1, patch.Files[0].Diff)
	assert.Equal(t, "file2.txt", patch.Files[0].NewFile)
	assert.Equal(t, "file1.txt", patch.Files[0].OldFile)
	assert.Equal(t, diffFile2, patch.Files[1].Diff)
	assert.Equal(t, "another_file.txt", patch.Files[1].NewFile)
	assert.Equal(t, "another_file.txt", patch.Files[1].OldFile)
}

func TestToFileDiff(t *testing.T) {
	// Test case 1: Valid diff string
	fileDiff, err := toFileDiff(diffFile1)

	assert.NoError(t, err)
	assert.Equal(t, "file1.txt", fileDiff.OldFile)
	assert.Equal(t, "file2.txt", fileDiff.NewFile)
	assert.Equal(t, diffFile1, fileDiff.Diff)

	// Test case 2: Diff string with less than 2 lines
	diffStr := "diff --git a/file1.txt b/file2.txt"
	fileDiff, err = toFileDiff(diffStr)

	assert.Error(t, err)
	assert.Equal(t, "diff contains less than 2 lines", err.Error())

	// Test case 3: Diff string with no file names
	diffStr = "diff --git\n\nindex 000000..111111 100644\n--- a/file1.txt\n+++ b/file2.txt\n@@ -1 +1 @@\n-hello\n+world\n"
	fileDiff, err = toFileDiff(diffStr)

	assert.Error(t, err)
	assert.Equal(t, "can't find file names", err.Error())
}
