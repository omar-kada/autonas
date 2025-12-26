package process

import (
	"context"
	"errors"
	"testing"
	"time"

	"omar-kada/autonas/internal/docker"
	"omar-kada/autonas/internal/events"
	"omar-kada/autonas/internal/git"
	"omar-kada/autonas/internal/storage"
	"omar-kada/autonas/models"

	"github.com/moby/moby/api/types/container"
	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type Mocker struct {
	mock.Mock
}

func (m *Mocker) WithCtx(_ context.Context) docker.Deployer {
	return m
}

func (m *Mocker) RemoveServices(services []string, servicesDir string) map[string]error {
	args := m.Called(services, servicesDir)
	return args.Get(0).(map[string]error)
}

func (m *Mocker) DeployServices(cfg models.Config, servicesDir string) map[string]error {
	args := m.Called(cfg, servicesDir)
	return args.Get(0).(map[string]error)
}

func (m *Mocker) RemoveAndDeployStacks(oldCfg, cfg models.Config, params models.DeploymentParams) error {
	args := m.Called(oldCfg, cfg, params)
	return args.Error(0)
}

func (m *Mocker) GetManagedStacks(servicesDir string) (map[string][]models.ContainerSummary, error) {
	args := m.Called(servicesDir)
	return args.Get(0).(map[string][]models.ContainerSummary), args.Error(1)
}

func (m *Mocker) GetNext() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}

func (m *Mocker) Schedule(fn func()) (*cron.Cron, error) {
	args := m.Called(fn)
	return args.Get(0).(*cron.Cron), args.Error(1)
}

func (m *Mocker) ClearRepo() error {
	args := m.Called()
	return args.Error(0)
}

func (m *Mocker) CheckoutBranch(branch string) error {
	args := m.Called(branch)
	return args.Error(0)
}

func (m *Mocker) PullBranch(branch string, commitSHA string) error {
	args := m.Called(branch, commitSHA)
	return args.Error(0)
}

func (m *Mocker) WithConfig(cfg models.Config) git.Fetcher {
	args := m.Called(cfg)
	return args.Get(0).(git.Fetcher)
}

func (m *Mocker) DiffWithRemote() (git.Patch, error) {
	args := m.Called()
	return args.Get(0).(git.Patch), args.Error(1)
}

var (
	mockConfigOld = models.Config{
		Repo:   "https://example.com/repo.git",
		Branch: "main",
		Services: map[string]models.ServiceConfig{
			"svc1": {
				// Extra: map[string]any{
				"Port":    8080,
				"Version": "v1",
				//},
			},
			"svc2": {},
		},
	}
	mockConfigNew = models.Config{
		Repo:   "https://example.com/repo.git",
		Branch: "main",
		Services: map[string]models.ServiceConfig{
			"svc2": {},
			"svc3": {},
		},
	}
)

func newServiceWithCurrentConfig(t *testing.T, mocker *Mocker, params models.DeploymentParams, currentCfg models.Config) *service {
	configStore := storage.NewConfigStore(t.TempDir() + "/config.yaml")
	configStore.Update(currentCfg)
	svc := NewService(
		params,
		mocker,
		mocker,
		mocker,
		storage.NewMemoryStorage(),
		configStore,
		events.NewVoidDispatcher(),
		mocker,
	).(*service)
	svc.currentCfg = currentCfg
	return svc
}

func newServiceWithMocks(t *testing.T, mocker *Mocker, params models.DeploymentParams) *service {
	return newServiceWithCurrentConfig(t, mocker, params, models.Config{})
}

var (
	ErrRemove   = errors.New("removeServices error")
	ErrDeploy   = errors.New("deployServices error")
	ErrGenerate = errors.New("generate file error")
	ErrFetch    = errors.New("sync config error")
)

func TestSync_Success(t *testing.T) {
	mocker := &Mocker{}
	service := newServiceWithMocks(t, mocker, models.DeploymentParams{
		ServicesDir: "/services",
		WorkingDir:  ".",
	})

	wantCfg := mockConfigOld
	service.configStore.Update(wantCfg)
	mocker.On("WithConfig", wantCfg).Return(service.fetcher)
	mocker.On("DiffWithRemote").Once().Return(git.Patch{Diff: "test"}, nil)
	mocker.On("PullBranch", "to_be_deployed", "").Once().Return(nil)
	mocker.On("RemoveAndDeployStacks", models.Config{}, wantCfg, service.params).Once().Return(nil)
	// signal when working branch pull completes
	done := make(chan struct{})
	mocker.On("PullBranch", "main", mock.Anything).Once().
		Return(nil).
		Run(func(_ mock.Arguments) { close(done) })

	dep, err := service.SyncDeployment()
	assert.NoError(t, err)
	assert.Equal(t, "Automatic Deploy", dep.Title)
	assert.Equal(t, "test", dep.Diff)
	assert.Equal(t, models.DeploymentStatusRunning, dep.Status)

	// wait for goroutine to finish its PullBranch step
	select {
	case <-done:
		// proceed with assertions
	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for background deployment goroutine")
	}

	newDep, err := service.store.GetDeployment(dep.ID)
	assert.NoError(t, err)
	assert.Equal(t, models.DeploymentStatusSuccess, newDep.Status)

	mocker.AssertExpectations(t)
}

