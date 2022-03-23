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

	InitDefaults(nil, configs)

	for _, use := range configs {
		for _, config := range use {
			require.Equal(t, config.DefaultValue, viper.Get(config.Name))
		}
	}

	cmd := &cobra.Command{
		Use: "hello",
	}

	InitFlagSet(nil, configs, cmd.Use, cmd.Flags())
}

func TestMultipleViper(t *testing.T) {
	configs := map[string][]Flag{
		"hello": {
			{"hello", "h", "world", "hello world"},
		},
	}

	configs2 := map[string][]Flag{
		"hello": {
			{"hello", "h", "world", "hello world"},
		},
	}

	cmd1 := &cobra.Command{
		Use: "hello",
	}

	cmd2 := &cobra.Command{
		Use: "hello",
	}

	v1 := viper.New()
	InitDefaults(v1, configs)
	InitFlagSet(v1, configs, cmd1.Use, cmd1.Flags())

	v2 := viper.New()
	InitDefaults(v2, configs2)
	InitFlagSet(v2, configs2, cmd2.Use, cmd2.Flags())
}
