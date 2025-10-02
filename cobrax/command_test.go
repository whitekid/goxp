package cobrax

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {
	t.Run("add to parent command", func(t *testing.T) {
		parent := &cobra.Command{Use: "parent"}
		cmd := &cobra.Command{Use: "child"}

		result := Add(parent, cmd, nil)

		require.Equal(t, cmd, result)
		require.Len(t, parent.Commands(), 1)
		require.Equal(t, cmd, parent.Commands()[0])
	})

	t.Run("create root command with nil parent", func(t *testing.T) {
		cmd := &cobra.Command{Use: "root"}

		result := Add(nil, cmd, nil)

		require.Equal(t, cmd, result)
	})

	t.Run("execute function callback", func(t *testing.T) {
		parent := &cobra.Command{Use: "parent"}
		cmd := &cobra.Command{Use: "child"}

		var called bool
		result := Add(parent, cmd, func(c *cobra.Command) {
			called = true
			require.Equal(t, cmd, c)
		})

		require.Equal(t, cmd, result)
		require.True(t, called)
	})

	t.Run("execute function callback for root command", func(t *testing.T) {
		cmd := &cobra.Command{Use: "root"}

		var called bool
		result := Add(nil, cmd, func(c *cobra.Command) {
			called = true
			require.Equal(t, cmd, c)
		})

		require.Equal(t, cmd, result)
		require.True(t, called)
	})
}

func ExampleAdd() {
	// Create a root command
	rootCmd := Add(nil, &cobra.Command{
		Use:   "myapp",
		Short: "My application",
	}, nil)

	// Add a subcommand with configuration
	Add(rootCmd, &cobra.Command{
		Use:   "version",
		Short: "Print version information",
	}, func(cmd *cobra.Command) {
		cmd.Run = func(cmd *cobra.Command, args []string) {
			// Version command implementation
		}
	})
}

func ExampleAdd_withFlags() {
	rootCmd := Add(nil, &cobra.Command{Use: "app"}, nil)

	// Add subcommand with flags configured via callback
	Add(rootCmd, &cobra.Command{
		Use:   "greet",
		Short: "Greet someone",
	}, func(cmd *cobra.Command) {
		cmd.Flags().StringP("name", "n", "World", "name to greet")
		cmd.Run = func(cmd *cobra.Command, args []string) {
			// Greet command implementation
		}
	})
}
