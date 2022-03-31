package fixtures

import (
	"encoding/base64"
	"os"
	"strings"

	"github.com/whitekid/goxp/fx"
	"github.com/whitekid/goxp/log"
)

// Env environment fixture
func Env(key, value string) Teardown {
	old, exists := os.LookupEnv(key)
	os.Setenv(key, value)

	cleared := false

	return func() {
		if cleared {
			return
		}

		if exists {
			os.Setenv(key, old)
		} else {
			os.Unsetenv(key)
		}

		cleared = true
	}
}

// Envs multiple environments fixture
func Envs(envs map[string]string) Teardown {
	log.Debugf("envs = %+v", envs)
	teardownsMap := fx.MapMap(envs, func(x string) Teardown {
		log.Debugf("%s = %s", x, envs[x])
		return Env(x, envs[x])
	})

	return Chain(fx.Values(teardownsMap)...)
}

// JSONEnv json environment fixture
func JSONEnv(key, value string) Teardown {
	encoded := base64.StdEncoding.EncodeToString([]byte(strings.TrimSpace(value)))
	return Env(key, encoded)
}
