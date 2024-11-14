package main

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/whitekid/goxp/flags"
	"github.com/whitekid/goxp/requests"
)

func init() {
	v := viper.New()

	cmd := &cobra.Command{
		Use:   "request url",
		Short: "request package example",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			resp, err := requests.Get(args[0]).
				Header(requests.HeaderUserAgent, v.GetString("user-agent")).
				Do(cmd.Context())
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

	fs := cmd.PersistentFlags()
	flags.String(fs, "user-agent", "user-agent", "A", "goxp requests agent "+GitTag, "use agent")
	flags.Bool(fs, "verbose", "verbose", "A", false, "verbose")

	rootCmd.AddCommand(cmd)
}
