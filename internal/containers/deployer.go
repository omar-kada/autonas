// Package containers handles the deployment and management of services.
package containers

import (
	"fmt"
	"log/slog"
	"omar-kada/autonas/internal/config"
	"omar-kada/autonas/internal/containers/docker"
	"omar-kada/autonas/internal/containers/model"
	"omar-kada/autonas/internal/files"
	"os"
	"path/filepath"
	"slices"
)

// Deployer abstracts service deployment operations
type Deployer interface {
	DeployServices(configDir, servicesDir string, currentCfg, cfg config.Config) error
	WithPermission(perm os.FileMode) Deployer
}

// NewDockerDeployer creates a new deployer that uses docker for containers
func NewDockerDeployer() Deployer {
	return newDeployer(docker.New())
}

// newDeployer creates a new Deployer instance
func newDeployer(containersManager model.Manager) *deployer {
	return &deployer{
		containersManager: containersManager,
		copyer:            files.NewCopier(),
	}
}

// deployer is responsible for deploying the services
type deployer struct {
	containersManager model.Manager
	copyer            files.Copier
	addPerm           os.FileMode
}

// DeployServices handles the deployment/removal of services based on the current and new configuration.
// It accepts a ServiceManager to allow injection in tests; callers can pass DefaultServices.
func (d *deployer) DeployServices(configDir, servicesDir string, currentCfg, cfg config.Config) error {
	toBeRemoved := getUnusedServices(currentCfg, cfg)
	// TODO : check if the stack is up before calling RemoveServices
	if err := d.containersManager.RemoveServices(toBeRemoved, servicesDir); err != nil {
		return err
	}

	slog.Debug("copying files from src to dst", "src", configDir+"/services", "dst", servicesDir)

	enabledServiecs := cfg.GetEnabledServices()
	for _, service := range enabledServiecs {
		src := filepath.Join(configDir, "services", service)
		dst := filepath.Join(servicesDir, service)
		if err := d.copyer.CopyWithAddPerm(src, dst, d.addPerm); err != nil {
			return fmt.Errorf("error while copying service "+service+" %w", err)
		}
	}

	slog.Debug("deploying enabled services", "services", enabledServiecs)
	if err := d.containersManager.DeployServices(cfg, servicesDir); err != nil {
		return err
	}
	return nil
}

// WithPermission adds permission to created files by the deployer
func (d *deployer) WithPermission(perm os.FileMode) Deployer {
	deployer := newDeployer(d.containersManager)
	deployer.addPerm = perm
	return deployer
}

func getUnusedServices(currentCfg, cfg config.Config) []string {
	var unusedServices []string
	currentlyEnabled := currentCfg.GetEnabledServices()
	shouldBeEnabled := cfg.GetEnabledServices()
	for _, serviceName := range currentlyEnabled {
		if !slices.Contains(shouldBeEnabled, serviceName) {
			unusedServices = append(unusedServices, serviceName)
		}
	}
	return unusedServices
}
