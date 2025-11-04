// Package containers handles the deployment and management of services.
package containers

import (
	"omar-kada/autonas/internal/config"
	"omar-kada/autonas/internal/containers/docker"
	"omar-kada/autonas/internal/containers/model"
	"slices"

	copydir "github.com/otiai10/copy"
)

// NewDockerDeployer creates a new deployer that uses docker for containers
func NewDockerDeployer() Deployer {
	return NewDeployer(docker.New())
}

// NewDeployer creates a new Deployer instance
func NewDeployer(containersManager model.Manager) Deployer {
	return Deployer{
		containersManager: containersManager,
		_copyFunc:         copydir.Copy,
	}
}

type Deployer struct {
	containersManager model.Manager
	_copyFunc         func(srcFolder, servicesPath string, _ ...copydir.Options) error
}

// DeployServices handles the deployment/removal of services based on the current and new configuration.
// It accepts a ServiceManager to allow injection in tests; callers can pass DefaultServices.
func (d *Deployer) DeployServices(configFolder string, currentCfg, cfg config.Config) error {
	toBeRemoved := getUnusedServices(currentCfg, cfg)
	if err := d.containersManager.RemoveServices(toBeRemoved, currentCfg.ServicesPath); err != nil {
		return err
	}

	if err := d._copyFunc(configFolder+"/services", cfg.ServicesPath); err != nil {
		return err
	}

	if err := d.containersManager.DeployServices(cfg); err != nil {
		return err
	}
	return nil
}

func getUnusedServices(currentCfg, cfg config.Config) []string {
	var unusedServices []string
	for _, serviceName := range currentCfg.EnabledServices {
		if !slices.Contains(cfg.EnabledServices, serviceName) {
			unusedServices = append(unusedServices, serviceName)
		}
	}
	return unusedServices
}
