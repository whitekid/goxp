package fixtures

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/whitekid/goxp"
	"github.com/whitekid/goxp/mapx"
	"github.com/whitekid/goxp/slicex"
)

func TestEnv(t *testing.T) {
	os.Setenv("HELLO", "OLD")
	defer os.Unsetenv("HELLO")

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

	mapx.Each(vars, func(k, v string) {
		require.Equal(t, v, os.Getenv(k))
	})
}

func TestUnsetEnv(t *testing.T) {
	defer Env("HELLO", "WORLD")()

	teardown := UnsetEnv("HELLO")
	require.False(t, goxp.EnvExists("HELLO"))

	teardown()
	value, exists := os.LookupEnv("HELLO")
	require.Equal(t, "WORLD", value)
	require.True(t, exists)
}

func TestUnsetEnvs(t *testing.T) {
	envs := []string{"HELLO", "SEOUL"}
	teardown := Chain(slicex.Map(envs, func(k string) Teardown { return Env(k, k+"_value") })...)
	defer teardown()

	slicex.Each(envs, func(_ int, k string) {
		require.Equal(t, k+"_value", os.Getenv(k))
	})

	teardown() // clear env
	slicex.Each(envs, func(_ int, k string) {
		require.False(t, goxp.EnvExists(k))
	})
}
