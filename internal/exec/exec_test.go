package exec

import (
	"errors"
	"omar-kada/autonas/internal/config"
	"omar-kada/autonas/internal/testutil"

	"testing"

	"github.com/moby/moby/api/types/container"
	copydir "github.com/otiai10/copy"
)

type Mocker struct {
	testutil.MockRecorder
	removeErr error
	deployErr error
	copyErr   error
}

func (m *Mocker) RemoveServices(services []string, servicesPath string) error {
	m.AddCall("removeServices", services, servicesPath)
	return m.removeErr
}

func (m *Mocker) DeployServices(cfg config.Config) error {
	m.AddCall("deployServices", cfg)
	return m.deployErr
}

func (m *Mocker) GetManagedContainers() (map[string][]container.Summary, error) {
	m.AddCall("getManagedContainers")
	return nil, nil
}

func (m *Mocker) Copy(srcFolder, servicesPath string, _ ...copydir.Options) error {
	m.AddCall("Copy", srcFolder, servicesPath)
	return m.copyErr
}

var (
	mockConfigOld = config.Config{
		EnabledServices: []string{"svc1", "svc2"},
		ServicesPath:    "/services",
	}
	mockConfigNew = config.Config{
		EnabledServices: []string{"svc2", "svc3"},
		ServicesPath:    "/services",
	}
)

func initMocks(useMocker *Mocker) *Mocker {
	defaultContainersHandler = useMocker
	copyFunc = useMocker.Copy
	return useMocker
}

func TestDeployServices_Success(t *testing.T) {
	mocker := initMocks(&Mocker{})
	deployer := New()
	err := deployer.DeployServices("configFolder", mockConfigOld, mockConfigNew)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	expectedCalls := [][]any{
		{"removeServices", []string{"svc1"}, "/services"},
		{"Copy", "configFolder/services", "/services"},
		{"deployServices", mockConfigNew},
	}
	mocker.AssertCalls(t, expectedCalls)
}

var (
	ErrRemove = errors.New("removeServices error")
	ErrDeploy = errors.New("deployServices error")
	ErrCopy   = errors.New("copyServices error")
)

func TestDeployServices_Errors(t *testing.T) {
	testCases := []struct {
		name          string
		mocker        Mocker
		expectedError error
	}{
		{
			name:          "removeServices error",
			mocker:        Mocker{removeErr: ErrRemove},
			expectedError: ErrRemove,
		},
		{
			name:          "deployServices error",
			mocker:        Mocker{deployErr: ErrDeploy},
			expectedError: ErrDeploy,
		},

		{
			name:          "copyServices error",
			mocker:        Mocker{copyErr: ErrCopy},
			expectedError: ErrCopy,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// create a mock that returns an error for the chosen method
			// replace the package default with the one returning an error
			initMocks(&tc.mocker)
			deployer := New()

			err := deployer.DeployServices("configFolder", mockConfigOld, mockConfigNew)
			if !errors.Is(err, tc.expectedError) {
				t.Fatalf("expected error %q, got %v", tc.name, err)
			}
		})
	}
}
