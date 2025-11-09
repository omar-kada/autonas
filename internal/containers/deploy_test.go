package containers

import (
	"errors"
	"fmt"
	"omar-kada/autonas/internal/config"
	"omar-kada/autonas/internal/containers/model"
	"omar-kada/autonas/internal/logger"

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

func newDeployerWithMocks(mocker *Mocker) *Deployer {
	return &Deployer{
		log:               logger.New(true),
		containersManager: mocker,
		_copyFunc:         mocker.Copy,
	}
}

func TestDeployServices_Success(t *testing.T) {
	mocker := &Mocker{}
	deployer := newDeployerWithMocks(mocker)
	mock.InOrder(
		mocker.On(
			"RemoveServices", []string{"svc1"}, "/services",
		).Return(nil),

		mocker.On(
			"Copy", "configFolder/services/svc2", "/services/svc2", []copydir.Options(nil),
		).Return(nil),

		mocker.On(
			"Copy", "configFolder/services/svc3", "/services/svc3", []copydir.Options(nil),
		).Return(nil),

		mocker.On(
			"DeployServices", mockConfigNew,
		).Return(nil),
	)
	err := deployer.DeployServices("configFolder", mockConfigOld, mockConfigNew)
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
			deployer := newDeployerWithMocks(mocker)

			mock.InOrder(
				mocker.On(
					"RemoveServices", []string{"svc1"}, "/services",
				).Return(tc.errors.removeErr),

				mocker.On(
					"Copy", "configFolder/services/svc2", "/services/svc2", []copydir.Options(nil),
				).Return(tc.errors.copyErr),

				mocker.On(
					"Copy", "configFolder/services/svc3", "/services/svc3", []copydir.Options(nil),
				).Return(tc.errors.copyErr),

				mocker.On(
					"DeployServices", mockConfigNew,
				).Return(tc.errors.deployErr),
			)

			err := deployer.DeployServices("configFolder", mockConfigOld, mockConfigNew)

			assert.ErrorContains(t, err, fmt.Sprint(tc.expectedError))
		})
	}
}
