package git

import (
	"fmt"
	"omar-kada/autonas/testutil"
	"os"
	"testing"

	git "github.com/go-git/go-git/v6"
	"github.com/stretchr/testify/assert"
)

func assertBranch(t *testing.T, clonePath string, branch string) {
	t.Helper()
	r, err := git.PlainOpen(clonePath)
	if err != nil {
		t.Fatalf("Failed to open cloned repo: %v", err)
	}

	ref, err := r.Head()
	if err != nil {
		t.Fatalf("Failed to get HEAD of cloned repo: %v", err)
	}

	if ref.Name().Short() != branch {
		t.Errorf("Expected branch '%s', got '%s'", branch, ref.Name().Short())
	}
}

func assertFileContent(t *testing.T, filePath string, wantContent string) {
	t.Helper()
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read %s in cloned repo: %v", filePath, err)
	}

	if string(content) != wantContent {
		t.Errorf("Expected %s to have %s, but got M%s", filePath, wantContent, content)
	}
}

func TestFetch_HappyPath(t *testing.T) {
	fetcher := NewFetcher(os.FileMode(0000))
	remoteRepoPath := testutil.SetupRemoteRepo(t)

	testutil.AddCommitToRepo(t, remoteRepoPath, "README.md", "dummy readme")
	clonePath := t.TempDir() + "/clone-repo"

	patch, err := fetcher.Fetch(remoteRepoPath, "main", clonePath)
	if err != nil {
		t.Fatalf("Fetch failed: %v", err)
	}

	assertFileContent(t, clonePath+"/README.md", "dummy readme")
	assertBranch(t, clonePath, "main")
	assert.True(t, assert.ObjectsAreEqual(Patch{}, patch),
		fmt.Sprintf("patch should be empty but is : %v", patch),
	)

	testutil.AddCommitToRepo(t, remoteRepoPath, "NEWFILE.txt", "new file content")

	patch, err = fetcher.Fetch(remoteRepoPath, "main", clonePath)
	if err != nil {
		t.Fatalf("Fetch failed: %v", err)
	}

	assertFileContent(t, clonePath+"/NEWFILE.txt", "new file content")
	assert.Equal(t, "add NEWFILE.txt", patch.Title)
	assert.Equal(t, "Test", patch.Author)

}

func TestFetch_NoChanges(t *testing.T) {
	fetcher := NewFetcher(os.FileMode(0000))

	remoteRepoPath := testutil.SetupRemoteRepo(t)
	clonePath := t.TempDir() + "/clone-repo"

	_, err := fetcher.Fetch(remoteRepoPath, "main", clonePath)
	if err != nil {
		t.Fatalf("Fetch failed: %v", err)
	}

	patch, err := fetcher.Fetch(remoteRepoPath, "main", clonePath)
	if err != NoErrAlreadyUpToDate {
		t.Fatalf("Expected NoErrAlreadyUpToDate, got: %v", err)
	}
	assert.Equal(t, "", patch.Diff)
}

func TestFetch_NonExistentRepo(t *testing.T) {
	fetcher := NewFetcher(os.FileMode(0000))
	clonePath := t.TempDir() + "/clone-repo"
	_, err := fetcher.Fetch("/path/does/not/exist", "main", clonePath)
	if err == nil {
		t.Fatalf("Expected error for non-existent repo, got nil")
	}
}

func TestFetch_NonExistentBranch(t *testing.T) {
	fetcher := NewFetcher(os.FileMode(0000))
	remoteRepoPath := testutil.SetupRemoteRepo(t)
	clonePath := t.TempDir() + "/clone-repo"

	_, err := fetcher.Fetch(remoteRepoPath, "non-existent-branch", clonePath)
	if err == nil {
		t.Fatalf("Expected error for non-existent branch, got nil")
	}
}

func TestFetch_WithAddPermissions(t *testing.T) {
	fetcher := NewFetcher(os.FileMode(0755))
	remoteRepoPath := testutil.SetupRemoteRepo(t)

	testutil.AddCommitToRepo(t, remoteRepoPath, "README.md", "dummy readme")
	clonePath := t.TempDir() + "/clone-repo"

	_, err := fetcher.Fetch(remoteRepoPath, "main", clonePath)
	if err != nil {
		t.Fatalf("Fetch failed: %v", err)
	}

	// Check file permissions
	fileInfo, err := os.Stat(clonePath + "/README.md")
	if err != nil {
		t.Fatalf("Failed to get file info: %v", err)
	}

	expectedPerm := os.FileMode(0755)
	if fileInfo.Mode().Perm() != expectedPerm {
		t.Errorf("Expected file permissions %v, got %v", expectedPerm, fileInfo.Mode().Perm())
	}
}
