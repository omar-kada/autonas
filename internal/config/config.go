package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type ServiceConfig struct {
	Port    int                    `mapstructure:"PORT"`
	Version string                 `mapstructure:"VERSION"`
	Extra   map[string]interface{} `mapstructure:",remain"`
}

type Config struct {
	AUTONAS_HOST    string                   `mapstructure:"AUTONAS_HOST"`
	PULL            int                      `mapstructure:"PULL"`
	STOP            int                      `mapstructure:"STOP"`
	SERVICES_PATH   string                   `mapstructure:"SERVICES_PATH"`
	DATA_PATH       string                   `mapstructure:"DATA_PATH"`
	EnabledServices []string                 `mapstructure:"enabled_services"`
	Services        map[string]ServiceConfig `mapstructure:"services"`
	Extra           map[string]interface{}   `mapstructure:",remain"`
}

// loadConfig loads and merges all config files into a typed Config struct using only Viper
func LoadConfig(files []string) (*Config, error) {
	v := viper.New()
	for _, file := range files {
		v.SetConfigFile(file)
		if err := v.MergeInConfig(); err != nil {
			return nil, fmt.Errorf("error reading config file %s: %w", file, err)
		}
		fmt.Printf("Loaded config file: %s\n", file)
	}
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config to struct: %w", err)
	}
	return &cfg, nil
}
