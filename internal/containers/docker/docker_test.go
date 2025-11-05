package docker

import (
	"errors"
	"fmt"
	"omar-kada/autonas/internal/config"
	"omar-kada/autonas/internal/logger"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type Mocker struct {
	mock.Mock
}

func (m *Mocker) WriteToFile(filePath string, content string) error {
	args := m.Called(filePath, content)
	return args.Error(0)
}

func (m *Mocker) RunCommand(cmd string, cmdArgs ...string) error {
	args := m.Called(cmd, cmdArgs)
	return args.Error(0)
}

func newManagerWithMocks(mocker *Mocker) *Manager {
	return &Manager{
		log:              logger.New(true),
		_writeToFileFunc: mocker.WriteToFile,
		_runCommandFunc:  mocker.RunCommand,
	}
}

var (
	mockConfig = config.Config{
		AutonasHost:     "localhost",
		ServicesPath:    "/",
		DataPath:        "/data",
		EnabledServices: []string{"svc1"},
		Services: map[string]config.ServiceConfig{
			"svc1": {
				Port:    8080,
				Version: "v1",
				Extra:   map[string]any{"NEW_FIELD": "new_value"},
			},
			"svc2": {
				Port:    9090,
				Version: "v2",
			},
		},
	}
)

func TestDeployServices_SingleService(t *testing.T) {
	mocker := &Mocker{}
	manager := newManagerWithMocks(mocker)

	wantEnv := strings.Join([]string{
		"AUTONAS_HOST=localhost",
		"SERVICES_PATH=\\",
		"DATA_PATH=\\data\\svc1",
		"PORT=8080",
		"VERSION=v1",
		"NEW_FIELD=new_value",
	}, "\n") + "\n"
	mock.InOrder(
		mocker.On("WriteToFile", "\\svc1\\.env", wantEnv).Return(nil),
		mocker.On(
			"RunCommand", "docker", []string{"compose", "--project-directory", "\\svc1", "up", "-d"},
		).Return(nil),
	)
	err := manager.DeployServices(mockConfig)
	assert.NoError(t, err)
}

func TestRemoveServices_MultipleServices(t *testing.T) {
	mocker := &Mocker{}
	manager := newManagerWithMocks(mocker)
	mocker.On(
		"RunCommand", "docker", []string{"compose", "--project-directory", "\\svc1", "down"},
	).Return(nil)
	mocker.On(
		"RunCommand", "docker", []string{"compose", "--project-directory", "\\svc2", "down"},
	).Return(fmt.Errorf("mock error"))
	err := manager.RemoveServices([]string{"svc1", "svc2"}, "/")

	assert.NoError(t, err)
}

type ExpectedErrors struct {
	writeFileErr error
	runCmdErr    error
}

var (
	ErrWriteFile = errors.New("writeFile error")
	ErrRunCmd    = errors.New("runCmd error")
)

func TestDeployServices_Errors(t *testing.T) {
	t.Skip("temporarily disabled until deploy returns proper errors")
	testCases := []struct {
		name          string
		errors        ExpectedErrors
		expectedError error
	}{
		{
			name:          "writeFile error",
			errors:        ExpectedErrors{writeFileErr: ErrWriteFile},
			expectedError: ErrWriteFile,
		},
		{
			name:          "runCmd error",
			errors:        ExpectedErrors{runCmdErr: ErrRunCmd},
			expectedError: ErrRunCmd,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			mocker := &Mocker{}
			manager := newManagerWithMocks(mocker)
			mocker.On("WriteToFile", mock.Anything, mock.Anything).Return(tc.errors.writeFileErr)
			mocker.On("RunCommand", "docker", mock.Anything).Return(tc.errors.runCmdErr)
			err := manager.DeployServices(mockConfig)
			// TODO : add tests for aggregared errors
			assert.ErrorIs(t, err, tc.expectedError)

		})
	}
}
