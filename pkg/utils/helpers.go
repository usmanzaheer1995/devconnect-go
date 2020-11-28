package utils

import (
	"errors"
	"fmt"
	"os"
)

// EnvString returns the environment variable or else, the fallback
func EnvString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}

func SetEnvString(key, value string) error {
	if value == "" {
		return errors.New(fmt.Sprintf("value for key %s not found", key))
	}
	return os.Setenv(key, value)
}
