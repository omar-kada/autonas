// Package defaults contain generic helper functions handling priority and defaults for cli variables
package defaults

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// VarKey defines a key for variables used in cli input
type VarKey string

// VariableInfoMap defines metadata related to cli variables
type VariableInfoMap map[VarKey]struct {
	EnvKey       string
	DefaultValue any
}

// GetDefaultString generates a function that will append "(default : <value>)" to prefix
// when defaultValue exists in varInfoMap
func (varInfoMap VariableInfoMap) GetDefaultString(prefix string, key VarKey) string {
	val, ok := varInfoMap[key]
	if ok && val.DefaultValue != nil {
		return fmt.Sprintf("%s (default : %v)", prefix, val.DefaultValue)
	}
	return prefix
}

// EnvOrDefault handles priority for cli variables while ignoring empty values,
// priority is : cli argument > env variable > default value
func (varInfoMap VariableInfoMap) EnvOrDefault(value string, key VarKey) string {
	if value != "" {
		return value
	}
	varInfo := varInfoMap[key]
	if envValue := os.Getenv(varInfo.EnvKey); envValue != "" {
		return envValue
	}
	if varInfo.DefaultValue == nil {
		return ""
	}
	return varInfo.DefaultValue.(string)
}

// EnvOrDefaultInt handles priority for cli variables while ignoring empty values,
// priority is : cli argument > env variable > default value
func (varInfoMap VariableInfoMap) EnvOrDefaultInt(value int, key VarKey) int {
	if value != 0 {
		return value
	}
	varInfo := varInfoMap[key]
	if envValue := os.Getenv(varInfo.EnvKey); envValue != "" {
		val, err := strconv.Atoi(envValue)
		if err != nil {
			return val
		}
	}
	if varInfo.DefaultValue == nil {
		return 0
	}
	return varInfo.DefaultValue.(int)
}

// EnvOrDefaultSlice handles priority for cli variables of type slice while ignoring empty values,
// priority is : cli argument > env variable > default value
func (varInfoMap VariableInfoMap) EnvOrDefaultSlice(value []string, key VarKey) []string {
	if len(value) > 0 {
		return value
	}
	varInfo := varInfoMap[key]
	if envStr := os.Getenv(varInfo.EnvKey); envStr != "" {
		envValue := strings.Split(envStr, ",")
		if len(envValue) > 0 {
			return envValue
		}
	}
	if varInfo.DefaultValue == nil {
		return []string{}
	}
	return varInfo.DefaultValue.([]string)
}