func TestSync_Success_RedploymentWithChangedConfig(t *testing.T) {
	mocker := &Mocker{}
	service := newServiceWithCurrentConfig(t, mocker, models.DeploymentParams{
		ServicesDir: "/services",
		WorkingDir:  ".",
	}, mockConfigOld)

	wantCfg := mockConfigNew
	service.configStore.Update(wantCfg)
	mocker.On("WithConfig", wantCfg).Return(service.fetcher)
	mocker.On("DiffWithRemote").Once().Return(git.Patch{Diff: "test"}, nil)
	mocker.On("PullBranch", "to_be_deployed", "").Once().Return(nil)
	mocker.On("RemoveAndDeployStacks", mockConfigOld, wantCfg, service.params).Once().Return(nil)
	// signal when working branch pull completes
	done := make(chan struct{})
	mocker.On("PullBranch", "main", mock.Anything).Once().
		Return(nil).
		Run(func(_ mock.Arguments) { close(done) })
	dep, err := service.SyncDeployment()
	assert.NoError(t, err)
	assert.Equal(t, "Automatic Deploy", dep.Title)
	assert.Equal(t, models.DeploymentStatusRunning, dep.Status)

	// wait for goroutine to finish its PullBranch step
	select {
	case <-done:
		// proceed with assertions
	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for background deployment goroutine")
	}

	newDep, err := service.store.GetDeployment(dep.ID)
	assert.NoError(t, err)
	assert.Equal(t, models.DeploymentStatusSuccess, newDep.Status)

	mocker.AssertExpectations(t)
}

func TestSync_ErrorsOnPullbranch(t *testing.T) {
	mocker := &Mocker{}
	service := newServiceWithCurrentConfig(t, mocker, models.DeploymentParams{
		ServicesDir: "/services",
		WorkingDir:  ".",
	}, mockConfigOld)
	wantCfg := mockConfigNew
	service.configStore.Update(wantCfg)
	mocker.On("WithConfig", wantCfg).Return(service.fetcher)
	mocker.On("DiffWithRemote").Once().Return(git.Patch{Diff: "test"}, nil)
	done := make(chan struct{})
	mocker.On("PullBranch", "to_be_deployed", "").Once().Return(ErrFetch).
		Run(func(_ mock.Arguments) { close(done) })

	dep, err := service.SyncDeployment()
	assert.NoError(t, err)
	// wait for goroutine to finish its PullBranch step
	select {
	case <-done:
		// proceed with assertions
	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for background deployment goroutine")
	}
	time.Sleep(10 * time.Millisecond)
	newDep, err := service.store.GetDeployment(dep.ID)
	assert.NoError(t, err)
	assert.Equal(t, models.DeploymentStatusError, newDep.Status)
}

func TestSync_Errors(t *testing.T) {
	mocker := &Mocker{}
	service := newServiceWithCurrentConfig(t, mocker, models.DeploymentParams{
		ServicesDir: "/services",
		WorkingDir:  ".",
	}, mockConfigOld)
	wantCfg := mockConfigNew
	service.configStore.Update(wantCfg)
	mocker.On("WithConfig", wantCfg).Return(service.fetcher)
	mocker.On("DiffWithRemote").Once().Return(git.Patch{Diff: "test"}, nil)
	mocker.On("PullBranch", "to_be_deployed", "").Once().Return(nil)

	done := make(chan struct{})
	mocker.On("RemoveAndDeployStacks", mockConfigOld, wantCfg, service.params).
		Once().Return(ErrDeploy).
		Run(func(_ mock.Arguments) { close(done) })
	dep, err := service.SyncDeployment()
	assert.NoError(t, err)
	// wait for goroutine to finish its PullBranch step
	select {
	case <-done:
		// proceed with assertions
	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for background deployment goroutine")
	}
	time.Sleep(10 * time.Millisecond)

	newDep, err := service.store.GetDeployment(dep.ID)
	assert.NoError(t, err)
	assert.Equal(t, models.DeploymentStatusError, newDep.Status)
}

