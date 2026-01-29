package cli

import (
	"fmt"
	"log/slog"

	"omar-kada/autonas/internal/docker"
	"omar-kada/autonas/internal/events"
	"omar-kada/autonas/internal/git"
	"omar-kada/autonas/internal/process"
	"omar-kada/autonas/internal/server"
	"omar-kada/autonas/internal/shell"
	"omar-kada/autonas/internal/storage"
	"omar-kada/autonas/models"

	"github.com/spf13/cobra"
)

type runCommand struct {
	executor     shell.Executor
	storeCreator func(params RunParams) (storage.Storage, error)

	cmd    *cobra.Command
	params RunParams
}

// NewRunCommand creates a new run
func NewRunCommand(executor shell.Executor, storeCreator func(params RunParams) (storage.Storage, error)) *cobra.Command {
	run := runCommand{
		params:       RunParams{},
		executor:     executor,
		storeCreator: storeCreator,
	}

	run.cmd = &cobra.Command{
		Use:   "run",
		Short: "Run with optional config files",
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := run.doRun(); err != nil {
				slog.Error(err.Error())
				return err
			}
			return nil
		},
	}
	run.cmd.Flags().StringVarP(&run.params.ConfigFile, string(_file), "f", "",
		varInfoMap.GetDefaultString("YAML config file", _file))
	run.cmd.Flags().StringVarP(&run.params.WorkingDir, string(_workingDir), "d", "",
		varInfoMap.GetDefaultString("directory where autonas data will be stored", _workingDir))
	run.cmd.Flags().StringVarP(&run.params.ServicesDir, string(_servicesDir), "s", "",
		varInfoMap.GetDefaultString("directory where services compose stacks will be stored", _servicesDir))
	run.cmd.Flags().StringVarP(&run.params.AddWritePerm, string(_addWritePerm), "w", "",
		varInfoMap.GetDefaultString("when true, the tool adds write permission to config files", _addWritePerm))
	run.cmd.Flags().IntVarP(&run.params.Port, string(_port), "p", 0,
		varInfoMap.GetDefaultString("port that will be used for exposing the API", _port))

	return run.cmd
}

func (run *runCommand) doRun() error {
	params := getParamsWithDefaults(run.params)
	store, err := run.storeCreator(params)
	if err != nil {
		return fmt.Errorf("couldn't init storage %w", err)
	}

	dispatcher := events.NewDefaultDispatcher(store)
	configStore := storage.NewConfigStore(params.ConfigFile)
	scheduler := process.NewConfigScheduler(configStore)
	configStore.SetOnChange(func(oldCfg, cfg models.Config) {
		slog.Debug("checking if cron changed", "oldCron", oldCfg.Settings.CronPeriod, "newCron", cfg.Settings.CronPeriod)
		if oldCfg.Settings.CronPeriod != cfg.Settings.CronPeriod {
			scheduler.ReSchedule()
		}
	})
	inspector, err := docker.NewInspector()
	if err != nil {
		return fmt.Errorf("couldn't init docker client %w", err)
	}
	service := process.NewService(
		params.DeploymentParams,
		docker.NewDeployer(dispatcher, run.executor),
		inspector,
		git.NewFetcher(params.GetAddWritePerm(), params.GetRepoDir()),
		store,
		configStore,
		dispatcher,
		scheduler)
	go func() {
		_, err = scheduler.Schedule(func() {
			_, err := service.SyncDeployment()
			if err != nil {
				slog.Error(err.Error())
			}
		})
		if err != nil {
			slog.Warn(err.Error())
		}
	}()
	server := server.NewServer(store, configStore, service)
	return server.Serve(params.Port)
}
