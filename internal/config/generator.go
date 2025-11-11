package config

import (
	"fmt"
	"os"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
)

// Generator is responsible for creating Config from sources
type Generator interface {
	FromFiles(files []string) (Config, error)
}

// generator is responsible for creating Config from sources
type generator struct {
}

// NewGenerator creates a new Generator and returns it
func NewGenerator() Generator {
	return generator{}
}

// FromFiles reads YAML files, merges them (later files override earlier ones),
func (g generator) FromFiles(files []string) (Config, error) {
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
	}

	return decodeConfig(merged)
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
