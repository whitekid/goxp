package fixtures

import (
	"encoding/base64"
	"os"
	"strings"
)

// Env environment fixture
func Env(key, value string) Teardown {
	old, exists := os.LookupEnv(key)
	os.Setenv(key, value)

	return func() {
		if exists {
			os.Setenv(key, old)
		} else {
			os.Unsetenv(key)
		}
	}
}

// Envs multiple environments fixture
func Envs(envs map[string]string) Teardown {
	teardowns := make([]Teardown, len(envs))
	for k, v := range envs {
		teardowns = append(teardowns, Env(k, v))
	}

	return Chain(teardowns...)
}

// JSONEnv json environment fixture
func JSONEnv(key, value string) Teardown {
	encoded := base64.StdEncoding.EncodeToString([]byte(strings.TrimSpace(value)))
	return Env(key, encoded)
}
