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
	v := viper.New()

	cmd := &cobra.Command{
		Use:   "request url",
		Short: "request package example",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			resp, err := request.Get(args[0]).Do(cmd.Context())
			if err != nil {
				return err
			}

			if v.GetBool("verbose") {
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

	configs := map[string][]flags.Flag{
		"request": {
			{"verbose", "v", false, "verbose"},
		},
	}

	flags.InitDefaults(v, configs)
	flags.InitFlagSet(v, configs, "request", cmd.Flags())

	rootCmd.AddCommand(cmd)
}
