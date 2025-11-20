// Package config provides functionality to load and manage application configuration.
package config

import (
	"fmt"
	"strings"

	"github.com/elliotchance/orderedmap/v3"
)

// ServiceConfig represents configuration for an individual service.
type ServiceConfig struct {
	Disabled bool           `mapstructure:"disabled"`
	Extra    map[string]any `mapstructure:",remain"`
}

// Config represents the overall configuration structure.
type Config struct {
	Repo       string                   `mapstructure:"repo"`
	Branch     string                   `mapstructure:"branch"`
	CronPeriod string                   `mapstructure:"cron"`
	Services   map[string]ServiceConfig `mapstructure:"services"`
	Extra      map[string]any           `mapstructure:",remain"`
}

// PerService generates a slice of configuration variables for a specific service
func (cfg Config) PerService(service string) *orderedmap.OrderedMap[string, string] {
	serviceConfig := orderedmap.NewOrderedMap[string, string]()

	for key, value := range cfg.Extra {
		serviceConfig.Set(strings.ToUpper(key), fmt.Sprint(value))
	}
	if svcVars, ok := cfg.Services[service]; ok {
		for key, value := range svcVars.Extra {
			serviceConfig.Set(strings.ToUpper(key), fmt.Sprint(value))
		}
	}
	return serviceConfig
}

// GetEnabledServices returns the list of enabled services on the configuration
func (cfg Config) GetEnabledServices() []string {
	var enabled []string
	for serviceName, value := range cfg.Services {
		if !value.Disabled {
			enabled = append(enabled, serviceName)
		}
	}
	return enabled
}
