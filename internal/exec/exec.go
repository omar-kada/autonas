package exec

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"omar-kada/autonas/internal/exec/files"
	"omar-kada/autonas/internal/exec/docker"
	"omar-kada/autonas/internal/config"
)

func DeployServices(configFolder string, currentCfg, cfg config.Config) error {
	
	if err := removeUnusedServices(currentCfg, cfg); err != nil {
		return err
	}

	if err := deployActivatedServices(configFolder, cfg); err != nil {
		return err
	}
	return nil
}

func removeUnusedServices(currentCfg, cfg config.Config) error {
	if len(currentCfg.EnabledServices) == 0 {
		fmt.Println("No previous services found. Skipping removal of unused services.")
		return nil
	}

	for _, serviceName := range currentCfg.EnabledServices {
		if !slices.Contains(cfg.EnabledServices, serviceName) {
			fmt.Printf("Service %s was previously enabled but is no longer in the config. It will be removed if running.\n", serviceName)
			err := docker.ComposeDown(filepath.Join(currentCfg.SERVICES_PATH, serviceName))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error running docker compose down for %s: %v\n", serviceName, err)
			}
		}
	}

	return nil
}

func deployActivatedServices(configFolder string, cfg config.Config) error {
	
	if len(cfg.EnabledServices) == 0 {
		fmt.Fprintln(os.Stderr, "No enabled_services specified in config. Skipping .env generation and compose up.")
		return nil
	}
	
	servicesPath, err := files.CopyServicesToPath(configFolder+"/services", cfg.SERVICES_PATH)
	if err != nil {
		return err
	}

	// For each enabled service, generate .env and run docker compose up
	for _, service := range cfg.EnabledServices {

		serviceCfg := config.ConfigPerService(cfg, service)
		err := files.WriteEnvFile(filepath.Join(servicesPath, service, ".env"), serviceCfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating env file for %s: %v\n", service, err)
		} else {
			fmt.Printf("Generated .env for service %s with config: %+v\n", service, serviceCfg)
		}
		// Run docker compose up for the service
		err = docker.ComposeUp(filepath.Join(servicesPath, service))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running docker compose for %s: %v\n", service, err)
		}
	}
	return nil
}
