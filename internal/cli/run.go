package cli

import (
	"fmt"
	"os"
	"omar-kada/autonas/internal/config"
	"omar-kada/autonas/internal/exec"
	"omar-kada/autonas/internal/exec/git"
	"github.com/robfig/cron/v3"
)

func RunCmd(configFiles []string, configRepo string) error {
	currentCfg := config.GetCurrentConfig()

	// TODO : add these to configuration
	configFolder := "."
	repoBranch := "main"
	
	err := git.SyncCode(configRepo, repoBranch, configFolder)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting config repo: %v\n", err)
		return err
	}
	
	cfg, err := config.LoadConfig(configFiles)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		return err
	}
	fmt.Printf("Final consolidated config: %+v\n", cfg)
	
	// Copy all files from ./services to SERVICES_PATH
	err = exec.DeployServices(configFolder, currentCfg, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error deploying services: %v\n", err)
		return err
	}
	return nil
}

func RunPeriocically(cronPeriod string, configFiles []string, configRepo string) {
	c := cron.New()

	c.AddFunc(cronPeriod, func() {
		RunCmd(configFiles, configRepo)
	})

	c.Start()
	select {} // keep running
}
