// Package containers handles the deployment and management of services.
package containers

import (
	"fmt"
	"omar-kada/autonas/internal/config"
	"omar-kada/autonas/internal/containers/docker"
	"omar-kada/autonas/internal/containers/model"
	"omar-kada/autonas/internal/files"
	"omar-kada/autonas/internal/logger"
	"os"
	"path/filepath"
	"slices"
)

// Deployer abstracts service deployment operations
type Deployer interface {
	DeployServices(configDir, servicesDir string, currentCfg, cfg config.Config) error
	AddPermission(perm os.FileMode)
}

// NewDockerDeployer creates a new deployer that uses docker for containers
func NewDockerDeployer(log logger.Logger) Deployer {
	return NewDeployer(docker.New(log), log)
}

// NewDeployer creates a new Deployer instance
func NewDeployer(containersManager model.Manager, log logger.Logger) Deployer {
	return &deployer{
		log:               log,
		containersManager: containersManager,
		copyer:            files.NewCopier(),
	}
}

// deployer is responsible for deploying the services
type deployer struct {
	log               logger.Logger
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

	d.log.Debugf("copying files from %s to %s", configDir+"/services", servicesDir)

	for _, service := range cfg.EnabledServices {
		src := filepath.Join(configDir, "services", service)
		dst := filepath.Join(servicesDir, service)
		if err := d.copyer.CopyWithAddPerm(src, dst, d.addPerm); err != nil {
			return fmt.Errorf("error while copying service "+service+" %w", err)
		}
	}

	d.log.Debugf("deploying enabled services: %v\n", cfg.EnabledServices)
	if err := d.containersManager.DeployServices(cfg, servicesDir); err != nil {
		return err
	}
	return nil
}

func (d *deployer) AddPermission(perm os.FileMode) {
	d.addPerm = perm
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
