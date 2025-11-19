package docker

import (
	"context"
	"fmt"
	"log/slog"
	"omar-kada/autonas/models"
	"strings"

	"github.com/moby/moby/client"
)

// Client defines the methods from the Docker client that are used by the Inspector
type Client interface {
	ContainerList(ctx context.Context, options client.ContainerListOptions) (client.ContainerListResult, error)
	ContainerInspect(ctx context.Context, containerID string, options client.ContainerInspectOptions) (client.ContainerInspectResult, error)
}

// Inspector implements information retrieval about docker stacks
type Inspector struct {
	dockerClient Client
}

// NewInspector creates new inspector given a docker client
func NewInspector(dockerClient Client) *Inspector {
	return &Inspector{
		dockerClient: dockerClient,
	}
}

// GetManagedContainers returns the list of containers (as returned by ContainerList)
// that are managed by AutoNAS
func (i Inspector) GetManagedContainers(servicesDir string) (map[string][]models.ContainerSummary, error) {
	ctx := context.Background()
	summaries, err := i.dockerClient.ContainerList(ctx, client.ContainerListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	matches := make(map[string][]models.ContainerSummary)
	for _, c := range summaries.Items {
		inspect, err := i.dockerClient.ContainerInspect(ctx, c.ID, client.ContainerInspectOptions{})
		if err != nil {
			// don't fail entirely on single-container inspect error; just log and continue
			slog.Error("Failed to inspect container",
				"containerId", c.ID,
				"names", c.Names,
				"error", err)
			continue
		}

		for key, value := range inspect.Container.Config.Labels {
			if strings.EqualFold(key, "com.docker.compose.project.working_dir") {
				after, found := strings.CutPrefix(value, servicesDir)
				if found {
					serviceName, _ := strings.CutPrefix(after, "/")
					matches[serviceName] = append(matches[serviceName], models.ContainerSummary{
						ID:     c.ID,
						Names:  c.Names,
						Image:  c.Image,
						State:  string(c.State),
						Status: c.Status,
					})
					break
				}
			}
		}
	}
	return matches, nil
}
