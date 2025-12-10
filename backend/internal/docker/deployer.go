// Package docker implements operations for docker containers
package docker

import (
	"context"
	"fmt"
	"omar-kada/autonas/internal/events"
	"omar-kada/autonas/internal/files"
	"omar-kada/autonas/internal/process"
	"omar-kada/autonas/internal/shell"
	"omar-kada/autonas/models"
	"path/filepath"
	"slices"
)

// NewDeployer creates an instance of Manager for docker containers
func NewDeployer(dispatcher events.Dispatcher) process.Deployer {
	return &Deployer{
		cmdRunner:    shell.NewRunner(),
		envGenerator: NewEnvGenerator(),
		copier:       files.NewCopier(),
		dispatcher:   dispatcher,
		ctx:          context.Background(),
	}
}

// Deployer manages Docker Compose services.
type Deployer struct {
	cmdRunner    shell.Runner
	envGenerator *EnvGenerator
	copier       files.Copier
	dispatcher   events.Dispatcher
	ctx          context.Context
}

// WithCtx sets the logger for the Deployer
func (d Deployer) WithCtx(ctx context.Context) process.Deployer {
	newDeployer := d
	newDeployer.ctx = ctx
	return newDeployer
}

// RemoveServices stops and removes Docker Compose services.
func (d Deployer) RemoveServices(services []string, servicesDir string) map[string]error {
	d.dispatcher.Debug(d.ctx, "these services will be removed if running.", "services", services)
	errors := make(map[string]error)
	for _, service := range services {
		err := d.composeDown(filepath.Join(servicesDir, service))
		if err != nil {
			d.dispatcher.Error(d.ctx, "Error running docker compose down for %s: %v", service, err)
			errors[service] = err
		}
	}

	return errors
}

// DeployServices generates .env files and runs Docker Compose for enabled services.
func (d Deployer) DeployServices(cfg models.Config, servicesDir string) map[string]error {
	enabledServices := cfg.GetEnabledServices()
	if len(enabledServices) == 0 {
		d.dispatcher.Warn(d.ctx, "No enabled services specified in config. Skipping .env generation and compose up.")
		return nil
	}

	errors := make(map[string]error)
	for _, service := range enabledServices {
		if err := d.envGenerator.generateEnvFile(cfg, servicesDir, service); err != nil {
			d.dispatcher.Error(d.ctx, "Error creating env file", "service", service, "error", err)
			errors[service] = err
			continue
		}
		if err := d.composeUp(filepath.Join(servicesDir, service)); err != nil {
			d.dispatcher.Error(d.ctx, "Error running docker compose", "service", service, "error", err)
			errors[service] = err
		}
	}
	return errors
}

func (d Deployer) composeUp(composePath string) error {
	args := []string{"compose", "--project-directory", composePath, "up", "-d"}
	if err := d.cmdRunner.Run("docker", args...); err != nil {
		return fmt.Errorf("failed to run docker compose up : %w", err)
	}
	return nil
}

func (d Deployer) composeDown(composePath string) error {
	args := []string{"compose", "--project-directory", composePath, "down"}
	if err := d.cmdRunner.Run("docker", args...); err != nil {
		return fmt.Errorf("failed to run docker compose down : %w", err)
	}
	return nil
}

// RemoveAndDeployStacks handles the deployment/removal of services based on the current and new configuration.
func (d Deployer) RemoveAndDeployStacks(oldCfg, cfg models.Config, params models.DeploymentParams) error {
	toBeRemoved := getUnusedServices(oldCfg, cfg)
	if len(toBeRemoved) > 0 {
		// TODO : check if the stack is up before calling RemoveServices
		if errs := d.RemoveServices(toBeRemoved, params.ServicesDir); len(errs) > 0 {
			return fmt.Errorf("error while removing services : %v", errs)
		}
	}

	d.dispatcher.Debug(d.ctx, "copying stacks config files", "src", params.WorkingDir+"/services", "dst", params.ServicesDir)

	enabledServiecs := cfg.GetEnabledServices()

	if errs := d.copyServicesFiles(enabledServiecs, params); len(errs) > 0 {
		return fmt.Errorf("error while copying services files : %v", errs)
	}

	d.dispatcher.Debug(d.ctx, "deploying enabled services", "services", enabledServiecs)
	if errs := d.DeployServices(cfg, params.ServicesDir); len(errs) > 0 {
		return fmt.Errorf("error while removing services : %v", errs)
	}
	return nil
}

func (d Deployer) copyServicesFiles(enabledServiecs []string, params models.DeploymentParams) map[string]error {
	errors := make(map[string]error)
	for _, service := range enabledServiecs {
		src := filepath.Join(params.WorkingDir, "services", service)
		dst := filepath.Join(params.ServicesDir, service)
		if err := d.copier.CopyWithAddPerm(src, dst, params.GetAddWritePerm()); err != nil {
			errors[service] = err
		}
	}
	return errors
}

func getUnusedServices(oldCfg, cfg models.Config) []string {
	var unusedServices []string
	currentlyEnabled := oldCfg.GetEnabledServices()
	shouldBeEnabled := cfg.GetEnabledServices()
	for _, serviceName := range currentlyEnabled {
		if !slices.Contains(shouldBeEnabled, serviceName) {
			unusedServices = append(unusedServices, serviceName)
		}
	}
	return unusedServices
}
