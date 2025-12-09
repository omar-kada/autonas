// Package process handles the deployment and management of services.
package process

import (
	"context"
	"fmt"
	"log/slog"
	"omar-kada/autonas/internal/events"
	"omar-kada/autonas/internal/git"
	"omar-kada/autonas/internal/storage"
	"omar-kada/autonas/models"
	"reflect"
)

// DeploymentID is the key used to store deployment ID in context
const DeploymentID = "deployment_id"

// Deployer defines methods for managing containerized services.
type Deployer interface {
	WithCtx(ctx context.Context) Deployer
	RemoveServices(services []string, servicesDir string) map[string]error
	DeployServices(cfg models.Config, servicesDir string) map[string]error
	RemoveAndDeployStacks(oldCfg, cfg models.Config, params models.DeploymentParams) error
}

// Inspector defined operations for info retreival on containers
type Inspector interface {
	GetManagerStacks(servicesDir string) (map[string][]models.ContainerSummary, error)
}

// Service abstracts service deployment operations
type Service interface {
	SyncDeployment(cfg models.Config) error

	GetManagerStacks() (map[string][]models.ContainerSummary, error)
}

// NewService creates a new process Service instance
func NewService(
	managerParams models.DeploymentParams,
	containersDeployer Deployer,
	containersInspector Inspector,
	fetcher git.Fetcher,
	store storage.DeploymentStorage,
	dispatcher events.Dispatcher,
) Service {
	return &service{
		containersDeployer:  containersDeployer,
		containersInspector: containersInspector,
		fetcher:             fetcher,
		store:               store,
		dispatcher:          dispatcher,
		params:              managerParams,
	}
}

// service is responsible for deploying the services
type service struct {
	containersDeployer  Deployer
	containersInspector Inspector
	fetcher             git.Fetcher
	store               storage.DeploymentStorage
	dispatcher          events.Dispatcher
	params              models.DeploymentParams

	currentCfg models.Config
}

func (m *service) SyncDeployment(cfg models.Config) (err error) {

	syncErr := m.fetcher.Fetch(cfg.Repo, cfg.Branch, m.params.WorkingDir)

	if syncErr != nil && syncErr != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("error getting config repo:  %w", syncErr)
	}

	// check if the config changed from last run
	if syncErr == git.NoErrAlreadyUpToDate && reflect.DeepEqual(m.currentCfg, cfg) {
		slog.Info("Configuration and repository are up to date. No changes detected.")
		return nil
	}

	deployment, err := m.store.InitDeployment("Auto deployment", "")
	if err != nil {
		return err
	}
	ctx := context.WithValue(context.Background(), events.ObjectID, deployment.Id)

	defer func() {
		if err != nil {
			m.dispatcher.Error(ctx, fmt.Sprintf("Deployment failed %v", err))
			m.store.UpdateStatus(deployment.Id, "failed")
		} else {
			m.dispatcher.Info(ctx, "Deployment done successfully")
			m.store.UpdateStatus(deployment.Id, "success")
		}
	}()

	// Copy all files from ./services to SERVICES_PATH
	err = m.containersDeployer.WithCtx(ctx).RemoveAndDeployStacks(m.currentCfg, cfg, m.params)
	if err != nil {
		return fmt.Errorf("error deploying services: %w", err)
	}

	m.currentCfg = cfg
	return nil
}

// GetManagerStacks returns a map of all containers managed by the tool
func (m *service) GetManagerStacks() (map[string][]models.ContainerSummary, error) {
	return m.containersInspector.GetManagerStacks(m.params.ServicesDir)
}
