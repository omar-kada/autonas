// Package defaults contain generic helper functions handling priority and defaults for cli variables
package defaults

import (
	"fmt"
	"os"
	"strings"
)

// VarKey defines a key for variables used in cli input
type VarKey string

// VariableInfoMap defines metadata related to cli variables
type VariableInfoMap map[VarKey]struct {
	EnvKey       string
	DefaultValue any
}

// GetDefaultStringFn generates a function that will append "(default : <value>)" to prefix
// when defaultValue exists in varInfoMap
func GetDefaultStringFn(varInfoMap VariableInfoMap) func(prefix string, key VarKey) string {
	return func(prefix string, key VarKey) string {
		val, ok := varInfoMap[key]
		if ok && val.DefaultValue != nil {
			return fmt.Sprintf("%s (default : %v)", prefix, val.DefaultValue)
		}
		return prefix
	}
}

// EnvOrDefaultFn handles priority for cli variables while ignoring empty values,
// priority is : cli argument > env variable > default value
func EnvOrDefaultFn(varInfoMap VariableInfoMap) func(value string, key VarKey) string {
	return func(value string, key VarKey) string {
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
}

// EnvOrDefaultSliceFn handles priority for cli variables of type slice while ignoring empty values,
// priority is : cli argument > env variable > default value
func EnvOrDefaultSliceFn(varInfoMap VariableInfoMap) func(value []string, key VarKey) []string {
	return func(value []string, key VarKey) []string {
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
}
