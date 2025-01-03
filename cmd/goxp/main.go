// Package main is goexp example runner
package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/whitekid/goxp"
	"github.com/whitekid/goxp/cobrax"
)

var rootCmd = &cobra.Command{
	Use:     "goxp",
	Short:   "goxp examples",
	Version: fmt.Sprintf("%s %s %s %s", GitTag, GitBranch, GitDirty, BuildTime),
}

func init() {
	cobrax.Add(rootCmd, &cobra.Command{
		Use:   "randomstring length",
		Short: "generate random string",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			n, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			fmt.Printf("%s\n", goxp.RandomString(n))
			return nil
		},
	}, nil)

	cobrax.Add(rootCmd, &cobra.Command{
		Use:   "randomuuid",
		Short: "generate random uuid",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s\n", uuid.New().String())
		},
	}, nil)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
