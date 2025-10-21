package cli

import (
	"fmt"
	"omar-kada/autonas/internal/config"
	"omar-kada/autonas/internal/exec"
	"os"
	"omar-kada/autonas/internal/exec/git"
	"github.com/robfig/cron/v3"
)

func RunCmd(configFiles []string, configRepo string) error {
	err := git.SyncCode(configRepo, "main", ".")
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
	exec.DeployServices(cfg)
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
