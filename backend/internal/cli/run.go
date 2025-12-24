package cli

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"omar-kada/autonas/internal/cli/defaults"
	"omar-kada/autonas/internal/docker"
	"omar-kada/autonas/internal/events"
	"omar-kada/autonas/internal/git"
	"omar-kada/autonas/internal/process"
	"omar-kada/autonas/internal/server"
	"omar-kada/autonas/internal/storage"
	"omar-kada/autonas/models"

	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	_file         defaults.VarKey = "file"
	_workingDir   defaults.VarKey = "working-dir"
	_servicesDir  defaults.VarKey = "services-dir"
	_addWritePerm defaults.VarKey = "add-write-perm"
	_port         defaults.VarKey = "port"
)

var varInfoMap = defaults.VariableInfoMap{
	_file:         {EnvKey: "AUTONAS_CONFIG_FILE", DefaultValue: "/data/config.yaml"},
	_workingDir:   {EnvKey: "AUTONAS_WORKING_DIR", DefaultValue: "./config"},
	_servicesDir:  {EnvKey: "AUTONAS_SERVICES_DIR", DefaultValue: "."},
	_addWritePerm: {DefaultValue: false},
	_port:         {EnvKey: "AUTONAS_PORT", DefaultValue: 5005},
}

// RunParams contain parameters of the run command
type RunParams struct {
	models.DeploymentParams
	models.ServerParams
	ConfigFile string
}

func getParamsWithDefaults(p RunParams) RunParams {
	return RunParams{
		ConfigFile: varInfoMap.EnvOrDefault(p.ConfigFile, _file),
		DeploymentParams: models.DeploymentParams{
			WorkingDir:   varInfoMap.EnvOrDefault(p.WorkingDir, _workingDir),
			ServicesDir:  varInfoMap.EnvOrDefault(p.ServicesDir, _servicesDir),
			AddWritePerm: p.AddWritePerm,
		},
		ServerParams: models.ServerParams{
			Port: varInfoMap.EnvOrDefaultInt(p.Port, _port),
		},
	}
}

// newRunCommand creates a new run
func newRunCommand() *cobra.Command {
	params := RunParams{}
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run with optional config files",
		Run: func(_ *cobra.Command, _ []string) {
			if err := doRun(getParamsWithDefaults(params)); err != nil {
				slog.Error(err.Error())
			}
		},
	}

	runCmd.Flags().StringVarP(&(params.ConfigFile), string(_file), "f", "",
		varInfoMap.GetDefaultString("YAML config file", _file))
	runCmd.Flags().StringVarP(&(params.WorkingDir), string(_workingDir), "d", "",
		varInfoMap.GetDefaultString("directory where autonas data will be stored", _workingDir))
	runCmd.Flags().StringVarP(&(params.ServicesDir), string(_servicesDir), "s", "",
		varInfoMap.GetDefaultString("directory where services compose stacks will be stored", _servicesDir))
	runCmd.Flags().BoolVar(&(params.AddWritePerm), string(_addWritePerm), false,
		varInfoMap.GetDefaultString("when true, the tool adds write permission to config files", _addWritePerm))
	runCmd.Flags().IntVarP(&(params.Port), string(_port), "p", 5005,
		varInfoMap.GetDefaultString("port that will be used for exposing the API", _port))
	return runCmd
}

func initGorm(dbFile string, addPerm os.FileMode) (*gorm.DB, error) {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(dbFile), 0o600|addPerm); err != nil {
			return nil, err
		}
	}
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	// pragmas and pooling
	db.Exec("PRAGMA journal_mode=WAL;")
	db.Exec("PRAGMA foreign_keys = ON;")
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxOpenConns(1)
		sqlDB.SetMaxIdleConns(1)
	}
	return db, err
}

func doRun(params RunParams) error {
	dbFile := filepath.Join(params.GetDBDir(), "autonas.db")
	db, err := initGorm(dbFile, params.GetAddWritePerm())
	if err != nil {
		return fmt.Errorf("couldn't init sqlite db %w", err)
	}
	store, err := storage.NewGormStorage(db)
	if err != nil {
		return fmt.Errorf("couldn't init gorm storage %w", err)
	}

	dispatcher := events.NewDefaultDispatcher(store)
	configStore := storage.NewConfigStore(params.ConfigFile)
	scheduler := process.NewConfigScheduler(configStore)
	inspector, err := docker.NewInspector()
	if err != nil {
		return fmt.Errorf("couldn't init docker client %w", err)
	}
	service := process.NewService(
		params.DeploymentParams,
		docker.NewDeployer(dispatcher),
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
	server := server.NewServer(store, service)
	return server.Serve(params.Port)
}
