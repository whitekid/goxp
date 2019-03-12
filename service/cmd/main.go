package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/whitekid/go-utils/service"
)

type echoOptions struct {
	message string
}

func (o *echoOptions) RegisterFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&o.message, "message", "m", "Hello World", "Message to echo")
}

type echoService struct {
	options *echoOptions
}

func (s *echoService) Serve(ctx context.Context, args ...string) error {
	fmt.Printf("echo: %s\n", s.options.message)
	fmt.Printf("args: %q\n", args)
	return nil
}

func newEchoService(options *echoOptions) service.Service {
	return &echoService{
		options: options,
	}
}

func newEchoCommand(ctx context.Context) *cobra.Command {
	options := &echoOptions{}

	cmd := &cobra.Command{
		Aliases: []string{"echo"},
		Use:     "echo",
		Long:    "Echo Service",
		RunE: func(cmd *cobra.Command, args []string) error {
			svc := newEchoService(options)
			return svc.Serve(ctx, args...)
		},
	}

	options.RegisterFlags(cmd.Flags())

	return cmd
}

func main() {
	root := &cobra.Command{
		Long: "Nucleo Service",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	ctx := service.SetupSignal(context.Background())
	root.AddCommand(newEchoCommand(ctx))

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
