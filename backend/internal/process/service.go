// Package process handles the deployment and management of services.
package process

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"sync"

	"omar-kada/autonas/internal/docker"
	"omar-kada/autonas/internal/events"
	"omar-kada/autonas/internal/git"
	"omar-kada/autonas/internal/storage"
	"omar-kada/autonas/models"

	"github.com/moby/moby/api/types/container"
)

// DeploymentID is the key used to store deployment ID in context
const DeploymentID = "deployment_id"

// Service abstracts service deployment operations
type Service interface {
	SyncDeployment() (models.Deployment, error)
	GetCurrentStats(days int) (models.Stats, error)
	GetDiff() ([]models.FileDiff, error)
	GetManagedStacks() (map[string][]models.ContainerSummary, error)
	GetDeployments(first int, after uint64) ([]models.Deployment, error)
}

// NewService creates a new process Service instance
func NewService(
	deployParams models.DeploymentParams,
	containersDeployer docker.Deployer,
	containersInspector docker.Inspector,
	fetcher git.Fetcher,
	store storage.DeploymentStorage,
	configStore storage.ConfigStore,
	dispatcher events.Dispatcher,
	scheduler ConfigScheduler,
) Service {
	cfg, _ := configStore.Get()
	return &service{
		containersDeployer:  containersDeployer,
		containersInspector: containersInspector,
		fetcher:             fetcher,
		store:               store,
		configStore:         configStore,
		dispatcher:          dispatcher,
		params:              deployParams,
		scheduler:           scheduler,
		currentCfg:          cfg,
	}
}

// service is responsible for deploying the services
type service struct {
	containersDeployer  docker.Deployer
	containersInspector docker.Inspector
	fetcher             git.Fetcher
	store               storage.DeploymentStorage
	configStore         storage.ConfigStore
	dispatcher          events.Dispatcher
	scheduler           ConfigScheduler
	params              models.DeploymentParams

	currentCfg models.Config
	mu         sync.Mutex
}

func (s *service) SyncDeployment() (models.Deployment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cfg, err := s.configStore.Get()
	if err != nil {
		return models.Deployment{}, fmt.Errorf("error getting current config:  %w", err)
	}
	fetcher := s.fetcher.WithConfig(cfg)
	slog.Info("deploying from " + cfg.Repo + "/" + cfg.Branch)

	patch, syncErr := fetcher.DiffWithRemote()

	if syncErr != nil && syncErr != git.NoErrAlreadyUpToDate {
		return models.Deployment{}, fmt.Errorf("error getting config repo:  %w", syncErr)
	}

	// check if the config changed from last run
	if patch.Diff == "" && reflect.DeepEqual(s.currentCfg, cfg) && s.areStacksHealthy(cfg) {
		slog.Info("Configuration and repository are up to date. No changes detected.",
			"oldConfig", s.currentCfg, "newConfig", cfg, "diff", patch.Diff)
		return models.Deployment{}, nil
	}
	title := patch.Title
	if title == "" {
		title = "Automatic Deploy"
	}

	deployment, err := s.store.InitDeployment(title, patch.Author, patch.Diff, patch.Files)
	if err != nil {
		return deployment, err
	}
	go func() {
		ctx := context.WithValue(context.Background(), events.ObjectID, deployment.ID)
		err := fetcher.PullBranch("to_be_deployed", "")
		if err != nil {
			s.updateDeploymentStatus(ctx, deployment, err)
			return
		}
		s.dispatcher.Info(ctx, "Pulled new changes into working branch")

		err = s.containersDeployer.WithCtx(ctx).RemoveAndDeployStacks(s.currentCfg, cfg, s.params)
		if err != nil {
			s.updateDeploymentStatus(ctx, deployment, err)
			return
		}
		s.dispatcher.Info(ctx, "Deployed changes to running stacks")

		err = fetcher.PullBranch(cfg.Branch, patch.CommitHash)
		s.updateDeploymentStatus(ctx, deployment, err)

		s.currentCfg = cfg
	}()

	return deployment, nil
}

func (s *service) areStacksHealthy(cfg models.Config) bool {
	runningStacks, err := s.containersInspector.GetManagedStacks(s.params.ServicesDir)
	if err != nil {
		return false
	}
	enabledServices := cfg.GetEnabledServices()
	for _, service := range enabledServices {
		serviceContainers := runningStacks[service]
		if len(serviceContainers) == 0 {
			return false
		}
		for _, ctr := range serviceContainers {
			if ctr.Health == container.Unhealthy {
				return false
			}
		}
	}
	slog.Info("services are healthy ", "stacks", runningStacks)
	return true
}

func (s *service) updateDeploymentStatus(ctx context.Context, deployment models.Deployment, err error) {
	if err != nil {
		s.dispatcher.Error(ctx, fmt.Sprintf("Deployment failed %v", err))
		slog.Error(err.Error())
		s.store.EndDeployment(deployment.ID, "error")
	} else {
		s.dispatcher.Info(ctx, "Deployment done successfully")
		s.store.EndDeployment(deployment.ID, "success")
	}
}

// GetManagedStacks returns a map of all containers managed by the tool
func (s *service) GetManagedStacks() (map[string][]models.ContainerSummary, error) {
	return s.containersInspector.GetManagedStacks(s.params.ServicesDir)
}

// GetCurrentStats returns the statistics of deployments for the last N days
func (s *service) GetCurrentStats(_ int) (models.Stats, error) {
	deps, err := s.store.GetDeployments(storage.NewIDCursor(100, 0))
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
	stats.NextDeploy = s.scheduler.GetNext()

	stacks, err := s.GetManagedStacks()
	if err != nil {
		return models.Stats{}, err
	}
	stats.Health = getGlobalHealth(stacks)
	return stats, nil
}

func getGlobalHealth(stacks map[string][]models.ContainerSummary) container.HealthStatus {
	for _, containers := range stacks {
		for _, ctnr := range containers {
			if container.Unhealthy == ctnr.Health {
				return container.Unhealthy
			}
		}
	}
	return container.Healthy
}

// GetDiff returns the changed files between what's deployed and the repo
func (s *service) GetDiff() ([]models.FileDiff, error) {
	cfg, err := s.getConfig()
	if err != nil {
		return []models.FileDiff{}, fmt.Errorf("error while getting config : %w", err)
	}
	patch, err := s.fetcher.WithConfig(cfg).DiffWithRemote()
	if err != nil {
		return []models.FileDiff{}, err
	}
	return patch.Files, nil
}

func (s *service) getConfig() (models.Config, error) {
	if s.currentCfg.Repo != "" {
		return s.currentCfg, nil
	}
	return s.configStore.Get()
}

// GetDeployments returns a paginated list of deployments.
func (s *service) GetDeployments(first int, after uint64) ([]models.Deployment, error) {
	return s.store.GetDeployments(storage.NewIDCursor(first, after))
}
