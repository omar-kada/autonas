package docker

import (
	"fmt"
	"os"
	"omar-kada/autonas/internal/util"
)

func ComposeUp(composePath string) error {
	fmt.Printf("Running: docker compose --project-directory %s up -d\n", composePath)
		cmdStr := fmt.Sprintf("docker compose --project-directory %s up -d", composePath)
		// TODO : replace shell cmd with docker client lib
		if err := util.RunShellCommand(cmdStr); err != nil {
			fmt.Fprintf(os.Stderr, "Error running docker compose for %s: %v\n", composePath, err)
			return err;
		}
		return nil
	}