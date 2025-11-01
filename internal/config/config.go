// Package config provides functionality to load and manage application configuration.
package config

import (
	"fmt"
	"os"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
)

// ServiceConfig represents configuration for an individual service.
type ServiceConfig struct {
	Port    int            `mapstructure:"PORT"`
	Version string         `mapstructure:"VERSION"`
	Extra   map[string]any `mapstructure:",remain"`
}

// Config represents the overall configuration structure.
type Config struct {
	AutonasHost     string                   `mapstructure:"AUTONAS_HOST"`
	ServicesPath    string                   `mapstructure:"SERVICES_PATH"`
	DataPath        string                   `mapstructure:"DATA_PATH"`
	EnabledServices []string                 `mapstructure:"enabled_services"`
	Services        map[string]ServiceConfig `mapstructure:"services"`
	Extra           map[string]any           `mapstructure:",remain"`
}

// FromFiles reads YAML files, merges them (later files override earlier ones),
func FromFiles(files []string) (Config, error) {
	merged := make(map[string]any)

	for _, file := range files {
		bs, err := os.ReadFile(file)
		if err != nil {
			return Config{}, fmt.Errorf("error reading config file %s: %w", file, err)
		}

		var m map[string]any
		if err := yaml.Unmarshal(bs, &m); err != nil {
			return Config{}, fmt.Errorf("error unmarshaling yaml %s: %w", file, err)
		}

		merged = mergeMaps(merged, m)
		fmt.Printf("Loaded config file: %s\n", file)
	}

	cfg, err := decodeConfig(merged)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func decodeConfig(configMap map[string]any) (Config, error) {
	var cfg Config
	decCfg := &mapstructure.DecoderConfig{
		TagName:          "mapstructure",
		Result:           &cfg,
		WeaklyTypedInput: true,
	}
	decoder, err := mapstructure.NewDecoder(decCfg)
	if err != nil {
		return Config{}, fmt.Errorf("failed to create decoder: %w", err)
	}
	if err := decoder.Decode(configMap); err != nil {
		return Config{}, fmt.Errorf("error decoding merged config: %w", err)
	}
	return cfg, nil
}

// mergeMaps merges src into dst recursively.
// Values from src override dst; original key case is preserved.
func mergeMaps(dst, src map[string]any) map[string]any {
	if dst == nil {
		dst = make(map[string]any)
	}
	for k, v := range src {
		if existing, ok := dst[k]; ok {
			oldMapValue, evok := existing.(map[string]any)
			newMapValue, svok := v.(map[string]any)
			if evok && svok {
				dst[k] = mergeMaps(oldMapValue, newMapValue)
				continue
			}
		}
		dst[k] = v
	}
	return dst
}

// PerService generates a configuration map for a specific service
func (cfg *Config) PerService(service string) map[string]any {
	serviceConfig := make(map[string]any)
	serviceConfig["AUTONAS_HOST"] = cfg.AutonasHost
	serviceConfig["SERVICES_PATH"] = cfg.ServicesPath
	serviceConfig["DATA_PATH"] = fmt.Sprintf("%s/%s", cfg.DataPath, service)
	for k, v := range cfg.Extra {
		serviceConfig[k] = v
	}
	if svcVars, ok := cfg.Services[service]; ok {
		if svcVars.Port != 0 {
			serviceConfig["PORT"] = svcVars.Port
		}
		if svcVars.Version != "" {
			serviceConfig["VERSION"] = svcVars.Version
		}
		for k, v := range svcVars.Extra {
			serviceConfig[k] = v
		}
	}
	return serviceConfig
}
