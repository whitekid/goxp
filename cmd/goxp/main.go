// Package main is goexp example runner
package main

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "goxp",
	Short: "goxp examples",
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
