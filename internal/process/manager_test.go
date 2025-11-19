package process

import (
	"errors"
	"fmt"
	"omar-kada/autonas/internal/config"
	"omar-kada/autonas/models"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type Mocker struct {
	mock.Mock
	Manager
}

func (m *Mocker) RemoveServices(services []string, servicesDir string) error {
	args := m.Called(services, servicesDir)
	return args.Error(0)
}

func (m *Mocker) DeployServices(cfg config.Config, servicesDir string) error {
	args := m.Called(cfg, servicesDir)
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
		Services: map[string]config.ServiceConfig{
			"svc1": {},
			"svc2": {},
			"svc3": {Disabled: true},
		},
	}
	mockConfigNew = config.Config{
		Services: map[string]config.ServiceConfig{
			"svc1": {Disabled: true},
			"svc2": {},
			"svc3": {},
		},
	}
)

func newDeployerWithMocks(mocker *Mocker, params ManagerParams) *manager {
	return &manager{
		containersDeployer: mocker,
		copier:             mocker,
		fetcher:            mocker,
		configGenerator:    mocker,
		params:             params,
	}
}

func TestDeployServices_Success(t *testing.T) {
	mocker := &Mocker{}
	deployer := newDeployerWithMocks(mocker, ManagerParams{
		ServicesDir: "/services",
		WorkingDir:  "configDir",
		ConfigFile:  "config.yaml",
	})
	deployer.currentCfg = mockConfigOld
	mock.InOrder(
		mocker.On(
			"RemoveServices", []string{"svc1"}, "/services",
		).Return(nil),

		mocker.On(
			"DeployServices", mockConfigNew, "/services",
		).Return(nil),
	)
	mocker.On(
		"CopyWithAddPerm", "configDir/services/svc2", "/services/svc2", os.FileMode(0000),
	).Return(nil)

	mocker.On(
		"CopyWithAddPerm", "configDir/services/svc3", "/services/svc3", os.FileMode(0000),
	).Return(nil)
	err := deployer.removeAndDeployStacks(mockConfigOld, mockConfigNew)
	assert.NoError(t, err)
}

var (
	ErrRemove = errors.New("removeServices error")
	ErrDeploy = errors.New("deployServices error")
	ErrCopy   = errors.New("copyServices error")
)

type ExpectedErrors struct {
	removeErr error
	deployErr error
	copyErr   error
}

func TestDeployServices_Errors(t *testing.T) {
	testCases := []struct {
		name          string
		errors        ExpectedErrors
		expectedError error
	}{
		{
			name:          "removeServices error",
			errors:        ExpectedErrors{removeErr: ErrRemove},
			expectedError: ErrRemove,
		},
		{
			name:          "deployServices error",
			errors:        ExpectedErrors{deployErr: ErrDeploy},
			expectedError: ErrDeploy,
		},

		{
			name:          "copyServices error",
			errors:        ExpectedErrors{copyErr: ErrCopy},
			expectedError: ErrCopy,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mocker := &Mocker{}
			deployer := newDeployerWithMocks(mocker, ManagerParams{
				ServicesDir: "/services",
				WorkingDir:  "configDir",
				ConfigFile:  "config.yaml",
			})
			deployer.currentCfg = mockConfigOld

			mock.InOrder(
				mocker.On(
					"RemoveServices", []string{"svc1"}, "/services",
				).Return(tc.errors.removeErr),

				mocker.On(
					"DeployServices", mockConfigNew, "/services",
				).Return(tc.errors.deployErr),
			)

			mocker.On(
				"CopyWithAddPerm", "configDir/services/svc2", "/services/svc2", os.FileMode(0000),
			).Return(tc.errors.copyErr)

			mocker.On(
				"CopyWithAddPerm", "configDir/services/svc3", "/services/svc3", os.FileMode(0000),
			).Return(tc.errors.copyErr)
			err := deployer.removeAndDeployStacks(mockConfigOld, mockConfigNew)

			assert.ErrorContains(t, err, fmt.Sprint(tc.expectedError))
		})
	}
}
