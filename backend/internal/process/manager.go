// Package process handles the deployment and management of services.
package process

import (
	"fmt"
	"log/slog"
	"omar-kada/autonas/api"
	"omar-kada/autonas/internal/config"
	"omar-kada/autonas/internal/files"
	"omar-kada/autonas/internal/git"
	"omar-kada/autonas/internal/storage"
	"omar-kada/autonas/models"
	"reflect"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

const DeploymentID = "deployment_id"

// Deployer defines methods for managing containerized services.
type Deployer interface {
	WithLogger(log *slog.Logger) Deployer
	RemoveServices(services []string, servicesDir string) map[string]error
	DeployServices(cfg config.Config, servicesDir string) map[string]error
	RemoveAndDeployStacks(oldCfg, cfg config.Config, params models.DeploymentParams) error
}

// Inspector defined operations for info retreival on containers
type Inspector interface {
	GetManagedContainers(servicesDir string) (map[string][]models.ContainerSummary, error)
}

// Manager abstracts service deployment operations
type Manager interface {
	SyncDeployment() error
	SyncPeriodically() *cron.Cron

	GetManagedContainers() (map[string][]models.ContainerSummary, error)
}

// NewManager creates a new Manager instance
func NewManager(
	managerParams models.DeploymentParams,
	containersDeployer Deployer,
	containersInspector Inspector,
	copier files.Copier,
	fetcher git.Fetcher,
	store storage.Storage,
	configGenerator config.Generator,
) Manager {
	return &manager{
		containersDeployer:  containersDeployer,
		containersInspector: containersInspector,
		copier:              copier,
		fetcher:             fetcher,
		configGenerator:     configGenerator,
		store:               store,
		params:              managerParams,
	}
}

// manager is responsible for deploying the services
type manager struct {
	containersDeployer  Deployer
	containersInspector Inspector
	copier              files.Copier
	fetcher             git.Fetcher
	configGenerator     config.Generator
	store               storage.Storage
	params              models.DeploymentParams

	currentCfg config.Config
	cron       *cron.Cron
	mu         sync.Mutex
}

func (m *manager) SyncDeployment() (err error) {
	// make sure only one sync job is running at a time
	m.mu.Lock()
	defer m.mu.Unlock()

	deployment := api.Deployment{
		Id:     fmt.Sprintf("%v", time.Now()),
		Title:  "Auto deployment",
		Time:   time.Now(),
		Status: "running",
		Diff:   "",
		Logs:   []string{},
	}
	m.store.SaveDeployment(
		deployment,
	)

	defer func() {
		if err != nil {
			m.store.UpdateStatus(
				deployment.Id, "failed",
			)
		} else {
			m.store.UpdateStatus(
				deployment.Id, "success",
			)
		}
	}()
	log := NewLoggerWith(func(record slog.Record) {
		m.store.AddLogRecord(deployment.Id, record)
	})

	cfg, err := m.configGenerator.FromFiles([]string{m.params.ConfigFile})
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}
	log.Debug("Final consolidated config", "config", cfg)

	oldConfig := m.currentCfg
	m.currentCfg = cfg

	if cfg.CronPeriod != oldConfig.CronPeriod {
		m.SyncPeriodically()
	}

	syncErr := m.fetcher.Fetch(cfg.Repo, cfg.Branch, m.params.WorkingDir)

	if syncErr != nil && syncErr != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("error getting config repo:  %w", syncErr)
	}

	// check if the config changed from last run
	if syncErr == git.NoErrAlreadyUpToDate && reflect.DeepEqual(oldConfig, cfg) {
		log.Info("Configuration and repository are up to date. No changes detected.")
		return nil
	}

	// Copy all files from ./services to SERVICES_PATH
	err = m.containersDeployer.WithLogger(log).RemoveAndDeployStacks(oldConfig, cfg, m.params)
	if err != nil {
		return fmt.Errorf("error deploying services: %w", err)
	}

	return nil
}

func (m *manager) SyncPeriodically() *cron.Cron {
	if m.cron != nil {
		m.cron.Stop()
		m.cron = nil
	}
	if m.currentCfg.CronPeriod == "" {
		slog.Warn("no cron period configured, will no schedule sync jobs")
		return nil
	}
	c := cron.New()

	c.AddFunc(m.currentCfg.CronPeriod, func() {
		err := m.SyncDeployment()
		if err != nil {
			slog.Error("error on run periodically", "error", err)
		}
	})

	c.Start()
	m.cron = c
	slog.Info("values for cron job", "entries", c.Entries())
	return c
}

// GetManagedContainers returns a map of all containers managed by the tool
func (m *manager) GetManagedContainers() (map[string][]models.ContainerSummary, error) {
	return m.containersInspector.GetManagedContainers(m.params.ServicesDir)
}
