package models

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/elliotchance/orderedmap/v3"
)

// DefaultBranch is the default branch name used when no branch is specified in the configuration.
const DefaultBranch = "main"

// Settings represents configuration of autonas.
type Settings struct {
	Repo       string `mapstructure:"repo"`
	Branch     string `mapstructure:"branch"`
	CronPeriod string `mapstructure:"cron"`
}

// Environment represents global environment variables.
type Environment map[string]string

// ServiceConfig represents configuration for an individual service.
type ServiceConfig map[string]string

// Config represents the overall configuration structure.
type Config struct {
	Settings    Settings                     `mapstructure:"settings"`
	Environment map[string]string            `mapstructure:"environment"`
	Services    map[string]map[string]string `mapstructure:"services"`
}

// PerService generates a slice of configuration variables for a specific service
func (cfg Config) PerService(service string) *orderedmap.OrderedMap[string, string] {
	serviceConfig := orderedmap.NewOrderedMap[string, string]()

	for key, value := range cfg.Environment {
		serviceConfig.Set(strings.ToUpper(key), fmt.Sprint(value))
	}
	if svcVars, ok := cfg.Services[service]; ok {
		for key, value := range svcVars {
			serviceConfig.Set(strings.ToUpper(key), fmt.Sprint(value))
		}
	}
	return serviceConfig
}

// GetEnabledServices returns the list of enabled services on the configuration
func (cfg Config) GetEnabledServices() []string {
	return slices.Collect(maps.Keys(cfg.Services))
}

// GetBranch returns the branch name from the configuration. If no branch is specified,
// it defaults to "main".
func (cfg Config) GetBranch() string {
	if cfg.Settings.Branch != "" {
		return cfg.Settings.Branch
	}
	return DefaultBranch
}
