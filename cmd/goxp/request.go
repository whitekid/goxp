package main

import (
	"github.com/spf13/cobra"

	"github.com/whitekid/goxp/cmd/goxp/request"
	"github.com/whitekid/goxp/cobrax"
)

func init() {
	cobrax.Add(rootCmd, &cobra.Command{
		Use:   "request url",
		Short: "request package example",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return request.Run(cmd.Context(), args[0])
		},
	}, func(cmd *cobra.Command) {
		request.SetFlags(cmd.PersistentFlags(), cmd.Flags(), GitTag)
	})
}
