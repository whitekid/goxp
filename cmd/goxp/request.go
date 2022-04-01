package main

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/whitekid/goxp/flags"
	"github.com/whitekid/goxp/request"
)

func init() {
	cmd := &cobra.Command{
		Use:  "request url",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			resp, err := request.Get(args[0]).Do(cmd.Context())

			if viper.GetBool("verbose") {
				fmt.Printf("%s\n", resp.Status)
				for k := range resp.Header {
					fmt.Printf("%s: %s\n", k, resp.Header.Get(k))
				}
				fmt.Printf("\n")
			}

			defer resp.Body.Close()
			io.Copy(os.Stdout, resp.Body)
			return err
		},
	}

	flags.InitFlagSet(nil, map[string][]flags.Flag{
		"request": {
			{"verbose", "v", false, "verbose"},
		},
	}, "request", cmd.Flags())

	rootCmd.AddCommand(cmd)
}
