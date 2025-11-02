// Package exec handles the deployment and management of services.
package exec

import (
	"omar-kada/autonas/internal/config"
	"omar-kada/autonas/internal/exec/containers"
	"omar-kada/autonas/internal/exec/files"
	"slices"
)

// Deployer abstracts service deployment operations
type Deployer interface {
	DeployServices(configFolder string, currentCfg, cfg config.Config) error
}

var defaultFileManager = files.NewManager()
var defaultContainersHandler = containers.NewHandler(defaultFileManager)

// New creates a new Deployer instance with default dependencies.
func New() Deployer {
	return &defaultDeployer{
		containersHandler: defaultContainersHandler,
		fileManager:       defaultFileManager,
	}
}

type defaultDeployer struct {
	containersHandler containers.Handler
	fileManager       files.Manager
}

// DeployServices handles the deployment/removal of services based on the current and new configuration.
// It accepts a ServiceManager to allow injection in tests; callers can pass DefaultServices.
func (d *defaultDeployer) DeployServices(configFolder string, currentCfg, cfg config.Config) error {
	toBeRemoved := getUnusedServices(currentCfg, cfg)
	if err := d.containersHandler.RemoveServices(toBeRemoved, currentCfg.ServicesPath); err != nil {
		return err
	}

	if err := d.fileManager.CopyToPath(configFolder+"/services", cfg.ServicesPath); err != nil {
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
