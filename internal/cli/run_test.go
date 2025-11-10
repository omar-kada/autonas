package cli

import (
	"errors"
	"omar-kada/autonas/internal/config"
	"omar-kada/autonas/internal/logger"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var mockConfig = config.Config{
	Extra: map[string]any{

		"AutonasHost":  "localhost",
		"ServicesPath": "/services",
		"DataPath":     "/data",
	},
	EnabledServices: []string{"svc1"},
	Services: map[string]config.ServiceConfig{
		"svc1": {
			Port:    8080,
			Version: "v1",
		},
	},
}

type Mocker struct {
	mock.Mock
	Deployer
}

func (m *Mocker) generateConfigFromFiles(files []string) (config.Config, error) {
	args := m.Called(files)
	return args.Get(0).(config.Config), args.Error(1)
}

func (m *Mocker) syncCode(repoURL string, branch string, path string) error {
	args := m.Called(repoURL, branch, path)
	return args.Error(0)
}

func (m *Mocker) DeployServices(configFolder, servicesDir string, currentCfg, cfg config.Config) error {
	args := m.Called(configFolder, servicesDir, currentCfg, cfg)
	return args.Error(0)
}

type ExpectedValues struct {
	generateConfig       config.Config
	generateErr          error
	syncErr              error
	deployErr            error
	deployInputOldConfig config.Config
}

func mockReturnValues(m *Mocker, val ExpectedValues) {
	mock.InOrder(
		m.On(
			"syncCode", "https://example.com/repo.git", "main", ".",
		).Once().Return(val.syncErr),
		m.On(
			"generateConfigFromFiles", []string{"config1.yaml", "config2.yaml"},
		).Once().Return(val.generateConfig, val.generateErr),
		m.On(
			"DeployServices", ".", "/services", val.deployInputOldConfig, val.generateConfig,
		).Once().Return(val.deployErr),
	)
}

func newRunnerWithMocks(mocker *Mocker) *Runner {
	return &Runner{
		log:                      logger.New(true),
		deployer:                 mocker,
		_generateConfigFromFiles: mocker.generateConfigFromFiles,
		_syncCode:                mocker.syncCode,
	}
}

func TestRunCmd_Success(t *testing.T) {
	mocker := &Mocker{}
	runner := newRunnerWithMocks(mocker)

	wantCfg := mockConfig
	mockReturnValues(mocker, ExpectedValues{
		generateConfig:       wantCfg,
		deployInputOldConfig: config.Config{},
	})

	err := runner.RunCmd(RunParams{
		ConfigFiles: []string{"config1.yaml", "config2.yaml"},
		Repo:        "https://example.com/repo.git",
		Branch:      "main",
		WorkingDir:  ".",
		ServicesDir: "/services",
	})
	assert.NoError(t, err)

	// Verify that the currentCfg in runner is updated
	mockReturnValues(mocker, ExpectedValues{
		generateConfig:       wantCfg,
		deployInputOldConfig: wantCfg,
	})

	err = runner.RunCmd(
		RunParams{
			ConfigFiles: []string{"config1.yaml", "config2.yaml"},
			Repo:        "https://example.com/repo.git",
			Branch:      "main",
			WorkingDir:  ".",
			ServicesDir: "/services",
		})
	assert.NoError(t, err)
}

var (
	ErrGenerate = errors.New("generate file error")
	ErrSync     = errors.New("sync config error")
	ErrDeploy   = errors.New("deploy error")
)

func TestRunCmd_Errors(t *testing.T) {

	testCases := []struct {
		name          string
		mockValues    ExpectedValues
		expectedError error
	}{
		{
			name:          "syncCode error",
			mockValues:    ExpectedValues{syncErr: ErrSync},
			expectedError: ErrSync,
		},
		{
			name:          "generateConfigFromFiles error",
			mockValues:    ExpectedValues{generateErr: ErrGenerate},
			expectedError: ErrGenerate,
		},
		{
			name:          "deployServices error",
			mockValues:    ExpectedValues{deployErr: ErrDeploy},
			expectedError: ErrDeploy,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mocker := &Mocker{}
			runner := newRunnerWithMocks(mocker)
			mockReturnValues(mocker, tc.mockValues)

			err := runner.RunCmd(RunParams{
				ConfigFiles: []string{"config1.yaml", "config2.yaml"},
				Repo:        "https://example.com/repo.git",
				Branch:      "main",
				WorkingDir:  ".",
				ServicesDir: "/services",
			})
			assert.ErrorIs(t, err, tc.expectedError, "want %s but got %s", tc.expectedError, err)
		})
	}
}
