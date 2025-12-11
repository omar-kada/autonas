// Package git provides functionality to operate on Git repositories.
package git

import (
	"log/slog"
	"omar-kada/autonas/internal/events"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
)

// NoErrAlreadyUpToDate is returned when the repository is already up to date.
var NoErrAlreadyUpToDate = git.NoErrAlreadyUpToDate

// Patch represent the difference between two commits
type Patch struct {
	Diff   string
	Title  string
	Author string
}

// Fetcher is responsible for syncing files from repo
type Fetcher interface {
	Fetch(repoURL, branch, path string) (Patch, error)
	ReFetch(repoURL, branch, path string) (Patch, error)
}

// Syncer is responsible for syncing files from repo
type fetcher struct {
	addPermissions os.FileMode
}

// NewFetcher creates a new Syncer and returns it
func NewFetcher(addPermissions os.FileMode) Fetcher {
	return fetcher{
		addPermissions: addPermissions,
	}
}

// Fetch clones or updates a Git repository at the specified path,
// checking out the specified branch.
// returns NoErrAlreadyUpToDate if the repository is already up to date.
func (f fetcher) Fetch(repoURL, branch, path string) (patch Patch, err error) {
	if branch == "" {
		branch = "main"
	}

	if repoExists(path) {
		patch, err = fetchAndPull(path, branch)
	} else {
		_, err = git.PlainClone(path, &git.CloneOptions{
			URL:           repoURL,
			ReferenceName: plumbing.NewBranchReferenceName(branch),
			SingleBranch:  true,
			Progress:      events.NewSlogWriter(slog.LevelInfo),
		})
	}
	if err == nil && f.addPermissions != os.FileMode(0000) {
		err = f.addPerm(path)
	}
	return patch, err
}

func repoExists(path string) bool {
	_, e := os.Stat(filepath.Join(path, ".git"))
	return e == nil
}

func (f fetcher) addPerm(path string) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		originalPerm := info.Mode().Perm()
		if err := os.Chmod(path, originalPerm|f.addPermissions); err != nil {
			return err
		}
		return nil
	})
}

func (f fetcher) ReFetch(repoURL, branch, path string) (Patch, error) {
	if err := os.RemoveAll(path); err != nil {
		return Patch{}, err
	}
	return f.Fetch(repoURL, branch, path)
}

func fetchAndPull(path string, branch string) (Patch, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return Patch{}, err
	}

	err = repo.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		Progress:   events.NewSlogWriter(slog.LevelInfo),
		Force:      true,
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return Patch{}, err
	}

	patch, err := getPatch(repo, branch)
	if err != nil {
		slog.Error("error while displaying patch : " + err.Error())
	}
	// Checkout branch
	wt, err := repo.Worktree()
	if err != nil {
		return Patch{}, err
	}

	err = wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branch),
		Force:  true, // like `--force`
	})
	if err != nil {
		return Patch{}, err
	}

	// Pull (with force)
	err = wt.Pull(&git.PullOptions{
		RemoteName: "origin",
		Force:      true,
	})

	return patch, err
}

func getPatch(repo *git.Repository, branch string) (Patch, error) {
	// Get local HEAD commit
	headRef, err := repo.Head()
	if err != nil {
		return Patch{}, err
	}

	localCommit, err := repo.CommitObject(headRef.Hash())
	if err != nil {
		return Patch{}, err
	}

	// Get remote HEAD commit (example: origin/main)
	remoteRefName := plumbing.NewRemoteReferenceName("origin", branch)

	remoteRef, err := repo.Reference(remoteRefName, true)
	if err != nil {
		return Patch{}, err
	}

	remoteCommit, err := repo.CommitObject(remoteRef.Hash())
	if err != nil {
		return Patch{}, err
	}

	// Extract trees for diff
	localTree, err := localCommit.Tree()
	if err != nil {
		return Patch{}, err
	}

	remoteTree, err := remoteCommit.Tree()
	if err != nil {
		return Patch{}, err
	}

	// Compute patch (the diff)
	patch, err := localTree.Patch(remoteTree)
	if err != nil {
		return Patch{}, err
	}
	slog.Info("displaying patch between current branch and repo", "patch", patch.String())
	return Patch{
		Title:  remoteCommit.Message,
		Diff:   patch.String(),
		Author: remoteCommit.Author.Name,
	}, nil
}
