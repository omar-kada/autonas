package cli

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"omar-kada/autonas/internal/docker"
	"omar-kada/autonas/internal/events"
	"omar-kada/autonas/internal/git"
	"omar-kada/autonas/internal/process"
	"omar-kada/autonas/internal/server"
	"omar-kada/autonas/internal/shell"
	"omar-kada/autonas/internal/storage"

	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type runCommand struct {
	executor shell.Executor

	cmd    *cobra.Command
	params RunParams
}

// NewRunCommand creates a new run
func NewRunCommand(executor shell.Executor) *cobra.Command {

	run := runCommand{
		params:   RunParams{},
		executor: executor,
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

func (*runCommand) initGorm(dbFile string, addPerm os.FileMode) (*gorm.DB, error) {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(dbFile), 0o700|addPerm); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	fmt.Printf("dbFile = %v\n", dbFile)

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

func (run *runCommand) doRun() error {
	params := getParamsWithDefaults(run.params)
	// create a configStore from the config file,
	// laod the config from there, and use it for the rest of the settings
	dbFile := filepath.Join(params.GetDBDir(), "autonas.db")
	db, err := run.initGorm(dbFile, params.GetAddWritePerm())
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
	server := server.NewServer(store, service)
	return server.Serve(params.Port)
}
