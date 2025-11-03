package cli

import (
	"errors"
	"omar-kada/autonas/internal/config"
	"omar-kada/autonas/internal/testutil"
	"testing"
)

var mockConfig = config.Config{
	AutonasHost:     "localhost",
	ServicesPath:    "/services",
	DataPath:        "/data",
	EnabledServices: []string{"svc1"},
	Services: map[string]config.ServiceConfig{
		"svc1": {
			Port:    8080,
			Version: "v1",
		},
	},
}

type Mocker struct {
	testutil.MockRecorder
	generateErr error
	syncErr     error
	deployErr   error
}

func (m *Mocker) generateConfigFromFiles(files []string) (config.Config, error) {
	m.AddCall("generateConfigFromFiles", files)
	return mockConfig, m.generateErr
}

func (m *Mocker) syncCode(repoURL string, branch string, path string) error {
	m.AddCall("syncCode", repoURL, branch, path)
	return m.syncErr
}

func (m *Mocker) DeployServices(configFolder string, currentCfg config.Config, cfg config.Config) error {
	m.AddCall("deployServices", configFolder, currentCfg, cfg)
	return m.deployErr
}

func initMocks(useMocker *Mocker) *Mocker {
	generateConfigFromFiles = useMocker.generateConfigFromFiles
	syncCode = useMocker.syncCode
	defaultDeployer = useMocker
	return useMocker
}

func TestRunCmd_Success(t *testing.T) {
	mocker := initMocks(&Mocker{})
	run := NewRunner()

	err := run.RunCmd([]string{"config1.yaml", "config2.yaml"}, "https://example.com/repo.git")
	if err != nil {
		t.Fatalf("RunCmd failed: %v", err)
	}

	wantCfg := mockConfig

	wantCalls := [][]any{
		{"generateConfigFromFiles", []string{"config1.yaml", "config2.yaml"}},
		{"syncCode", "https://example.com/repo.git", "main", "."},
		{"deployServices", ".", config.Config{}, wantCfg},
	}
	mocker.AssertCalls(t, wantCalls)

	// Verify that the currentCfg in runner is updated
	mocker.Reset()

	err = run.RunCmd([]string{"config1.yaml", "config2.yaml"}, "https://example.com/repo.git")
	if err != nil {
		t.Fatalf("RunCmd failed: %v", err)
	}

	wantCalls = [][]any{
		{"generateConfigFromFiles", []string{"config1.yaml", "config2.yaml"}},
		{"syncCode", "https://example.com/repo.git", "main", "."},
		{"deployServices", ".", wantCfg, wantCfg},
	}
	mocker.AssertCalls(t, wantCalls)
}

var (
	ErrGenerate = errors.New("generate file error")
	ErrSync     = errors.New("sync config error")
	ErrDeploy   = errors.New("deploy error")
)

func TestRunCmd_Errors(t *testing.T) {

	testCases := []struct {
		name          string
		errorMocker   Mocker
		expectedError error
	}{
		{
			name:          "syncCode error",
			errorMocker:   Mocker{syncErr: ErrSync},
			expectedError: ErrSync,
		},
		{
			name:          "generateConfigFromFiles error",
			errorMocker:   Mocker{generateErr: ErrGenerate},
			expectedError: ErrGenerate,
		},
		{
			name:          "deployServices error",
			errorMocker:   Mocker{deployErr: ErrDeploy},
			expectedError: ErrDeploy,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			initMocks(&tc.errorMocker)
			runner := NewRunner()

			err := runner.RunCmd([]string{"config1.yaml"}, "https://example.com/repo.git")
			if !errors.Is(err, tc.expectedError) {
				t.Fatalf("expected error %q, got %v", tc.expectedError, err)
			}
		})
	}
}
