package config

import (
	"fmt"
	"os"
	"strings"
)

var (
	IsDebug     = false
	_IsDebugSet = false
)

func GetApiPrefix(defaultprefix string) string {
	prefix := os.Getenv("API_PREFIX")
	if prefix == "" {
		return defaultprefix
	}

	if prefix[0] != '/' {
		panic("the API_PREFIX must start with a '/'")
	}
	if prefix[len(prefix)-1] == '/' {
		panic("the API_PREFIX must **not** end with a '/'")
	}

	// If the given prefix is "/", the API_PREFIX will be a no-op
	if prefix == "/" {
		return ""
	}
	return prefix
}

func MustGetEnvValue(key string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		panic(fmt.Sprintf("the environment variable %s must be set", key))
	}
	return v
}

func GetEnvValueOrDefault(key, defvalue string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return defvalue
	}
	return v
}

func GetIsDebug() bool {
	if !_IsDebugSet {
		d := GetEnvValueOrDefault("DEBUG", "false")
		IsDebug = d == "true"
		_IsDebugSet = true
	}
	return IsDebug
}
