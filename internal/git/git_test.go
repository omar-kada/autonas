package git

import (
	"omar-kada/autonas/internal/testutil"
	"os"
	"testing"

	git "github.com/go-git/go-git/v5"
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

func TestSyncCode_HappyPath(t *testing.T) {
	remoteRepoPath := testutil.SetupRemoteRepo(t)

	testutil.AddCommitToRepo(t, remoteRepoPath, "README.md", "dummy readme")
	clonePath := t.TempDir() + "/clone-repo"

	err := SyncCode(remoteRepoPath, "main", clonePath)
	if err != nil {
		t.Fatalf("SyncCode failed: %v", err)
	}

	assertFileContent(t, clonePath+"/README.md", "dummy readme")
	assertBranch(t, clonePath, "main")

	testutil.AddCommitToRepo(t, remoteRepoPath, "NEWFILE.txt", "new file content")

	err = SyncCode(remoteRepoPath, "main", clonePath)
	if err != nil {
		t.Fatalf("SyncCode failed: %v", err)
	}

	assertFileContent(t, clonePath+"/NEWFILE.txt", "new file content")

}

func TestSyncCode_NoChanges(t *testing.T) {
	remoteRepoPath := testutil.SetupRemoteRepo(t)
	clonePath := t.TempDir() + "/clone-repo"

	err := SyncCode(remoteRepoPath, "main", clonePath)
	if err != nil {
		t.Fatalf("SyncCode failed: %v", err)
	}

	err = SyncCode(remoteRepoPath, "main", clonePath)
	if err != NoErrAlreadyUpToDate {
		t.Fatalf("Expected NoErrAlreadyUpToDate, got: %v", err)
	}
}

func TestSyncCode_NonExistentRepo(t *testing.T) {
	clonePath := t.TempDir() + "/clone-repo"
	err := SyncCode("/path/does/not/exist", "main", clonePath)
	if err == nil {
		t.Fatalf("Expected error for non-existent repo, got nil")
	}
}

func TestSyncCode_NonExistentBranch(t *testing.T) {
	remoteRepoPath := testutil.SetupRemoteRepo(t)
	clonePath := t.TempDir() + "/clone-repo"

	err := SyncCode(remoteRepoPath, "non-existent-branch", clonePath)
	if err == nil {
		t.Fatalf("Expected error for non-existent branch, got nil")
	}
}
