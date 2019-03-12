package fixtures

import (
	"encoding/base64"
	"os"
	"strings"
)

// Env setup fixture environment
func Env(key, data string) Fixture {
	os.Setenv(key, data)

	return &envFixture{key: key}
}

type envFixture struct {
	key string
}

func (e *envFixture) Teardown() {
	os.Unsetenv(e.key)
}

// JSONEnv ...
func JSONEnv(key, data string) Fixture {
	encoded := base64.StdEncoding.EncodeToString([]byte(strings.TrimSpace(data)))
	return Env(key, encoded)
}
