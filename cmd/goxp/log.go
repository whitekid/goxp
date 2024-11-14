package main

import (
	"github.com/spf13/cobra"

	"github.com/whitekid/goxp/log"
)

func init() {
	cmd := &cobra.Command{
		Use: "log",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Debug("DEBUG")
			log.Info("INFO")
			return nil
		},
	}

	rootCmd.AddCommand(cmd)

}
