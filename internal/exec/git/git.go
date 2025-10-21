package git


import (
	"fmt"
  "os"
	"github.com/go-git/go-git/v5"
  "github.com/go-git/go-git/v5/plumbing"
)


func SyncCode(repoURL, branch, path string) error {

    // Clone
    repo, err := git.PlainClone(path, false, &git.CloneOptions{
        URL:           repoURL,
        ReferenceName: plumbing.NewBranchReferenceName(branch),
        SingleBranch:  true,
        Progress:      os.Stdout,
    })
    if err == git.ErrRepositoryAlreadyExists {
        repo, err = git.PlainOpen(path)
        if err != nil {
			return err
        }
    } else if err != nil {
		return err
    }

    // Fetch
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
    if err != nil && err != git.NoErrAlreadyUpToDate {
		return err
    }

    fmt.Println("Done!")
	return nil
}
