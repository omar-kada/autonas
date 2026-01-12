package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"omar-kada/autonas/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type Mocker struct {
	mock.Mock
}

func (m *Mocker) Exec(cmd string, cmdArgs ...string) error {
	args := m.Called(cmd, cmdArgs)
	return args.Error(0)
}

func TestRunCommand_CmdParams(t *testing.T) {
	baseDir := t.TempDir()
	mocker := &Mocker{}
	cmd := NewRunCommand(mocker)

	servicesDir := filepath.Join(baseDir, "services")
	dataDir := filepath.Join(baseDir, "data")
	workingDir := filepath.Join(baseDir, "work")
	configFile := filepath.Join(workingDir, "config.yaml")
	os.MkdirAll(servicesDir, 0o750)
	os.MkdirAll(dataDir, 0o750)
	os.MkdirAll(workingDir, 0o750)

	mocker.On(
		"Exec", "docker",
		[]string{"compose", "--project-directory", filepath.Join(servicesDir, "homepage"), "up", "-d"},
	).Return(nil)

	err := os.WriteFile(configFile,
		[]byte(strings.Join([]string{
			"SERVICES_PATH: " + servicesDir,
			"DATA_PATH: " + dataDir,
			"repo: \"https://github.com/omar-kada/autonas-config\"",
			"cron: \"* * * * *\"",
			"services:",
			"  homepage:",
			"    port : 12345",
		}, "\n")), 0o750)
	assert.NoError(t, err, "error while creating config file")

	t.Setenv("AUTONAS_CONFIG_FILE", "")
	t.Setenv("AUTONAS_WORKING_DIR", "")

	go func() {
		cmd.SetArgs([]string{
			"-f", configFile,
			"-d", workingDir,
			"-s", servicesDir,
			"-w", "true",
			"-p", "5008",
		})
		cmd.Execute()
	}()

	targetFile := filepath.Join(servicesDir, "homepage", ".env")

	ok := testutil.WaitForFile(targetFile, 1*time.Minute)
	assert.True(t, ok, "homepage files should be created")
}

func TestRunCommand_EnvParams(t *testing.T) {
	baseDir := t.TempDir()
	mocker := &Mocker{}
	cmd := NewRunCommand(mocker)

	customServicesDir := filepath.Join(baseDir, "custom_services")
	customDataDir := filepath.Join(baseDir, "custom_data")
	customWorkingDir := filepath.Join(baseDir, "custom_work")
	customConfigFile := filepath.Join(customWorkingDir, "custom_config.yaml")
	os.MkdirAll(customServicesDir, 0o750)
	os.MkdirAll(customDataDir, 0o750)
	os.MkdirAll(customWorkingDir, 0o750)

	mocker.On(
		"Exec", "docker",
		[]string{"compose", "--project-directory", filepath.Join(customServicesDir, "homepage"), "up", "-d"},
	).Return(nil)

	err := os.WriteFile(customConfigFile,
		[]byte(strings.Join([]string{
			"SERVICES_PATH: " + customServicesDir,
			"DATA_PATH: " + customDataDir,
			"repo: \"https://github.com/omar-kada/autonas-config\"",
			"cron: \"* * * * *\"",
			"services:",
			"  homepage:",
			"    port : 54321",
		}, "\n")), 0o750)
	assert.NoError(t, err, "error while creating custom config file")

	t.Setenv("AUTONAS_CONFIG_FILE", customConfigFile)
	t.Setenv("AUTONAS_WORKING_DIR", customWorkingDir)
	t.Setenv("AUTONAS_SERVICES_DIR", customServicesDir)
	t.Setenv("AUTONAS_ADD_WRITE_PERM", "true")
	t.Setenv("AUTONAS_PORT", "0")

	go func() {
		cmd.Execute()
	}()

	targetFile := filepath.Join(customServicesDir, "homepage", ".env")

	ok := testutil.WaitForFile(targetFile, 1*time.Minute)
	assert.True(t, ok, "custom homepage files should be created")
}

func TestRunCommand_WithInvalidConfig(t *testing.T) {
	mocker := &Mocker{}
	cmd := NewRunCommand(mocker)

	t.Setenv("AUTONAS_WORKING_DIR", "/invalid")

	// Create a channel to capture the command's exit status
	done := make(chan error, 1)

	go func() {
		err := cmd.Execute()
		done <- err
	}()

	select {

	case cmdErr := <-done:
		assert.ErrorContains(t, cmdErr, "couldn't init sqlite db")
	case <-time.After(1 * time.Second):
		assert.Fail(t, "timeout while waiting for command error")
	}
}
