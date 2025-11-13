// Package config provides functionality to load and manage application configuration.
package config

import (
	"fmt"
	"strconv"

	"github.com/elliotchance/orderedmap/v3"
)

// ServiceConfig represents configuration for an individual service.
type ServiceConfig struct {
	Port    int            `mapstructure:"PORT"`
	Version string         `mapstructure:"VERSION"`
	Extra   map[string]any `mapstructure:",remain"`
}

// Config represents the overall configuration structure.
type Config struct {
	EnabledServices []string                 `mapstructure:"enabled_services"`
	Services        map[string]ServiceConfig `mapstructure:"services"`
	Extra           map[string]any           `mapstructure:",remain"`
}

// Variable represent an environement variable
type Variable struct {
	Key   string
	Value string
}

// PerService generates a slice of configuration variables for a specific service
func (cfg Config) PerService(service string) *orderedmap.OrderedMap[string, string] {
	serviceConfig := orderedmap.NewOrderedMap[string, string]()

	for key, value := range cfg.Extra {
		serviceConfig.Set(key, fmt.Sprint(value))
	}
	if svcVars, ok := cfg.Services[service]; ok {
		if svcVars.Port != 0 {
			serviceConfig.Set("PORT", strconv.Itoa(svcVars.Port))
		}
		if svcVars.Version != "" {
			serviceConfig.Set("VERSION", svcVars.Version)
		}

		for key, value := range svcVars.Extra {
			serviceConfig.Set(key, fmt.Sprint(value))
		}
	}
	return serviceConfig
}
