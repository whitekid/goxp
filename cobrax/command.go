package cobrax

import "github.com/spf13/cobra"

func Add(parent, cmd *cobra.Command, fn func(cmd *cobra.Command)) *cobra.Command {
	parent.AddCommand(cmd)

	if fn != nil {
		fn(cmd)
	}

	return cmd
}
