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
	GetCurrentStats(days int) (models.Stats, error)
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
	scheduler ConfigScheduler,
) Service {
	return &service{
		containersDeployer:  containersDeployer,
		containersInspector: containersInspector,
		fetcher:             fetcher,
		store:               store,
		dispatcher:          dispatcher,
		params:              deployParams,
		scheduler:           scheduler,
	}
}

// service is responsible for deploying the services
type service struct {
	containersDeployer  docker.Deployer
	containersInspector docker.Inspector
	fetcher             git.Fetcher
	store               storage.DeploymentStorage
	dispatcher          events.Dispatcher
	scheduler           ConfigScheduler
	params              models.DeploymentParams

	currentCfg models.Config
	mu         sync.Mutex
}

func (s *service) SyncDeployment(cfg models.Config) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	slog.Info("deploying from " + cfg.Repo + "/" + cfg.Branch)

	var syncErr error
	var patch git.Patch
	if s.currentCfg.Repo == "" || cfg.Repo == s.currentCfg.Repo {
		patch, syncErr = s.fetcher.Fetch(cfg.Repo, cfg.Branch, s.params.GetRepoDir())
	} else {
		patch, syncErr = s.fetcher.ReFetch(cfg.Repo, cfg.Branch, s.params.GetRepoDir())
	}

	if syncErr != nil && syncErr != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("error getting config repo:  %w", syncErr)
	}

	// check if the config changed from last run
	if patch.Diff == "" && reflect.DeepEqual(s.currentCfg, cfg) {
		slog.Info("Configuration and repository are up to date. No changes detected.")
		return nil
	}
	title := patch.Title
	if title == "" {
		title = "Automatic Deploy"
	}

	deployment, err := s.store.InitDeployment(title, patch.Author, patch.Diff, patch.Files)
	if err != nil {
		return err
	}
	ctx := context.WithValue(context.Background(), events.ObjectID, deployment.ID)

	defer func() {
		if err != nil {
			s.dispatcher.Error(ctx, fmt.Sprintf("Deployment failed %v", err))
			s.store.EndDeployment(deployment.ID, "error")
		} else {
			s.dispatcher.Info(ctx, "Deployment done successfully")
			err = s.store.EndDeployment(deployment.ID, "success")
		}
	}()

	err = s.containersDeployer.WithCtx(ctx).RemoveAndDeployStacks(s.currentCfg, cfg, s.params)
	if err != nil {
		return fmt.Errorf("error deploying services: %w", err)
	}

	s.currentCfg = cfg
	return nil
}

// GetManagedStacks returns a map of all containers managed by the tool
func (s *service) GetManagedStacks() (map[string][]models.ContainerSummary, error) {
	return s.containersInspector.GetManagedStacks(s.params.ServicesDir)
}

// GetCurrentStats returns the statistics of deployments for the last N days
func (s *service) GetCurrentStats(_ int) (models.Stats, error) {
	deps, err := s.store.GetDeployments()
	if err != nil {
		return models.Stats{}, err
	}
	var stats models.Stats
	for _, d := range deps {
		switch d.Status {
		case models.DeploymentStatusSuccess:
			stats.Success++
		case models.DeploymentStatusError:
			stats.Error++
		}
	}
	last, err := s.store.GetLastDeployment()
	if err == nil {
		stats.Author = last.Author
		stats.LastDeploy = last.Time
		stats.LastStatus = last.Status
	}
	stats.NextDeploy = s.scheduler.getNext()
	return stats, nil
}
