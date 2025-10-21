package config

import (
	"fmt"
	"os"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
)

type ServiceConfig struct {
	Port    int                    `mapstructure:"PORT"`
	Version string                 `mapstructure:"VERSION"`
	Extra   map[string]interface{} `mapstructure:",remain"`
}

type Config struct {
	AUTONAS_HOST    string                   `mapstructure:"AUTONAS_HOST"`
	SERVICES_PATH   string                   `mapstructure:"SERVICES_PATH"`
	DATA_PATH       string                   `mapstructure:"DATA_PATH"`
	EnabledServices []string                 `mapstructure:"enabled_services"`
	Services        map[string]ServiceConfig `mapstructure:"services"`
	Extra           map[string]interface{}   `mapstructure:",remain"`
}

// LoadConfig reads YAML files, merges them (later files override earlier ones),
// preserves key case for unknown keys, and decodes into Config using mapstructure.
func LoadConfig(files []string) (Config, error) {
	merged := make(map[string]interface{})

	for _, file := range files {
		bs, err := os.ReadFile(file)
		if err != nil {
			return Config{}, fmt.Errorf("error reading config file %s: %w", file, err)
		}

		var m map[string]interface{}
		if err := yaml.Unmarshal(bs, &m); err != nil {
			return Config{}, fmt.Errorf("error unmarshaling yaml %s: %w", file, err)
		}

		merged = mergeMaps(merged, m)
		fmt.Printf("Loaded config file: %s\n", file)
	}

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
	if err := decoder.Decode(merged); err != nil {
		return Config{}, fmt.Errorf("error decoding merged config: %w", err)
	}
	return cfg, nil
}

// mergeMaps merges src into dst recursively.
// Values from src override dst; original key case is preserved.
func mergeMaps(dst, src map[string]interface{}) map[string]interface{} {
	if dst == nil {
		dst = make(map[string]interface{})
	}
	for k, v := range src {
		if existing, ok := dst[k]; ok {
			ev, evok := existing.(map[string]interface{})
			sv, svok := v.(map[string]interface{})
			if evok && svok {
				dst[k] = mergeMaps(ev, sv)
				continue
			}
		}
		dst[k] = v
	}
	return dst
}

func ConfigPerService(cfg Config, service string) map[string]interface{} {
	serviceConfig := make(map[string]interface{})
	serviceConfig["AUTONAS_HOST"] = cfg.AUTONAS_HOST
	serviceConfig["SERVICES_PATH"] = cfg.SERVICES_PATH
	serviceConfig["DATA_PATH"] = fmt.Sprintf("%s/%s", cfg.DATA_PATH, service)
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
