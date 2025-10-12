package cli

import (
	"fmt"
	"omar-kada/autonas/internal/config"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	copydir "github.com/otiai10/copy"
	"github.com/go-git/go-git/v5"
    "github.com/go-git/go-git/v5/plumbing"
)

func RunCmd(configFiles []string, configRepo string) error {
	err := getGit(configRepo, "main", ".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting config repo: %v\n", err)
		return err
	}

	cfg, err := config.LoadConfig(configFiles)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		return err
	}
	fmt.Printf("Final consolidated config: %+v\n", cfg)

	// Copy all files from ./services to SERVICES_PATH
	servicesPath, err := copyServicesToPath(cfg.SERVICES_PATH)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}

	if len(cfg.EnabledServices) == 0 {
		fmt.Fprintln(os.Stderr, "No enabled_services specified in config. Skipping .env generation and compose up.")
		return nil
	}

	// For each enabled service, generate .env and run docker compose up
	for _, service := range cfg.EnabledServices {
		envFilePath := filepath.Join(servicesPath, service, ".env")
		f, err := os.Create(envFilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating .env for %s: %v\n", service, err)
			continue
		}
		// Write global env
		fmt.Fprintf(f, "AUTONAS_MANAGED=true\n")
		fmt.Fprintf(f, "AUTONAS_HOST=%v\n", cfg.AUTONAS_HOST)
		fmt.Fprintf(f, "SERVICES_PATH=%v\n", cfg.SERVICES_PATH)
		fmt.Fprintf(f, "DATA_PATH=%v/%v\n", cfg.DATA_PATH, service)
		for k, v := range cfg.Extra {
			fmt.Fprintf(f, "%s=%v\n", k, v)
		}

		// Write service-specific env
		if svcVars, ok := cfg.Services[service]; ok {
			if svcVars.Port != 0 {
				fmt.Fprintf(f, "PORT=%v\n", svcVars.Port)
			}
			if svcVars.Version != "" {
				fmt.Fprintf(f, "VERSION=%v\n", svcVars.Version)
			}
			for k, v := range svcVars.Extra {
				fmt.Fprintf(f, "%s=%v\n", strings.ToUpper(k), v)
			}
		}
		f.Close()

		// Run docker compose up for the service
		composePath := filepath.Join(servicesPath, service)
		composeFile := filepath.Join(composePath, "compose.yaml")
		fmt.Printf("Running: docker compose -f %s up -d\n", composeFile)
		cmdStr := fmt.Sprintf("docker compose -f %s up -d", composeFile)
		if err := runShellCommand(cmdStr); err != nil {
			fmt.Fprintf(os.Stderr, "Error running docker compose for %s: %v\n", service, err)
		}
	}
	return nil
}

func copyServicesToPath(servicesPath string) (string, error) {
	if servicesPath == "" {
		return "", fmt.Errorf("SERVICES_PATH not set in config. Aborting copy")
	}
	err := copydir.Copy("./services", servicesPath)
	if err != nil {
		return "", fmt.Errorf("error copying services: %v", err)
	}
	fmt.Printf("Copied all files from ./services to %s\n", servicesPath)
	return servicesPath, nil
}

// runShellCommand runs a shell command and returns error if any
func runShellCommand(cmdStr string) error {
	var c *exec.Cmd
	if runtime.GOOS == "windows" {
		c = execCommand("cmd", "/C", cmdStr)
	} else {
		shell := os.Getenv("SHELL")
		if shell == "" {
			shell = "bash"
		}
		c = execCommand(shell, "-c", cmdStr)
	}
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

// execCommand is a wrapper for exec.Command for testability
var execCommand = defaultExecCommand

func defaultExecCommand(name string, arg ...string) *exec.Cmd {
	return exec.Command(name, arg...)
}

func getGit(repoURL, branch, path string) error {

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
