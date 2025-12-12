// Package process handles the deployment and management of services.
package process

import (
	"context"
	"fmt"
	"log/slog"
	"omar-kada/autonas/internal/docker"
	"omar-kada/autonas/internal/events"
	"omar-kada/autonas/internal/git"
	"omar-kada/autonas/internal/storage"
	"omar-kada/autonas/models"
	"reflect"
	"sync"
)

// DeploymentID is the key used to store deployment ID in context
const DeploymentID = "deployment_id"

// Service abstracts service deployment operations
type Service interface {
	SyncDeployment(cfg models.Config) error

	GetManagedStacks() (map[string][]models.ContainerSummary, error)
}

// NewService creates a new process Service instance
func NewService(
	deployParams models.DeploymentParams,
	containersDeployer docker.Deployer,
	containersInspector docker.Inspector,
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
		params:              deployParams,
	}
}

// service is responsible for deploying the services
type service struct {
	containersDeployer  docker.Deployer
	containersInspector docker.Inspector
	fetcher             git.Fetcher
	store               storage.DeploymentStorage
	dispatcher          events.Dispatcher
	params              models.DeploymentParams

	currentCfg models.Config
	mu         sync.Mutex
}

func (m *service) SyncDeployment(cfg models.Config) (err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	slog.Info("deploying from " + cfg.Repo + "/" + cfg.Branch)

	var syncErr error
	var patch git.Patch
	if m.currentCfg.Repo == "" || cfg.Repo == m.currentCfg.Repo {
		patch, syncErr = m.fetcher.Fetch(cfg.Repo, cfg.Branch, m.params.GetRepoDir())
	} else {
		patch, syncErr = m.fetcher.ReFetch(cfg.Repo, cfg.Branch, m.params.GetRepoDir())
	}

	if syncErr != nil && syncErr != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("error getting config repo:  %w", syncErr)
	}

	// check if the config changed from last run
	if patch.Diff == "" && reflect.DeepEqual(m.currentCfg, cfg) {
		slog.Info("Configuration and repository are up to date. No changes detected.")
		return nil
	}

	deployment, err := m.store.InitDeployment(patch.Title, patch.Author, patch.Diff)
	if err != nil {
		return err
	}
	ctx := context.WithValue(context.Background(), events.ObjectID, deployment.Id)

	defer func() {
		if err != nil {
			m.dispatcher.Error(ctx, fmt.Sprintf("Deployment failed %v", err))
			m.store.UpdateStatus(deployment.Id, "error")
		} else {
			m.dispatcher.Info(ctx, "Deployment done successfully")
			m.store.UpdateStatus(deployment.Id, "success")
		}
	}()

	err = m.containersDeployer.WithCtx(ctx).RemoveAndDeployStacks(m.currentCfg, cfg, m.params)
	if err != nil {
		return fmt.Errorf("error deploying services: %w", err)
	}

	m.currentCfg = cfg
	return nil
}

// GetManagedStacks returns a map of all containers managed by the tool
func (m *service) GetManagedStacks() (map[string][]models.ContainerSummary, error) {
	return m.containersInspector.GetManagedStacks(m.params.ServicesDir)
}
