package process

import (
	"errors"
	"log/slog"
	"omar-kada/autonas/internal/config"
	"omar-kada/autonas/models"
	"os"
	"testing"
	"time"

	"github.com/moby/moby/api/types/container"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type Mocker struct {
	mock.Mock
}

func (m *Mocker) WithLogger(log *slog.Logger) Deployer {
	args := m.Called(log)
	return args.Get(0).(Deployer)
}

func (m *Mocker) RemoveServices(services []string, servicesDir string) map[string]error {
	args := m.Called(services, servicesDir)
	return args.Get(0).(map[string]error)
}

func (m *Mocker) DeployServices(cfg config.Config, servicesDir string) map[string]error {
	args := m.Called(cfg, servicesDir)
	return args.Get(0).(map[string]error)
}

func (m *Mocker) RemoveAndDeployStacks(oldCfg, cfg config.Config, params models.DeploymentParams) error {
	args := m.Called(oldCfg, cfg, params)
	return args.Error(0)
}

func (m *Mocker) Copy(srcDir, servicesDir string) error {
	args := m.Called(srcDir, servicesDir)
	return args.Error(0)
}

func (m *Mocker) CopyWithAddPerm(srcDir, servicesDir string, permission os.FileMode) error {
	args := m.Called(srcDir, servicesDir, permission)
	return args.Error(0)
}

func (m *Mocker) GetManagedContainers(servicesDir string) (map[string][]models.ContainerSummary, error) {
	args := m.Called(servicesDir)
	return args.Get(0).(map[string][]models.ContainerSummary), args.Error(1)
}

func (m *Mocker) Fetch(repo string, branch string, dir string) error {
	args := m.Called(repo, branch, dir)
	return args.Error(0)
}

func (m *Mocker) FromFiles(files []string) (config.Config, error) {
	args := m.Called(files)
	return args.Get(0).(config.Config), args.Error(1)
}

var (
	mockConfigOld = config.Config{
		Repo:   "https://example.com/repo.git",
		Branch: "main",
		Services: map[string]config.ServiceConfig{
			"svc1": {
				Extra: map[string]any{
					"Port":    8080,
					"Version": "v1",
				},
			},
			"svc2": {},
			"svc3": {Disabled: true},
		},
	}
	mockConfigNew = config.Config{
		Repo:   "https://example.com/repo.git",
		Branch: "main",
		Services: map[string]config.ServiceConfig{
			"svc1": {
				Disabled: true,
				Extra: map[string]any{
					"Port":    8080,
					"Version": "v1",
				},
			},
			"svc2": {},
			"svc3": {},
		},
	}
)

func newManagerWithCurrentConfig(mocker *Mocker, params models.DeploymentParams, currentCfg config.Config) *manager {
	return &manager{
		containersDeployer:  mocker,
		containersInspector: mocker,
		copier:              mocker,
		fetcher:             mocker,
		configGenerator:     mocker,
		params:              params,
		currentCfg:          currentCfg,
	}
}

func newManagerWithMocks(mocker *Mocker, params models.DeploymentParams) *manager {
	return newManagerWithCurrentConfig(mocker, params, config.Config{})
}

var (
	ErrRemove   = errors.New("removeServices error")
	ErrDeploy   = errors.New("deployServices error")
	ErrCopy     = errors.New("copyServices error")
	ErrGenerate = errors.New("generate file error")
	ErrFetch    = errors.New("sync config error")
)

func TestSync_Success(t *testing.T) {
	mocker := &Mocker{}
	manager := newManagerWithMocks(mocker, models.DeploymentParams{
		ServicesDir: "/services",
		WorkingDir:  ".",
		ConfigFile:  "config.yaml",
	})

	wantCfg := mockConfigOld

	mocker.On("FromFiles", []string{"config.yaml"}).Once().Return(wantCfg, nil)
	mocker.On("Fetch", wantCfg.Repo, wantCfg.Branch, ".").Once().Return(nil)
	mocker.On("CopyWithAddPerm", "services/svc1", "/services/svc1", os.FileMode(0000)).Once().Return(nil)
	mocker.On("CopyWithAddPerm", "services/svc2", "/services/svc2", os.FileMode(0000)).Once().Return(nil)
	mocker.On("DeployServices", wantCfg, "/services").Once().Return(map[string]error{})
	err := manager.SyncDeployment()
	assert.NoError(t, err)
	mocker.AssertExpectations(t)
}

func TestSync_Success_RedploymentWithChangedConfig(t *testing.T) {
	mocker := &Mocker{}
	manager := newManagerWithCurrentConfig(mocker, models.DeploymentParams{
		ServicesDir: "/services",
		WorkingDir:  ".",
		ConfigFile:  "config.yaml",
	}, mockConfigOld)

	wantCfg := mockConfigNew
	manager.params.AddWritePerm = true
	mocker.On("FromFiles", []string{"config.yaml"}).Once().Return(wantCfg, nil)
	mocker.On("Fetch", wantCfg.Repo, wantCfg.Branch, ".").Once().Return(nil)
	mocker.On("CopyWithAddPerm", "services/svc2", "/services/svc2", os.FileMode(0666)).Once().Return(nil)
	mocker.On("CopyWithAddPerm", "services/svc3", "/services/svc3", os.FileMode(0666)).Once().Return(nil)
	mocker.On("RemoveServices", []string{"svc1"}, "/services").Once().Return(map[string]error{})
	mocker.On("DeployServices", wantCfg, "/services").Once().Return(map[string]error{})

	err := manager.SyncDeployment()
	assert.NoError(t, err)
	mocker.AssertExpectations(t)
}

func TestRunCmd_Errors(t *testing.T) {
	type ExpectedErrors struct {
		fetchErr    error
		generateErr error
		deployErr   map[string]error
		removeErr   map[string]error
		copyErr     error
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
			name:          "generateConfigFromFiles error",
			mockValues:    ExpectedErrors{generateErr: ErrGenerate},
			expectedError: ErrGenerate,
		},
		{
			name:          "deployServices error",
			mockValues:    ExpectedErrors{deployErr: map[string]error{"svc1": ErrDeploy}},
			expectedError: ErrDeploy,
		},
		{
			name:          "removeServices error",
			mockValues:    ExpectedErrors{removeErr: map[string]error{"svc1": ErrRemove}},
			expectedError: ErrRemove,
		},
		{
			name:          "copy files error",
			mockValues:    ExpectedErrors{copyErr: ErrCopy},
			expectedError: ErrCopy,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mocker := &Mocker{}
			manager := newManagerWithCurrentConfig(mocker, models.DeploymentParams{
				ServicesDir: "/services",
				WorkingDir:  ".",
				ConfigFile:  "config.yaml",
			}, mockConfigOld)
			wantCfg := mockConfigNew
			mocker.On("FromFiles", []string{"config.yaml"}).Once().Return(wantCfg, tc.mockValues.generateErr)
			mocker.On("Fetch", wantCfg.Repo, wantCfg.Branch, ".").Once().Return(tc.mockValues.fetchErr)
			mocker.On("CopyWithAddPerm", "services/svc2", "/services/svc2", os.FileMode(0000)).Once().Return(tc.mockValues.copyErr)
			mocker.On("CopyWithAddPerm", "services/svc3", "/services/svc3", os.FileMode(0000)).Once().Return(tc.mockValues.copyErr)
			mocker.On("RemoveServices", []string{"svc1"}, "/services").Once().Return(tc.mockValues.removeErr)
			mocker.On("DeployServices", wantCfg, "/services").Once().Return(tc.mockValues.deployErr)

			err := manager.SyncDeployment()
			assert.ErrorContains(t, err, tc.expectedError.Error())
		})
	}
}

