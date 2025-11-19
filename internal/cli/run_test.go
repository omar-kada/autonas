package cli

import (
	"omar-kada/autonas/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetParamsWithDefaults_AllCliValuesProvided(t *testing.T) {
	// When all CLI values are provided, they should be returned as-is
	params := RunParams{
		DeploymentParams: models.DeploymentParams{
			ConfigFile:   "custom.yaml",
			WorkingDir:   "/custom/work",
			ServicesDir:  "/custom/services",
			AddWritePerm: true,
		},
		ServerParams: models.ServerParams{
			Port: 9090,
		},
	}

	result := getParamsWithDefaults(params)

	assert.Equal(t, "custom.yaml", result.ConfigFile)
	assert.Equal(t, "/custom/work", result.WorkingDir)
	assert.Equal(t, "/custom/services", result.ServicesDir)
	assert.Equal(t, true, result.AddWritePerm)
	assert.Equal(t, 9090, result.Port)
}

func TestGetParamsWithDefaults_UseEnvVariablesWhenCliEmpty(t *testing.T) {
	t.Setenv("AUTONAS_WORKING_DIR", "/env/work")
	t.Setenv("AUTONAS_SERVICES_DIR", "/env/services")
	t.Setenv("AUTONAS_CONFIG_FILE", "env1.yaml")
	t.Setenv("AUTONAS_PORT", "8080")

	params := RunParams{}

	result := getParamsWithDefaults(params)

	assert.Equal(t, "env1.yaml", result.ConfigFile)
	assert.Equal(t, "/env/work", result.WorkingDir)
	assert.Equal(t, "/env/services", result.ServicesDir)
	assert.Equal(t, false, result.AddWritePerm) // default Value
	assert.Equal(t, 8080, result.Port)          // Value from env
}

func TestGetParamsWithDefaults_UseDefaultsWhenCliAndEnvEmpty(t *testing.T) {
	// Clear environment variables
	t.Setenv("AUTONAS_WORKING_DIR", "")
	t.Setenv("AUTONAS_SERVICES_DIR", "")
	t.Setenv("AUTONAS_CONFIG_FILE", "")
	t.Setenv("AUTONAS_PORT", "")

	params := RunParams{}

	result := getParamsWithDefaults(params)

	// Check defaults are applied
	assert.Equal(t, "/config/config.yaml", result.ConfigFile)
	assert.Equal(t, "./config", result.WorkingDir)
	assert.Equal(t, ".", result.ServicesDir)
	assert.Equal(t, false, result.AddWritePerm) // Default value
	assert.Equal(t, 5005, result.Port)          // Default value
}

func TestGetParamsWithDefaults_CliPriority(t *testing.T) {
	// CLI values should take priority over env variables and defaults
	t.Setenv("AUTONAS_CONFIG_BRANCH", "env-branch")
	t.Setenv("AUTONAS_PORT", "8080")

	params := RunParams{
		DeploymentParams: models.DeploymentParams{
			ServicesDir:  "/s",
			AddWritePerm: true,
		},
		ServerParams: models.ServerParams{
			Port: 9090,
		},
	}

	result := getParamsWithDefaults(params)

	// CLI value should win
	assert.Equal(t, "/s", result.ServicesDir)
	assert.Equal(t, true, result.AddWritePerm)
	assert.Equal(t, 9090, result.Port)
}

func TestGetParamsWithDefaults_MixedSources(t *testing.T) {
	// Test a mix of CLI values, env variables, and defaults
	t.Setenv("AUTONAS_CONFIG_BRANCH", "env-branch")
	t.Setenv("AUTONAS_WORKING_DIR", "")
	t.Setenv("AUTONAS_PORT", "8080")

	params := RunParams{
		DeploymentParams: models.DeploymentParams{
			ConfigFile:   "cli.yaml", // From CLI
			ServicesDir:  "/s",       // From CLI (overrides env)
			AddWritePerm: true,       // From CLI
			// WorkingDir not provided, should use env or default
		},
		ServerParams: models.ServerParams{
			Port: 9090, // From CLI
		},
	}

	result := getParamsWithDefaults(params)

	assert.Equal(t, "cli.yaml", result.ConfigFile)
	assert.Equal(t, "/s", result.ServicesDir)
	assert.Equal(t, "./config", result.WorkingDir) // Should use default
	assert.Equal(t, true, result.AddWritePerm)
	assert.Equal(t, 9090, result.Port)
}
