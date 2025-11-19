// Package process handles the deployment and management of services.
package process

import (
	"fmt"
	"log/slog"
	"omar-kada/autonas/internal/config"
	"omar-kada/autonas/internal/files"
	"omar-kada/autonas/internal/git"
	"omar-kada/autonas/models"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"sync"

	"github.com/robfig/cron/v3"
)

// Deployer defines methods for managing containerized services.
type Deployer interface {
	RemoveServices(services []string, servicesDir string) error
	DeployServices(cfg config.Config, servicesDir string) error
	GetManagedContainers(servicesDir string) (map[string][]models.ContainerSummary, error)
}

// Manager abstracts service deployment operations
type Manager interface {
	GetManagedContainers() (map[string][]models.ContainerSummary, error)
	GetCurrentCfg() config.Config

	SyncDeployment() error
	SyncPeriodically() error
}

// NewManager creates a new Manager instance
func NewManager(
	managerParams ManagerParams,
	containersDeployer Deployer,
	copier files.Copier,
	fetcher git.Fetcher,
	configGenerator config.Generator,
) Manager {
	return &manager{
		containersDeployer: containersDeployer,
		copier:             copier,
		fetcher:            fetcher,
		configGenerator:    configGenerator,
		params:             managerParams,
	}
}

// ManagerParams are tool params that doesn't change in runtime
type ManagerParams struct {
	AddPerm     os.FileMode
	ServicesDir string
	WorkingDir  string
	ConfigFile  string
}

// manager is responsible for deploying the services
type manager struct {
	containersDeployer Deployer
	copier             files.Copier
	fetcher            git.Fetcher
	configGenerator    config.Generator
	params             ManagerParams
	currentCfg         config.Config
	cron               *cron.Cron
	mu                 sync.Mutex
}

func (m *manager) SyncDeployment() error {
	// make sure only one sync job is running at a time
	m.mu.Lock()
	defer m.mu.Unlock()
	slog.Warn(fmt.Sprintf("generating config from %s", m.params.ConfigFile))
	cfg, err := m.configGenerator.FromFiles([]string{m.params.ConfigFile})
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}
	slog.Debug("Final consolidated config", "config", cfg)

	oldConfig := m.currentCfg
	m.currentCfg = cfg

	if cfg.CronPeriod != oldConfig.CronPeriod {
		err = m.SyncPeriodically()
		if err != nil {
			return err
		}
	}

	syncErr := m.fetcher.Fetch(cfg.Repo, cfg.Branch, m.params.WorkingDir)

	if syncErr != nil && syncErr != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("error getting config repo:  %w", syncErr)
	}

	// check if the config changed from last run
	if syncErr == git.NoErrAlreadyUpToDate && reflect.DeepEqual(m.GetCurrentCfg(), cfg) {
		slog.Info("Configuration and repository are up to date. No changes detected.")
		return nil
	}

	// Copy all files from ./services to SERVICES_PATH
	err = m.removeAndDeployStacks(oldConfig, cfg)
	if err != nil {
		return fmt.Errorf("error deploying services: %w", err)
	}
	return nil
}

func (m *manager) SyncPeriodically() error {
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
	return nil
}

// removeAndDeployStacks handles the deployment/removal of services based on the current and new configuration.
func (m *manager) removeAndDeployStacks(oldCfg, cfg config.Config) error {
	toBeRemoved := getUnusedServices(oldCfg, cfg)
	// TODO : check if the stack is up before calling RemoveServices
	if err := m.containersDeployer.RemoveServices(toBeRemoved, m.params.ServicesDir); err != nil {
		return err
	}

	slog.Debug("copying files from src to dst", "src", m.params.WorkingDir+"/services", "dst", m.params.ServicesDir)

	enabledServiecs := cfg.GetEnabledServices()
	for _, service := range enabledServiecs {
		src := filepath.Join(m.params.WorkingDir, "services", service)
		dst := filepath.Join(m.params.ServicesDir, service)
		if err := m.copier.CopyWithAddPerm(src, dst, m.params.AddPerm); err != nil {
			return fmt.Errorf("error while copying service "+service+" %w", err)
		}
	}

	slog.Debug("deploying enabled services", "services", enabledServiecs)
	if err := m.containersDeployer.DeployServices(cfg, m.params.ServicesDir); err != nil {
		return err
	}
	return nil
}

// GetManagedContainers returns a map of all containers managed by the tool
func (m *manager) GetManagedContainers() (map[string][]models.ContainerSummary, error) {
	return m.containersDeployer.GetManagedContainers(m.params.ServicesDir)
}

func getUnusedServices(oldCfg, cfg config.Config) []string {
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

func (m *manager) GetCurrentCfg() config.Config {
	return m.currentCfg
}
