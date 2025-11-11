package defaults

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDefaultStringFn_WithDefault(t *testing.T) {
	varMap := VariableInfoMap{
		"k1": {EnvKey: "E1", DefaultValue: "d1"},
	}

	got := GetDefaultStringFn(varMap)("prefix", "k1")
	want := "prefix (default : d1)"
	assert.Equal(t, want, got)
}

func TestGetDefaultStringFn_NoDefault(t *testing.T) {
	varMap := VariableInfoMap{
		"k1": {EnvKey: "E1", DefaultValue: nil},
	}

	got := GetDefaultStringFn(varMap)("prefix", "k1")
	want := "prefix"
	assert.Equal(t, want, got)
}

func TestEnvOrDefaultFn_Priorities(t *testing.T) {
	key := VarKey("varName")
	varMap := VariableInfoMap{
		key: {EnvKey: "ENV_KEY", DefaultValue: "default"},
	}
	envOrDefault := EnvOrDefaultFn(varMap)

	orig := os.Getenv("ENV_KEY")
	defer os.Setenv("ENV_KEY", orig)
	os.Setenv("ENV_KEY", "envValue")

	assert.Equal(t, "cliValue", envOrDefault("cliValue", key), "CLI value has most priority")
	assert.Equal(t, "envValue", envOrDefault("", key), "ENV value has 2nd most priority")
	os.Unsetenv("ENV_KEY")
	assert.Equal(t, "default", envOrDefault("", key), "default value has least priority")
}

func TestEnvOrDefaultSliceFn_Priorities(t *testing.T) {
	key := VarKey("sliceVarName")
	cliValue := []string{"x", "y"}
	defaultValue := []string{"a"}

	varMap := VariableInfoMap{
		key: {EnvKey: "ENV_KEY", DefaultValue: defaultValue},
	}
	envOrDefaultSlice := EnvOrDefaultSliceFn(varMap)

	orig := os.Getenv("ENV_KEY")
	defer os.Setenv("ENV_KEY", orig)
	os.Setenv("ENV_KEY", "e1,e2")
	envValue := []string{"e1", "e2"}

	assert.Equal(t, cliValue, envOrDefaultSlice(cliValue, key), "CLI value has most priority")
	assert.Equal(t, envValue, envOrDefaultSlice(nil, key), "ENV value has 2nd most priority")
	os.Unsetenv("ENV_KEY")
	assert.Equal(t, defaultValue, envOrDefaultSlice(nil, key), "default value has least priority")
}
