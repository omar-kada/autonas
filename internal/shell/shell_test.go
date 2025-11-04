package shell

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunCommand_Success(t *testing.T) {
	err := RunCommand("go", "help")
	assert.NoError(t, err)
}

func TestRunCommand_NoArgs(t *testing.T) {
	err := RunCommand("go")
	assert.ErrorContains(t, err, "exit status 2")
}

func TestRunCommand_NotFound(t *testing.T) {
	err := RunCommand("dummyCmd")
	assert.ErrorContains(t, err, "executable not found")
}
