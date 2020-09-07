package flags

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestDefaults(t *testing.T) {
	configs := map[string][]Flag{
		"hello": {
			{"hello", "h", "world", "hello world"},
		},
	}

	InitDefaults(configs)

	for _, use := range configs {
		for _, config := range use {
			require.Equal(t, config.DefaultValue, viper.Get(config.Key))
		}
	}

	cmd := &cobra.Command{
		Use: "hello",
	}

	InitFlagSet(configs, cmd.Use, cmd.Flags())
}