func TestGetCurrentStats_NoDeployments(t *testing.T) {
	mocker := &Mocker{}
	service := newServiceWithMocks(t, mocker, models.DeploymentParams{})

	next := time.Now().Add(1 * time.Hour)
	mocker.On("GetNext").Return(next)
	mocker.On("GetManagedStacks", mock.Anything).Return(map[string][]models.ContainerSummary{}, nil)

	stats, err := service.GetCurrentStats(7)
	assert.NoError(t, err)
	assert.Equal(t, int32(0), stats.Success)
	assert.Equal(t, int32(0), stats.Error)
	assert.Equal(t, "", stats.Author)
	assert.True(t, stats.LastDeploy.IsZero())
	assert.Equal(t, next, stats.NextDeploy)
}

func TestGetCurrentStats_WithDeployments(t *testing.T) {
	mocker := &Mocker{}
	service := newServiceWithMocks(t, mocker, models.DeploymentParams{})

	next := time.Now().Add(30 * time.Minute)
	mocker.On("GetNext").Return(next)
	mocker.On("GetManagedStacks", mock.Anything).Return(map[string][]models.ContainerSummary{}, nil)

	// create a successful deployment
	dep1, err := service.store.InitDeployment("first", "alice", "diff1", nil)
	assert.NoError(t, err)
	err = service.store.EndDeployment(dep1.ID, models.DeploymentStatusSuccess)
	assert.NoError(t, err)

	// create a failed (last) deployment
	dep2, err := service.store.InitDeployment("second", "bob", "diff2", nil)
	assert.NoError(t, err)
	err = service.store.EndDeployment(dep2.ID, models.DeploymentStatusError)
	assert.NoError(t, err)
	stats, err := service.GetCurrentStats(7)
	assert.NoError(t, err)
	assert.Equal(t, int32(1), stats.Success)
	assert.Equal(t, int32(1), stats.Error)
	assert.Equal(t, "bob", stats.Author)
	assert.Equal(t, models.DeploymentStatusError, stats.LastStatus)
	assert.Equal(t, next, stats.NextDeploy)
}

func TestGetManagedStacks(t *testing.T) {
	testCases := []struct {
		name           string
		servicesDir    string
		expectedResult map[string][]models.ContainerSummary
		expectedError  error
	}{
		{
			name:        "Success",
			servicesDir: "/services",
			expectedResult: map[string][]models.ContainerSummary{
				"service1": {
					{
						ID:     "container1",
						Name:   "container1",
						Image:  "image1",
						State:  container.StateRunning,
						Health: container.Healthy,
					},
				},
			},
			expectedError: nil,
		},
		{
			name:           "Error",
			servicesDir:    "/services",
			expectedResult: nil,
			expectedError:  errors.New("failed to get managed containers"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mocker := &Mocker{}
			service := newServiceWithMocks(t, mocker, models.DeploymentParams{
				ServicesDir: tc.servicesDir,
			})

			mocker.On("GetManagedStacks", tc.servicesDir).Return(tc.expectedResult, tc.expectedError)

			result, err := service.GetManagedStacks()

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}

func TestGetDiff_Success(t *testing.T) {
	mocker := &Mocker{}
	service := newServiceWithCurrentConfig(t, mocker, models.DeploymentParams{}, mockConfigOld)

	expectedDiff := []models.FileDiff{
		{OldFile: "file1.txt", NewFile: "file1.txt", Diff: "diff1"},
		{OldFile: "file2.txt", NewFile: "file2.txt", Diff: "diff2"},
	}

	mocker.On("WithConfig", mockConfigOld).Return(service.fetcher)
	mocker.On("DiffWithRemote").Return(git.Patch{Files: expectedDiff}, nil)

	diff, err := service.GetDiff()

	assert.NoError(t, err)
	assert.Equal(t, expectedDiff, diff)
	mocker.AssertExpectations(t)
}

func TestGetDiff_ErrorDiffWithRemote(t *testing.T) {
	mocker := &Mocker{}
	service := newServiceWithCurrentConfig(t, mocker, models.DeploymentParams{}, mockConfigOld)

	mocker.On("WithConfig", mockConfigOld).Return(service.fetcher)
	mocker.On("DiffWithRemote").Return(git.Patch{}, ErrFetch)

	diff, err := service.GetDiff()

	assert.Error(t, err)
	assert.Equal(t, 0, len(diff))
	mocker.AssertExpectations(t)
}
