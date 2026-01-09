package cli

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"time"

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
	"gorm.io/gorm/logger"
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
	_addWritePerm: {EnvKey: "AUTONAS_ADD_WRITE_PERM", DefaultValue: "false"},
	_port:         {EnvKey: "AUTONAS_PORT", DefaultValue: 5005},
}

// RunParams contain parameters of the run command
type RunParams struct {
	models.DeploymentParams
	models.ServerParams
	ConfigFile string
}

func getParamsWithDefaults(p RunParams, addWritePerm string) RunParams {
	return RunParams{
		ConfigFile: varInfoMap.EnvOrDefault(p.ConfigFile, _file),
		DeploymentParams: models.DeploymentParams{
			WorkingDir:   varInfoMap.EnvOrDefault(p.WorkingDir, _workingDir),
			ServicesDir:  varInfoMap.EnvOrDefault(p.ServicesDir, _servicesDir),
			AddWritePerm: varInfoMap.EnvOrDefault(addWritePerm, _addWritePerm),
		},
		ServerParams: models.ServerParams{
			Port: varInfoMap.EnvOrDefaultInt(p.Port, _port),
		},
	}
}

// newRunCommand creates a new run
func newRunCommand() *cobra.Command {
	params := RunParams{}
	addWritePerm := ""
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run with optional config files",
		Run: func(_ *cobra.Command, _ []string) {
			if err := doRun(getParamsWithDefaults(params, addWritePerm)); err != nil {
				slog.Error(err.Error())
			}
		},
	}
	runCmd.Flags().StringVarP(&params.ConfigFile, string(_file), "f", "",
		varInfoMap.GetDefaultString("YAML config file", _file))
	runCmd.Flags().StringVarP(&params.WorkingDir, string(_workingDir), "d", "",
		varInfoMap.GetDefaultString("directory where autonas data will be stored", _workingDir))
	runCmd.Flags().StringVarP(&params.ServicesDir, string(_servicesDir), "s", "",
		varInfoMap.GetDefaultString("directory where services compose stacks will be stored", _servicesDir))
	runCmd.Flags().StringVar(&params.AddWritePerm, string(_addWritePerm), "",
		varInfoMap.GetDefaultString("when true, the tool adds write permission to config files", _addWritePerm))
	runCmd.Flags().IntVarP(&params.Port, string(_port), "p", 5005,
		varInfoMap.GetDefaultString("port that will be used for exposing the API", _port))

	return runCmd
}

func initGorm(dbFile string, addPerm os.FileMode) (*gorm.DB, error) {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(dbFile), 0o600|addPerm); err != nil {
			return nil, err
		}
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Enable color
		},
	)

	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{
		Logger: newLogger,
	})
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
	// create a configStore from the config file,
	// laod the config from there, and use it for the rest of the settings
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
