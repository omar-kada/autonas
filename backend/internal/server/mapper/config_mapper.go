package mapper

import (
	"omar-kada/autonas/api"
	"omar-kada/autonas/models"
)

// ConfigMapper maps models.Config to api.Config
type ConfigMapper struct{}

// Map converts a models.Config to an api.Config
func (ConfigMapper) Map(config models.Config) api.Config {
	convertedMap := make(map[string]map[string]string)

	for key, innerMap := range config.Services {
		convertedMap[key] = innerMap
	}

	return api.Config{
		GlobalVariables: config.Environment,
		Services:        convertedMap,
	}
}
