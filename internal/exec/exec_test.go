package exec

import (
	"errors"
	"omar-kada/autonas/internal/config"
	"omar-kada/autonas/internal/testutil"

	"testing"

	"github.com/moby/moby/api/types/container"
)

type mockContainersHandler struct {
	testutil.Mocker
	removeErr error
	deployErr error
}

func (m *mockContainersHandler) RemoveServices(services []string, servicesPath string) error {
	m.AddCall("removeServices", services, servicesPath)
	return m.removeErr
}

func (m *mockContainersHandler) DeployServices(cfg config.Config) error {
	m.AddCall("deployServices", cfg)
	return m.deployErr
}

func (m *mockContainersHandler) GetManagedContainers() (map[string][]container.Summary, error) {
	m.AddCall("getManagedContainers")
	return nil, nil
}

type mockFileManager struct {
	testutil.Mocker
	copyErr error
}

func (m *mockFileManager) CopyToPath(srcFolder, servicesPath string) error {
	m.AddCall("CopyToPath", srcFolder, servicesPath)
	return m.copyErr
}

func (m *mockFileManager) WriteToFile(filePath string, content string) error {
	m.AddCall("writeToFile", filePath, content)
	return nil
}

var mockConfigOld = config.Config{
	EnabledServices: []string{"svc1", "svc2"},
	ServicesPath:    "/services",
}
var mockConfigNew = config.Config{
	EnabledServices: []string{"svc2", "svc3"},
	ServicesPath:    "/services",
}

var (
	fileManager       = &mockFileManager{}
	containersHandler = &mockContainersHandler{}
)

func TestDeployServices_Success(t *testing.T) {
	defaultFileManager = fileManager
	defaultContainersHandler = containersHandler
	deployer := New()
	err := deployer.DeployServices("configFolder", mockConfigOld, mockConfigNew)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedCalls := [][]any{
		{"CopyToPath", "configFolder/services", "/services"},
	}
	fileManager.AssertCalls(t, expectedCalls)
	expectedContainerCalls := [][]any{
		{"removeServices", []string{"svc1"}, "/services"},
		{"deployServices", mockConfigNew},
	}
	containersHandler.AssertCalls(t, expectedContainerCalls)
}

var ErrRemove = errors.New("removeServices error")
var ErrDeploy = errors.New("deployServices error")
var ErrCopy = errors.New("copyServices error")

func TestDeployServices_Errors(t *testing.T) {
	testCases := []struct {
		name              string
		containersHandler mockContainersHandler
		fileManager       mockFileManager
		expectedError     error
	}{
		{
			name:              "removeServices error",
			containersHandler: mockContainersHandler{removeErr: ErrRemove},
			fileManager:       mockFileManager{},
			expectedError:     ErrRemove,
		},
		{
			name:              "deployServices error",
			containersHandler: mockContainersHandler{deployErr: ErrDeploy},
			fileManager:       mockFileManager{},
			expectedError:     ErrDeploy,
		},

		{
			name:              "copyServices error",
			containersHandler: mockContainersHandler{},
			fileManager:       mockFileManager{copyErr: ErrCopy},
			expectedError:     ErrCopy,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// create a mock that returns an error for the chosen method
			// replace the package default with the one returning an error
			defaultContainersHandler = &tc.containersHandler
			defaultFileManager = &tc.fileManager
			deployer := New()

			err := deployer.DeployServices("configFolder", mockConfigOld, mockConfigNew)
			if !errors.Is(err, tc.expectedError) {
				t.Fatalf("expected error %q, got %v", tc.name, err)
			}
		})
	}
}
