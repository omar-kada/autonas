package storage

import (
	"fmt"
	"log/slog"
	"os"

	"omar-kada/autonas/models"

	"github.com/go-viper/mapstructure/v2"
	"go.yaml.in/yaml/v3"
)

// ConfigStore stores and retreives the configuration
type ConfigStore interface {
	Update(cfg models.Config) error
	Get() (models.Config, error)
	SetOnChange(fn func(oldConfig, newConfig models.Config))
}

type configStore struct {
	OnConfigUpdate func(oldConfig, newConfig models.Config)
	configFilePath string
}

// NewConfigStore creates a new config file storage
func NewConfigStore(filePath string) ConfigStore {
	return &configStore{
		configFilePath: filePath,
	}
}

func (s *configStore) Update(cfg models.Config) (err error) {

	slog.Debug("updating configuration file")
	if s.OnConfigUpdate != nil {
		oldCfg, err := s.Get()
		if err != nil {
			return err
		}
		defer func() {
			if err != nil { // check no error occurred when updating the config
				return
			}
			s.OnConfigUpdate(oldCfg, cfg)
		}()
	}

	var m map[string]any
	encCfg := &mapstructure.DecoderConfig{
		TagName:          "mapstructure",
		Result:           &m,
		WeaklyTypedInput: true,
	}
	encoder, err := mapstructure.NewDecoder(encCfg)
	if err != nil {
		return fmt.Errorf("failed to create encoder: %w", err)
	}
	if err := encoder.Decode(cfg); err != nil {
		return fmt.Errorf("error encoding config: %w", err)
	}

	bs, err := yaml.Marshal(m)
	if err != nil {
		return fmt.Errorf("error marshaling config: %w", err)
	}

	if err := os.WriteFile(s.configFilePath, bs, 0o644); err != nil {
		return fmt.Errorf("error writing config file %s: %w", s.configFilePath, err)
	}

	return nil
}

func (s *configStore) SetOnChange(fn func(oldConfig, newConfig models.Config)) {
	slog.Debug("setting OnConfigUpdate")
	s.OnConfigUpdate = fn
}

// Get reads the configuration from the config file
func (s *configStore) Get() (models.Config, error) {
	bs, err := os.ReadFile(s.configFilePath)
	if err != nil {
		return models.Config{}, fmt.Errorf("error reading config file %s: %w", s.configFilePath, err)
	}

	var m map[string]any
	if err := yaml.Unmarshal(bs, &m); err != nil {
		return models.Config{}, fmt.Errorf("error unmarshaling yaml %s: %w", s.configFilePath, err)
	}

	return decodeConfig(m)
}

func decodeConfig(configMap map[string]any) (models.Config, error) {
	var cfg models.Config
	decCfg := &mapstructure.DecoderConfig{
		TagName:          "mapstructure",
		Result:           &cfg,
		WeaklyTypedInput: true,
	}
	decoder, err := mapstructure.NewDecoder(decCfg)
	if err != nil {
		return models.Config{}, fmt.Errorf("failed to create decoder: %w", err)
	}
	if err := decoder.Decode(configMap); err != nil {
		return models.Config{}, fmt.Errorf("error decoding merged config: %w", err)
	}
	return cfg, nil
}
