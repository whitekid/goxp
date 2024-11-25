package main

import (
	"github.com/spf13/cobra"

	"github.com/whitekid/goxp/cobrax"
	"github.com/whitekid/goxp/log"
)

func init() {
	cobrax.Add(rootCmd, &cobra.Command{
		Use: "log",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Debug("DEBUG")
			log.Info("INFO")
			return nil
		},
	}, nil)
}
