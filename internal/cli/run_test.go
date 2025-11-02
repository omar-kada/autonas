package cli

import (
	"fmt"
	"omar-kada/autonas/internal/config"
	"reflect"
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
	Calls       [][]any
	ShouldError bool
}

func (m *Mocker) generateConfigFromFiles(files []string) (config.Config, error) {
	args := []any{"generateConfigFromFiles"}
	for _, f := range files {
		args = append(args, f)
	}
	m.Calls = append(m.Calls, args)
	if m.ShouldError {
		return config.Config{}, fmt.Errorf("mock error generateConfigFromFiles")
	}
	return mockConfig, nil
}

func (m *Mocker) syncCode(repoURL string, branch string, path string) error {
	m.Calls = append(m.Calls, []any{"syncCode", repoURL, branch, path})
	if m.ShouldError {
		return fmt.Errorf("mock error syncCode")
	}
	return nil
}

func (m *Mocker) DeployServices(configFolder string, currentCfg config.Config, cfg config.Config) error {
	m.Calls = append(m.Calls, []any{"deployServices", configFolder, currentCfg, cfg})
	if m.ShouldError {
		return fmt.Errorf("mock error deployServices")
	}
	return nil
}

var mocks Mocker
var errorMocks Mocker

func initMock() {
	mocks = Mocker{}
	errorMocks = Mocker{ShouldError: true}
	generateConfigFromFiles = mocks.generateConfigFromFiles
	syncCode = mocks.syncCode
	defaultDeployer = &mocks
}

func TestRunCmd_Success(t *testing.T) {
	initMock()
	runner := NewRunner()

	err := runner.RunCmd([]string{"config1.yaml", "config2.yaml"}, "https://example.com/repo.git")
	if err != nil {
		t.Fatalf("RunCmd failed: %v", err)
	}

	wantCfg := mockConfig

	currentCfg := runner.GetCurrentConfig()
	if !reflect.DeepEqual(currentCfg, wantCfg) {
		t.Fatalf("currentCfg mismatch\nwant=%+v\ngot =%+v", wantCfg, currentCfg)
	}
	wantCalls := [][]any{
		{"syncCode", "https://example.com/repo.git", "main", "."},
		{"generateConfigFromFiles", "config1.yaml", "config2.yaml"},
		{"deployServices", ".", config.Config{}, wantCfg},
	}

	if !reflect.DeepEqual(mocks.Calls, wantCalls) {
		t.Fatalf("Calls mismatch\nwant=%+v\ngot =%+v", wantCalls, mocks.Calls)
	}
}

func TestRunCmd_Errors(t *testing.T) {

	testCases := []struct {
		name          string
		mockErrors    func() // which function to mock to return error
		expectedError string
	}{
		{
			name:          "syncCode error",
			mockErrors:    func() { syncCode = errorMocks.syncCode },
			expectedError: "mock error syncCode",
		},
		{
			name:          "generateConfigFromFiles error",
			mockErrors:    func() { generateConfigFromFiles = errorMocks.generateConfigFromFiles },
			expectedError: "mock error generateConfigFromFiles",
		},
		{
			name:          "deployServices error",
			mockErrors:    func() { defaultDeployer = &errorMocks },
			expectedError: "mock error deployServices",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			initMock()

			// Set the specific function to return error
			tc.mockErrors()
			runner := NewRunner()

			err := runner.RunCmd([]string{"config1.yaml"}, "https://example.com/repo.git")
			if err == nil || err.Error() != tc.expectedError {
				t.Fatalf("expected error %q, got %v", tc.expectedError, err)
			}
		})
	}
}
