package process

import (
	"context"
	"errors"
	"omar-kada/autonas/internal/events"
	"omar-kada/autonas/internal/storage"
	"omar-kada/autonas/models"
	"testing"

	"github.com/moby/moby/api/types/container"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type Mocker struct {
	mock.Mock
}

func (m *Mocker) WithCtx(_ context.Context) Deployer {
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

func (m *Mocker) Fetch(repo string, branch string, dir string) error {
	args := m.Called(repo, branch, dir)
	return args.Error(0)
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

func newServiceWithCurrentConfig(mocker *Mocker, params models.DeploymentParams, currentCfg models.Config) *service {
	return &service{
		containersDeployer:  mocker,
		containersInspector: mocker,
		fetcher:             mocker,
		dispatcher:          events.NewVoidDispatcher(),
		store:               storage.NewMemoryStorage(),
		params:              params,
		currentCfg:          currentCfg,
	}
}

func newServiceWithMocks(mocker *Mocker, params models.DeploymentParams) *service {
	return newServiceWithCurrentConfig(mocker, params, models.Config{})
}

var (
	ErrRemove   = errors.New("removeServices error")
	ErrDeploy   = errors.New("deployServices error")
	ErrGenerate = errors.New("generate file error")
	ErrFetch    = errors.New("sync config error")
)

func TestSync_Success(t *testing.T) {
	mocker := &Mocker{}
	service := newServiceWithMocks(mocker, models.DeploymentParams{
		ServicesDir: "/services",
		WorkingDir:  ".",
	})

	wantCfg := mockConfigOld

	mocker.On("Fetch", wantCfg.Repo, wantCfg.Branch, ".").Once().Return(nil)
	mocker.On("RemoveAndDeployStacks", models.Config{}, wantCfg, service.params).Once().Return(nil)
	err := service.SyncDeployment(wantCfg)
	assert.NoError(t, err)
	mocker.AssertExpectations(t)
}

func TestSync_Success_RedploymentWithChangedConfig(t *testing.T) {
	mocker := &Mocker{}
	service := newServiceWithCurrentConfig(mocker, models.DeploymentParams{
		ServicesDir: "/services",
		WorkingDir:  ".",
	}, mockConfigOld)

	wantCfg := mockConfigNew
	service.params.AddWritePerm = true
	mocker.On("Fetch", wantCfg.Repo, wantCfg.Branch, ".").Once().Return(nil)
	mocker.On("RemoveAndDeployStacks", mockConfigOld, wantCfg, service.params).Once().Return(nil)

	err := service.SyncDeployment(wantCfg)
	assert.NoError(t, err)
	mocker.AssertExpectations(t)
}

func TestRunCmd_Errors(t *testing.T) {
	type ExpectedErrors struct {
		fetchErr    error
		generateErr error
		deployErr   error
	}
	testCases := []struct {
		name          string
		mockValues    ExpectedErrors
		expectedError error
	}{
		{
			name:          "fetch from repo error",
			mockValues:    ExpectedErrors{fetchErr: ErrFetch},
			expectedError: ErrFetch,
		},
		{
			name:          "deployServices error",
			mockValues:    ExpectedErrors{deployErr: ErrDeploy},
			expectedError: ErrDeploy,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mocker := &Mocker{}
			service := newServiceWithCurrentConfig(mocker, models.DeploymentParams{
				ServicesDir: "/services",
				WorkingDir:  ".",
			}, mockConfigOld)
			wantCfg := mockConfigNew
			mocker.On("FromFiles", []string{"config.yaml"}).Once().Return(wantCfg, tc.mockValues.generateErr)
			mocker.On("Fetch", wantCfg.Repo, wantCfg.Branch, ".").Once().Return(tc.mockValues.fetchErr)
			mocker.On("WithLogger", mock.Anything).Once().Return(mocker)
			mocker.On("RemoveAndDeployStacks", mockConfigOld, wantCfg, service.params).Once().Return(tc.mockValues.deployErr)

			err := service.SyncDeployment(wantCfg)
			assert.ErrorContains(t, err, tc.expectedError.Error())
		})
	}
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
			service := newServiceWithMocks(mocker, models.DeploymentParams{
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
