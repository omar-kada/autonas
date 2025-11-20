package docker

import (
	"context"
	"fmt"
	"log/slog"
	"omar-kada/autonas/models"
	"strings"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

// Client defines the methods from the Docker client that are used by the Inspector
type Client interface {
	ContainerList(ctx context.Context, options client.ContainerListOptions) (client.ContainerListResult, error)
	ContainerInspect(ctx context.Context, containerID string, options client.ContainerInspectOptions) (client.ContainerInspectResult, error)
}

// Inspector implements information retrieval about docker stacks
type Inspector struct{}

// NewInspector creates new inspector given a docker client
func NewInspector() *Inspector {
	return &Inspector{}
}

// GetManagedContainers returns the list of containers (as returned by ContainerList)
// that are managed by AutoNAS
func (Inspector) GetManagedContainers(servicesDir string) (map[string][]models.ContainerSummary, error) {
	dockerClient, err := client.New(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}
	return getManagedContainersWithClient(dockerClient, servicesDir)
}

func getManagedContainersWithClient(dockerClient Client, servicesDir string) (map[string][]models.ContainerSummary, error) {
	ctx := context.Background()
	summaries, err := dockerClient.ContainerList(ctx, client.ContainerListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	matches := make(map[string][]models.ContainerSummary)
	for _, c := range summaries.Items {
		serviceName, err := getServiceNameFromLabel(ctx, dockerClient, c, servicesDir)
		if err != nil {
			slog.Error("Failed to inspect container",
				"containerId", c.ID, "names", c.Names, "error", err)
			// don't fail entirely on single-container inspect error; just log and continue
			continue
		}
		if serviceName != "" {
			matches[serviceName] = append(matches[serviceName], models.ContainerSummary{
				ID:     c.ID,
				Names:  c.Names,
				Image:  c.Image,
				State:  string(c.State),
				Status: c.Status,
			})
		}
	}
	return matches, nil
}

func getServiceNameFromLabel(ctx context.Context, dockerClient Client, c container.Summary, servicesDir string) (string, error) {
	inspect, err := dockerClient.ContainerInspect(ctx, c.ID, client.ContainerInspectOptions{})
	if err != nil {
		return "", err
	}
	for key, value := range inspect.Container.Config.Labels {
		if strings.EqualFold(key, "com.docker.compose.project.working_dir") {
			if after, found := strings.CutPrefix(value, servicesDir); found {
				return strings.TrimPrefix(after, "/"), nil
			}
			return "", nil
		}
	}
	return "", nil
}
