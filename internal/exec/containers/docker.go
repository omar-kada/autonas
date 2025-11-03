package containers

import (
	"context"
	"fmt"
	"omar-kada/autonas/internal/config"
	"omar-kada/autonas/internal/exec/files"
	"omar-kada/autonas/internal/shell"
	"os"
	"path/filepath"
	"strings"

	"github.com/moby/moby/client"
)

func newDockerHandler() *dockerHandler {
	return &dockerHandler{
		_writeToFileFunc: files.WriteToFile,
		_runCommandFunc:  shell.RunCommand,
	}
}

// dockerHandler manages Docker Compose services.
type dockerHandler struct {
	_writeToFileFunc func(filePath string, content string) error
	_runCommandFunc  func(cmdAndArgs ...string) error
}

// RemoveServices stops and removes Docker Compose services.
func (d *dockerHandler) RemoveServices(services []string, servicesPath string) error {

	fmt.Printf("Services %s will be removed if running.\n", services)
	for _, serviceName := range services {
		err := d.composeDown(filepath.Join(servicesPath, serviceName))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running docker compose down for %s: %v\n", serviceName, err)
		}
	}

	// TODO : return aggregated error instead of nil
	return nil
}

// DeployServices generates .env files and runs Docker Compose for enabled services.
func (d *dockerHandler) DeployServices(cfg config.Config) error {

	if len(cfg.EnabledServices) == 0 {
		fmt.Fprintln(os.Stderr, "No enabled_services specified in config. Skipping .env generation and compose up.")
		return nil
	}

	for _, service := range cfg.EnabledServices {
		if err := d.generateEnvFile(cfg, service); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating env file for %s: %v\n", service, err)
		}
		if err := d.composeUp(filepath.Join(cfg.ServicesPath, service)); err != nil {
			fmt.Fprintf(os.Stderr, "Error running docker compose for %s: %v\n", service, err)
		}
	}
	// TODO : return aggregated error instead of nil
	return nil
}

func (d *dockerHandler) composeUp(composePath string) error {
	cmd := []string{"docker", "compose", "--project-directory", composePath, "up", "-d"}
	if err := d._runCommandFunc(cmd...); err != nil {
		return fmt.Errorf("failed to run docker compose up : %w", err)
	}
	return nil
}

func (d *dockerHandler) composeDown(composePath string) error {
	cmd := []string{"docker", "compose", "--project-directory", composePath, "down"}
	if err := d._runCommandFunc(cmd...); err != nil {
		return fmt.Errorf("failed to run docker compose down : %w", err)
	}
	return nil
}

func (d *dockerHandler) generateEnvFile(cfg config.Config, service string) error {
	serviceCfg := cfg.PerService(service)

	var content strings.Builder
	for k, v := range serviceCfg {
		content.WriteString(fmt.Sprintf("%s=%v\n", k, v))
	}

	envFilePath := filepath.Join(cfg.ServicesPath, service, ".env")
	return d._writeToFileFunc(envFilePath, content.String())
}

// GetManagedContainers returns the list of containers (as returned by ContainerList)
// that are managed by AutoNAS
func (d *dockerHandler) GetManagedContainers() (map[string][]Summary, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}
	defer cli.Close()

	summaries, err := cli.ContainerList(ctx, client.ContainerListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	matches := make(map[string][]Summary)
	for _, c := range summaries {
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
				matches[serviceName] = append(matches[serviceName], Summary{
					ID:     c.ID,
					Names:  c.Names,
					Image:  c.Image,
					State:  c.State,
					Status: c.Status,
				})
			}
		}
	}
	return matches, nil
}
