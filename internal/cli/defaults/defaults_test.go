package defaults

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDefaultString_WithDefault(t *testing.T) {
	varMap := VariableInfoMap{
		"k1": {EnvKey: "E1", DefaultValue: "d1"},
	}

	got := varMap.GetDefaultString("prefix", "k1")
	want := "prefix (default : d1)"
	assert.Equal(t, want, got)
}

func TestGetDefaultString_NoDefault(t *testing.T) {
	varMap := VariableInfoMap{
		"k1": {EnvKey: "E1", DefaultValue: nil},
	}

	got := varMap.GetDefaultString("prefix", "k1")
	want := "prefix"
	assert.Equal(t, want, got)
}

func TestEnvOrDefault_Priorities(t *testing.T) {
	key := VarKey("varName")
	varMap := VariableInfoMap{
		key: {EnvKey: "ENV_KEY", DefaultValue: "default"},
	}

	t.Setenv("ENV_KEY", "envValue")

	assert.Equal(t, "cliValue", varMap.EnvOrDefault("cliValue", key), "CLI value has most priority")
	assert.Equal(t, "envValue", varMap.EnvOrDefault("", key), "ENV value has 2nd most priority")
	t.Setenv("ENV_KEY", "")
	assert.Equal(t, "default", varMap.EnvOrDefault("", key), "default value has least priority")
}

func TestEnvOrDefaultInt_Priorities(t *testing.T) {
	key := VarKey("intVarName")
	cliValue := 123
	defaultValue := 456

	varMap := VariableInfoMap{
		key: {EnvKey: "ENV_KEY", DefaultValue: defaultValue},
	}

	t.Setenv("ENV_KEY", "789")
	envValue := 789

	assert.Equal(t, cliValue, varMap.EnvOrDefaultInt(cliValue, key), "CLI value has most priority")
	assert.Equal(t, envValue, varMap.EnvOrDefaultInt(0, key), "ENV value has 2nd most priority")
	t.Setenv("ENV_KEY", "")
	assert.Equal(t, defaultValue, varMap.EnvOrDefaultInt(0, key), "default value has least priority")
}

func TestEnvOrDefaultInt_InvalidEnvValue(t *testing.T) {
	key := VarKey("intVarName")
	defaultValue := 456

	varMap := VariableInfoMap{
		key: {EnvKey: "ENV_KEY", DefaultValue: defaultValue},
	}

	t.Setenv("ENV_KEY", "invalid")

	assert.Equal(t, defaultValue, varMap.EnvOrDefaultInt(0, key), "Invalid ENV value should fall back to default")
}
