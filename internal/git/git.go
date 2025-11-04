// Package git provides functionality to operate on Git repositories.
package git

import (
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

var (
	// NoErrAlreadyUpToDate is returned when the repository is already up to date.
	NoErrAlreadyUpToDate = git.NoErrAlreadyUpToDate
)

// SyncCode clones or updates a Git repository at the specified path,
// checking out the specified branch.
// returns NoErrAlreadyUpToDate if the repository is already up to date.
func SyncCode(repoURL, branch, path string) error {

	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:           repoURL,
		ReferenceName: plumbing.NewBranchReferenceName(branch),
		SingleBranch:  true,
		Progress:      os.Stdout,
	})
	if err == git.ErrRepositoryAlreadyExists {
		return fetchAndPull(path, branch)
	}
	return err
}

func fetchAndPull(path string, branch string) error {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return err
	}

	err = repo.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		Progress:   os.Stdout,
		Force:      true,
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return err
	}

	// Checkout branch
	wt, _ := repo.Worktree()
	err = wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branch),
		Force:  true, // like `--force`
	})
	if err != nil {
		return err
	}

	// Pull (with force)
	err = wt.Pull(&git.PullOptions{
		RemoteName: "origin",
		Force:      true,
	})

	return err
}
