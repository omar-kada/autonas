package run

import (
	"errors"
	"omar-kada/autonas/internal/config"
	"omar-kada/autonas/internal/logger"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var mockConfig = config.Config{
	Extra: map[string]any{
		"AutonasHost": "localhost",
		"ServicesDir": "/services",
		"DataPath":    "/data",
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
}

func (m *Mocker) FromFiles(files []string) (config.Config, error) {
	args := m.Called(files)
	return args.Get(0).(config.Config), args.Error(1)
}

func (m *Mocker) Sync(repoURL string, branch string, path string) error {
	args := m.Called(repoURL, branch, path)
	return args.Error(0)
}

func (m *Mocker) DeployServices(configDir, servicesDir string, currentCfg, cfg config.Config) error {
	args := m.Called(configDir, servicesDir, currentCfg, cfg)
	return args.Error(0)
}

func (m *Mocker) AddPermission(perm os.FileMode) {
	m.Called(perm)
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
			"Sync", "https://example.com/repo.git", "main", ".",
		).Once().Return(val.syncErr),
		m.On(
			"FromFiles", []string{"config1.yaml", "config2.yaml"},
		).Once().Return(val.generateConfig, val.generateErr),
		m.On(
			"DeployServices", ".", "/services", val.deployInputOldConfig, val.generateConfig,
		).Once().Return(val.deployErr),
	)
}

func newRunnerWithMocks(mocker *Mocker) *Cmd {
	return &Cmd{
		Log:             logger.New(true),
		Deployer:        mocker,
		ConfigGenerator: mocker,
		Syncer:          mocker,
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

	err := runner.RunOnce(runParams{
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

	err = runner.RunOnce(
		runParams{
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

			err := runner.RunOnce(runParams{
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

func TestGetParamsWithDefaults_AllCliValuesProvided(t *testing.T) {
	// When all CLI values are provided, they should be returned as-is
	params := runParams{
		ConfigFiles: []string{"custom.yaml"},
		Repo:        "https://custom.com/repo.git",
		Branch:      "develop",
		WorkingDir:  "/custom/work",
		ServicesDir: "/custom/services",
		CronPeriod:  "0 * * * *",
	}

	result := getParamsWithDefaults(params)

	assert.Equal(t, []string{"custom.yaml"}, result.ConfigFiles)
	assert.Equal(t, "https://custom.com/repo.git", result.Repo)
	assert.Equal(t, "develop", result.Branch)
	assert.Equal(t, "/custom/work", result.WorkingDir)
	assert.Equal(t, "/custom/services", result.ServicesDir)
	assert.Equal(t, "0 * * * *", result.CronPeriod)
}

func TestGetParamsWithDefaults_UseEnvVariablesWhenCliEmpty(t *testing.T) {
	t.Setenv("AUTONAS_CONFIG_REPO", "https://env.com/repo.git")
	t.Setenv("AUTONAS_CONFIG_BRANCH", "env-branch")
	t.Setenv("AUTONAS_WORKING_DIR", "/env/work")
	t.Setenv("AUTONAS_SERVICES_DIR", "/env/services")
	t.Setenv("AUTONAS_CRON_PERIOD", "*/5 * * * *")
	t.Setenv("AUTONAS_CONFIG_FILES", "env1.yaml,env2.yaml")

	params := runParams{}

	result := getParamsWithDefaults(params)

	assert.Equal(t, []string{"env1.yaml", "env2.yaml"}, result.ConfigFiles)
	assert.Equal(t, "https://env.com/repo.git", result.Repo)
	assert.Equal(t, "env-branch", result.Branch)
	assert.Equal(t, "/env/work", result.WorkingDir)
	assert.Equal(t, "/env/services", result.ServicesDir)
	assert.Equal(t, "*/5 * * * *", result.CronPeriod)
}

func TestGetParamsWithDefaults_UseDefaultsWhenCliAndEnvEmpty(t *testing.T) {
	// Clear environment variables
	t.Setenv("AUTONAS_CONFIG_REPO", "")
	t.Setenv("AUTONAS_CONFIG_BRANCH", "")
	t.Setenv("AUTONAS_WORKING_DIR", "")
	t.Setenv("AUTONAS_SERVICES_DIR", "")
	t.Setenv("AUTONAS_CRON_PERIOD", "")
	t.Setenv("AUTONAS_CONFIG_FILES", "")

	params := runParams{}

	result := getParamsWithDefaults(params)

	// Check defaults are applied
	assert.Equal(t, []string{"config.yaml"}, result.ConfigFiles)
	assert.Equal(t, "main", result.Branch)
	assert.Equal(t, "./config", result.WorkingDir)
	assert.Equal(t, ".", result.ServicesDir)
	// Repo and CronPeriod have nil defaults, so they should be empty strings
	assert.Equal(t, "", result.Repo)
	assert.Equal(t, "", result.CronPeriod)
}

func TestGetParamsWithDefaults_CliPriority(t *testing.T) {
	// CLI values should take priority over env variables and defaults
	t.Setenv("AUTONAS_CONFIG_BRANCH", "env-branch")

	params := runParams{
		Branch: "cli-branch",
	}

	result := getParamsWithDefaults(params)

	// CLI value should win
	assert.Equal(t, "cli-branch", result.Branch)
}

func TestGetParamsWithDefaults_MixedSources(t *testing.T) {
	// Test a mix of CLI values, env variables, and defaults
	t.Setenv("AUTONAS_CONFIG_BRANCH", "env-branch")
	t.Setenv("AUTONAS_WORKING_DIR", "")

	params := runParams{
		ConfigFiles: []string{"cli.yaml"}, // From CLI
		Branch:      "cli-branch",         // From CLI (overrides env)
		// WorkingDir not provided, should use env or default
	}

	result := getParamsWithDefaults(params)

	assert.Equal(t, []string{"cli.yaml"}, result.ConfigFiles)
	assert.Equal(t, "cli-branch", result.Branch)
	assert.Equal(t, "./config", result.WorkingDir) // Should use default
}