func TestSyncPeriodically(t *testing.T) {
	testCases := []struct {
		name          string
		currentCfg    config.Config
		expectedCron  bool
		expectedError error
	}{
		{
			name: "CronPeriod is set",
			currentCfg: config.Config{
				CronPeriod: "*/5 * * * *",
			},
			expectedCron:  true,
			expectedError: nil,
		},
		{
			name: "CronPeriod is not set",
			currentCfg: config.Config{
				CronPeriod: "",
			},
			expectedCron:  false,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mocker := &Mocker{}
			manager := newManagerWithCurrentConfig(mocker, models.DeploymentParams{}, tc.currentCfg)

			cron := manager.SyncPeriodically()

			if tc.expectedCron {
				assert.NotNil(t, cron)
				assert.Equal(t, 1, len(cron.Entries()))
				assert.Less(t, time.Now(), cron.Entries()[0].Next)
			} else {
				assert.Nil(t, cron)
			}
		})
	}
}

func TestGetManagedContainers(t *testing.T) {
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
			manager := newManagerWithMocks(mocker, models.DeploymentParams{
				ServicesDir: tc.servicesDir,
			})

			mocker.On("GetManagedContainers", tc.servicesDir).Return(tc.expectedResult, tc.expectedError)

			result, err := manager.GetManagedContainers()

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
