// Package containers handles the deployment and management of services.
package containers

import (
	"fmt"
	"io/fs"
	"omar-kada/autonas/internal/config"
	"omar-kada/autonas/internal/containers/docker"
	"omar-kada/autonas/internal/containers/model"
	"omar-kada/autonas/internal/logger"
	"os"
	"path/filepath"
	"slices"

	copydir "github.com/otiai10/copy"
)

// NewDockerDeployer creates a new deployer that uses docker for containers
func NewDockerDeployer(log logger.Logger) Deployer {
	return NewDeployer(docker.New(log), log)
}

// NewDeployer creates a new Deployer instance
func NewDeployer(containersManager model.Manager, log logger.Logger) Deployer {
	return Deployer{
		log:               log,
		containersManager: containersManager,

		_copyFunc: copydir.Copy,
	}
}

// Deployer is responsible for deploying the services
type Deployer struct {
	log               logger.Logger
	containersManager model.Manager

	_copyFunc func(srcFolder, servicesPath string, _ ...copydir.Options) error
}

// DeployServices handles the deployment/removal of services based on the current and new configuration.
// It accepts a ServiceManager to allow injection in tests; callers can pass DefaultServices.
func (d *Deployer) DeployServices(configFolder string, currentCfg, cfg config.Config) error {
	toBeRemoved := getUnusedServices(currentCfg, cfg)
	if err := d.containersManager.RemoveServices(toBeRemoved, currentCfg.ServicesPath); err != nil {
		return err
	}

	d.log.Debugf("copying files from %s to %s", configFolder+"/services", cfg.ServicesPath)

	for _, service := range cfg.EnabledServices {
		src := filepath.Join(configFolder, "services", service)
		dst := filepath.Join(cfg.ServicesPath, service)
		if err := d._copyFunc(src, dst); err != nil {
			return fmt.Errorf("error while copying service "+service+" %w", err)
		}
	}

	if os.Getenv("ENV") == "DEV" {
		// allow auto removing of copied services while testing
		err := filepath.WalkDir(cfg.ServicesPath, func(path string, _ fs.DirEntry, _ error) error {
			err := os.Chmod(path, 0777)
			return err
		})
		if err != nil {
			return err
		}
	}

	d.log.Debugf("deploying enabled services: %v\n", cfg.EnabledServices)
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
