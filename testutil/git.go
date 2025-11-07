package testutil

import (
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func SetupRemoteRepo(t *testing.T) string {
	t.Helper()

	dir := t.TempDir() + "/remote-repo"
	_, err := git.PlainInitWithOptions(dir, &git.PlainInitOptions{
		InitOptions: git.InitOptions{
			DefaultBranch: plumbing.NewBranchReferenceName("main"),
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	AddCommitToRepo(t, dir, "README.md", "initial commit")
	return dir
}

func AddCommitToRepo(t *testing.T, repoPath string, fileName string, content string) {
	t.Helper()
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		t.Fatalf("Failed to open repo: %v", err)
	}

	w, err := r.Worktree()
	if err != nil {
		t.Fatalf("Failed to get worktree: %v", err)
	}

	file, err := w.Filesystem.Create(fileName)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	file.Write([]byte(content))
	file.Close()

	w.Add(fileName)
	_, err = w.Commit("add "+fileName, &git.CommitOptions{
		Author: &object.Signature{Name: "Test", Email: "test@test.com"},
	})
	if err != nil {
		t.Fatalf("Failed to commit changes: %v", err)
	}
}
