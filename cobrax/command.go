// Package cobrax provides utility functions for working with cobra commands.
package cobrax

import "github.com/spf13/cobra"

// Add creates or adds a cobra command with optional configuration.
// When parent is nil, cmd is treated as a root command.
// When parent is non-nil, cmd is added as a subcommand.
// The optional fn callback allows for command configuration.
//
// Example usage:
//
//	rootCmd := cobrax.Add(nil, &cobra.Command{Use: "app"}, nil)
//	subCmd := cobrax.Add(rootCmd, &cobra.Command{Use: "sub"}, func(cmd *cobra.Command) {
//	    cmd.Flags().StringP("name", "n", "", "name flag")
//	})
func Add(parent, cmd *cobra.Command, fn func(cmd *cobra.Command)) *cobra.Command {
	if parent != nil {
		parent.AddCommand(cmd)
	}

	if fn != nil {
		fn(cmd)
	}

	return cmd
}
