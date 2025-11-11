// Package config provides functionality to load and manage application configuration.
package config

import (
	"fmt"
	"strconv"
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
func (cfg Config) PerService(service string) []Variable {
	serviceConfig := []Variable{}

	for key, value := range cfg.Extra {
		serviceConfig = append(serviceConfig,
			Variable{Key: key, Value: fmt.Sprint(value)})
	}
	if svcVars, ok := cfg.Services[service]; ok {
		if svcVars.Port != 0 {
			serviceConfig = append(serviceConfig,
				Variable{Key: "PORT", Value: strconv.Itoa(svcVars.Port)})
		}
		if svcVars.Version != "" {
			serviceConfig = append(serviceConfig,
				Variable{Key: "VERSION", Value: svcVars.Version})
		}

		for key, value := range svcVars.Extra {
			serviceConfig = append(serviceConfig,
				Variable{Key: key, Value: fmt.Sprint(value)})
		}
	}
	return serviceConfig
}
