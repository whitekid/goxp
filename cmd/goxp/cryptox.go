package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/whitekid/goxp/cryptox"
)

func init() {
	cmd := &cobra.Command{
		Use: "cryptox",
	}

	cmd.AddCommand(&cobra.Command{
		Use:  "encrypt key data",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := cryptox.Encrypt(args[0], args[1])
			if err != nil {
				return err
			}

			fmt.Printf("%s\n", s)
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:  "decrypt key data",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := cryptox.Decrypt(args[0], args[1])
			if err != nil {
				return err
			}

			fmt.Printf("%s\n", s)
			return nil
		},
	})

	rootCmd.AddCommand(cmd)
}
