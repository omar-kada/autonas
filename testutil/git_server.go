package testutil

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// GitTestServer represents a test Git HTTP server
type GitTestServer struct {
	Server  *httptest.Server
	URL     string
	TempDir string
	Repo    *git.Repository
}

// NewGitTestServer creates and starts a test Git server
func NewGitTestServer() (*GitTestServer, error) {
	tempDir, err := os.MkdirTemp("", "git-test-*")
	if err != nil {
		return nil, err
	}

	// Initialize Git repo
	repo, err := git.PlainInit(tempDir, false)
	if err != nil {
		return nil, err
	}

	// Create test server
	handler := createGitHandler(tempDir, repo)
	server := httptest.NewServer(handler)

	return &GitTestServer{
		Server:  server,
		URL:     server.URL,
		TempDir: tempDir,
		Repo:    repo,
	}, nil
}

// Close shuts down the server and cleans up
func (g *GitTestServer) Close() {
	g.Server.Close()
	os.RemoveAll(g.TempDir)
}

// AddFile creates a file and commits it to the test repo
func (g *GitTestServer) AddFile(path, content string) error {
	fullPath := filepath.Join(g.TempDir, path)
	err := os.WriteFile(fullPath, []byte(content), 0644)
	if err != nil {
		return err
	}

	worktree, err := g.Repo.Worktree()
	if err != nil {
		return err
	}

	_, err = worktree.Add(path)
	if err != nil {
		return err
	}

	_, err = worktree.Commit("Add "+path, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test",
			Email: "test@example.com",
			When:  time.Now(),
		},
	})
	return err
}

// GetGitURL returns the Git HTTP URL for clients
func (g *GitTestServer) GetGitURL() string {
	return g.URL + "/repo.git"
}

func createGitHandler(repoPath string, repo *git.Repository) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/repo.git/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(repoPath, ".git"))
	})

	return mux
}
