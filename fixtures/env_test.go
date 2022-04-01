package fixtures

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/whitekid/goxp/fx"
)

func TestEnv(t *testing.T) {
	t.Parallel()

	os.Setenv("HELLO", "OLD")

	teardown := Env("HELLO", "WORLD")
	defer teardown()
	require.Equal(t, "WORLD", os.Getenv("HELLO"))

	teardown()
	require.Equal(t, "OLD", os.Getenv("HELLO"))
}

func TestEnvs(t *testing.T) {
	vars := map[string]string{
		"HELLO": "WORLD",
		"SEOUL": "KOREA",
	}
	teardown := Envs(vars)
	defer teardown()

	fx.ForEachMap(vars, func(k, v string) {
		require.Equal(t, v, os.Getenv(k))
	})
}
