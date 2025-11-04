package containers

import (
	"errors"
	"omar-kada/autonas/internal/config"
	"omar-kada/autonas/internal/containers/model"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	copydir "github.com/otiai10/copy"
)

type Mocker struct {
	mock.Mock
	model.Manager
}

func (m *Mocker) RemoveServices(services []string, servicesPath string) error {
	args := m.Called(services, servicesPath)
	return args.Error(0)
}

func (m *Mocker) DeployServices(cfg config.Config) error {
	args := m.Called(cfg)
	return args.Error(0)
}

func (m *Mocker) Copy(srcFolder, servicesPath string, opts ...copydir.Options) error {
	args := m.Called(srcFolder, servicesPath, opts)
	return args.Error(0)
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

type ExpectedErrors struct {
	removeErr error
	deployErr error
	copyErr   error
}

func initMocker(errors ExpectedErrors) *Mocker {
	mocker := &Mocker{}
	mock.InOrder(
		mocker.On(
			"RemoveServices", []string{"svc1"}, "/services",
		).Return(errors.removeErr),

		mocker.On(
			"Copy", "configFolder/services", "/services", []copydir.Options(nil),
		).Return(errors.copyErr),

		mocker.On(
			"DeployServices", mockConfigNew,
		).Return(errors.deployErr),
	)
	return mocker
}

func newDeployerWithMocks(mocker *Mocker) *defaultDeployer {
	return &defaultDeployer{
		containersManager: mocker,
		_copyFunc:         mocker.Copy,
	}
}

func TestDeployServices_Success(t *testing.T) {
	mocker := initMocker(ExpectedErrors{})
	deployer := newDeployerWithMocks(mocker)

	err := deployer.DeployServices("configFolder", mockConfigOld, mockConfigNew)
	assert.NoError(t, err)
}

var (
	ErrRemove = errors.New("removeServices error")
	ErrDeploy = errors.New("deployServices error")
	ErrCopy   = errors.New("copyServices error")
)

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

			mocker := initMocker(tc.errors)
			deployer := newDeployerWithMocks(mocker)

			err := deployer.DeployServices("configFolder", mockConfigOld, mockConfigNew)

			assert.ErrorIs(t, err, tc.expectedError, "want %s but got %s", tc.expectedError, err)
		})
	}
}
