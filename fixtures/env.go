package fixtures

import (
	"encoding/base64"
	"maps"
	"os"
	"slices"
	"strings"
	"sync"

	"github.com/whitekid/goxp/mapx"
	"github.com/whitekid/goxp/slicex"
)

// Env environment fixture
func Env(key, value string) Teardown {
	old, exists := os.LookupEnv(key)
	os.Setenv(key, value)

	var once sync.Once

	return func() {
		once.Do(func() {
			if exists {
				os.Setenv(key, old)
			} else {
				os.Unsetenv(key)
			}
		})
	}
}

// Envs multiple environments fixture
func Envs(envs map[string]string) Teardown {
	teardowns := mapx.MapValue(envs, func(k, v string) Teardown {
		return Env(k, v)
	})

	return Chain(slices.Collect(maps.Values(teardowns))...)
}

// JSONEnv json environment fixture
func JSONEnv(key, value string) Teardown {
	encoded := base64.RawStdEncoding.EncodeToString([]byte(strings.TrimSpace(value)))
	return Env(key, encoded)
}

// UnsetEnv ensure environment was unset and return recover environment recover teardown
func UnsetEnv(key string) Teardown {
	value, exists := os.LookupEnv(key)
	if exists {
		os.Unsetenv(key)
	}

	var once sync.Once

	return func() {
		once.Do(func() {
			if exists {
				os.Setenv(key, value)
			} else {
				os.Unsetenv(key)
			}
		})
	}
}

// UnsetEnvs multiple environments fixture
func UnsetEnvs(envs []string) Teardown {
	return Chain(slicex.Map(envs, func(k string) Teardown {
		return UnsetEnv(k)
	})...)
}
