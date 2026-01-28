package models

import (
	"os"
	"strconv"
)

// Features represents the feature flags that can be enabled or disabled.
type Features struct {
	DisplayConfig bool
	EditConfig    bool
	EditSettings  bool
}

// LoadFeatures loads feature flags from environment variables.
func LoadFeatures() Features {
	return Features{
		DisplayConfig: getBool("AUTONAS_DISPLAY_CONFIG", false),
		EditConfig:    getBool("AUTONAS_EDIT_CONFIG", false),
		EditSettings:  getBool("AUTONAS_EDIT_SETTINGS", false),
	}
}

func getBool(key string, def bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return def
	}
	return b
}
