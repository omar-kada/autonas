package docker

import (
	"context"
	"fmt"
	"os"
	"strings"

	"omar-kada/autonas/internal/util"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

func ComposeUp(composePath string) error {
	cmdStr := fmt.Sprintf("docker compose --project-directory %s up -d", composePath)
	fmt.Printf("Running: %s \n", cmdStr)
	// TODO : replace shell cmd with docker client lib
	if err := util.RunShellCommand(cmdStr); err != nil {
		return fmt.Errorf("failed to run docker compose up : %w", err)
	}
	return nil
}

func ComposeDown(composePath string) error {
	cmdStr := fmt.Sprintf("docker compose --project-directory %s down", composePath)
	fmt.Printf("Running: %s \n", cmdStr)
	// TODO : replace shell cmd with docker client lib
	if err := util.RunShellCommand(cmdStr); err != nil {
		return fmt.Errorf("failed to run docker compose down : %w", err)
	}
	return nil
}

// RemoveContainer removes the container with the given ID.
// If force is true, the container will be removed even if running.
// Volumes attached to the container will also be removed.
func RemoveContainer(containerID string, force bool) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("failed to create docker client: %w", err)
	}
	defer cli.Close()

	opts := client.ContainerRemoveOptions{
		Force:         force,
		RemoveVolumes: true,
	}

	if err := cli.ContainerRemove(ctx, containerID, opts); err != nil {
		return fmt.Errorf("failed to remove container %s: %w", containerID, err)
	}
	return nil
}

// GetManagedContainers returns the list of containers (as returned by ContainerList)
// that are managed by AutoNAS
func GetManagedContainers() (map[string][]container.Summary, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}
	defer cli.Close()

	containers, err := cli.ContainerList(ctx, client.ContainerListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	matches := make(map[string][]container.Summary)
	for _, c := range containers {
		inspect, err := cli.ContainerInspect(ctx, c.ID)
		if err != nil {
			// don't fail entirely on single-container inspect error; just log and continue
			fmt.Fprintf(os.Stderr, "warning: failed to inspect container %s: %v\n", c.ID, err)
			continue
		}
		var managed bool
		var serviceName string
		for _, env := range inspect.Config.Env {
			if strings.HasPrefix(env, "AUTONAS_MANAGED=") {
				managed = true
			}
			if strings.HasPrefix(env, "AUTONAS_SERVICE_NAME=") {
				serviceName = strings.TrimPrefix(env, "AUTONAS_SERVICE_NAME=")
			}
			if managed && serviceName != "" {
				break
			}
		}
		if managed {
			if serviceName == "" {
				fmt.Fprintf(os.Stderr, "warning: container %s marked as AUTONAS_MANAGED but missing AUTONAS_SERVICE_NAME\n", c.ID)
			} else {
				matches[serviceName] = append(matches[serviceName], c)
			}
		}
	}
	return matches, nil
}
