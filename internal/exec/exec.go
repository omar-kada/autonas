// Package exec handles the deployment and management of services.
package exec

import (
	"omar-kada/autonas/internal/config"
	"omar-kada/autonas/internal/exec/containers"
	"slices"

	copydir "github.com/otiai10/copy"
)

// Deployer abstracts service deployment operations
type Deployer interface {
	DeployServices(configFolder string, currentCfg, cfg config.Config) error
}

// New creates a new Deployer instance with default dependencies.
func New() Deployer {
	return &defaultDeployer{
		containersHandler: containers.New(),
		_copyFunc:         copydir.Copy,
	}
}

type defaultDeployer struct {
	containersHandler containers.Handler
	_copyFunc         func(srcFolder, servicesPath string, _ ...copydir.Options) error
}

// DeployServices handles the deployment/removal of services based on the current and new configuration.
// It accepts a ServiceManager to allow injection in tests; callers can pass DefaultServices.
func (d *defaultDeployer) DeployServices(configFolder string, currentCfg, cfg config.Config) error {
	toBeRemoved := getUnusedServices(currentCfg, cfg)
	if err := d.containersHandler.RemoveServices(toBeRemoved, currentCfg.ServicesPath); err != nil {
		return err
	}

	if err := d._copyFunc(configFolder+"/services", cfg.ServicesPath); err != nil {
		return err
	}

	if err := d.containersHandler.DeployServices(cfg); err != nil {
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
