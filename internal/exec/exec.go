package exec

import (
	"fmt"
	"os"
	"path/filepath"
	"omar-kada/autonas/internal/exec/files"
	"omar-kada/autonas/internal/exec/docker"
	"omar-kada/autonas/internal/config"
)

func DeployServices(cfg config.Config) error {
	servicesPath, err := files.CopyServicesToPath(cfg.SERVICES_PATH)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}

	if len(cfg.EnabledServices) == 0 {
		fmt.Fprintln(os.Stderr, "No enabled_services specified in config. Skipping .env generation and compose up.")
		return nil
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
