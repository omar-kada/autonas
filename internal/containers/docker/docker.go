// Package docker implements a manager for docker containers
package docker

import (
	"context"
	"fmt"
	"omar-kada/autonas/internal/config"
	"omar-kada/autonas/internal/containers/model"
	"omar-kada/autonas/internal/files"
	"omar-kada/autonas/internal/logger"
	"omar-kada/autonas/internal/shell"
	"path/filepath"
	"strings"

	"github.com/moby/moby/client"
)

// New creates an instance of Manager for docker containers
func New(log logger.Logger) *Manager {
	return &Manager{
		log:              log,
		_writeToFileFunc: files.WriteToFile,
		_runCommandFunc:  shell.RunCommand,
	}
}

// Manager manages Docker Compose services.
type Manager struct {
	log              logger.Logger
	_writeToFileFunc func(filePath string, content string) error
	_runCommandFunc  func(cmd string, args ...string) error
}

// RemoveServices stops and removes Docker Compose services.
func (d Manager) RemoveServices(services []string, servicesPath string) error {

	d.log.Debugf("services %s will be removed if running.", services)
	for _, serviceName := range services {
		err := d.composeDown(filepath.Join(servicesPath, serviceName))
		if err != nil {
			d.log.Errorf("Error running docker compose down for %s: %v", serviceName, err)
		}
	}

	// TODO : return aggregated error instead of nil
	return nil
}

// DeployServices generates .env files and runs Docker Compose for enabled services.
func (d Manager) DeployServices(cfg config.Config) error {

	if len(cfg.EnabledServices) == 0 {
		d.log.Warnf("No enabled_services specified in config. Skipping .env generation and compose up.")
		return nil
	}

	for _, service := range cfg.EnabledServices {
		if err := d.generateEnvFile(cfg, service); err != nil {
			d.log.Errorf("Error creating env file for %s: %v", service, err)
		}
		if err := d.composeUp(filepath.Join(cfg.ServicesPath, service)); err != nil {
			d.log.Errorf("Error running docker compose for %s: %v", service, err)
		}
	}
	// TODO : return aggregated error instead of nil
	return nil
}

func (d Manager) composeUp(composePath string) error {
	args := []string{"compose", "--project-directory", composePath, "up", "-d"}
	if err := d._runCommandFunc("docker", args...); err != nil {
		return fmt.Errorf("failed to run docker compose up : %w", err)
	}
	return nil
}

func (d Manager) composeDown(composePath string) error {
	args := []string{"compose", "--project-directory", composePath, "down"}
	if err := d._runCommandFunc("docker", args...); err != nil {
		return fmt.Errorf("failed to run docker compose down : %w", err)
	}
	return nil
}

func (d Manager) generateEnvFile(cfg config.Config, service string) error {
	serviceCfg := cfg.PerService(service)

	var content strings.Builder
	for _, v := range serviceCfg {
		content.WriteString(fmt.Sprintf("%s=%v\n", v.Key, v.Value))
	}

	envFilePath := filepath.Join(cfg.ServicesPath, service, ".env")
	return d._writeToFileFunc(envFilePath, content.String())
}

// GetManagedContainers returns the list of containers (as returned by ContainerList)
// that are managed by AutoNAS
func (d Manager) GetManagedContainers() (map[string][]model.Summary, error) {
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

	matches := make(map[string][]model.Summary)
	for _, c := range summaries {
		inspect, err := cli.ContainerInspect(ctx, c.ID)
		if err != nil {
			// don't fail entirely on single-container inspect error; just log and continue
			d.log.Errorf("Failed to inspect container %s, %s: %v", c.ID, c.Names, err)
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
				d.log.Errorf("container %s marked as AUTONAS_MANAGED but missing AUTONAS_SERVICE_NAME", c.ID)
			} else {
				matches[serviceName] = append(matches[serviceName], model.Summary{
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
