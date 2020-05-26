package fixtures

import (
	"encoding/base64"
	"os"
	"strings"
)

// Env environment fixture
func Env(key, value string) Teardown {
	os.Setenv(key, value)

	return func() {
		os.Unsetenv(key)
	}
}

// JSONEnv json environment fixture
func JSONEnv(key, value string) Teardown {
	encoded := base64.StdEncoding.EncodeToString([]byte(strings.TrimSpace(value)))
	return Env(key, encoded)
}
