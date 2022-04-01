package main

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/whitekid/goxp"
	"github.com/whitekid/goxp/service"
)

func init() {
	cmd := &cobra.Command{
		Use:   "service",
		Short: "service example",
		Run: func(cmd *cobra.Command, args []string) {
			srv := newTimerService()

			ttl := 10 * time.Second
			ctx, cancel := context.WithTimeout(cmd.Context(), ttl)
			defer cancel()

			fmt.Printf("timer will be terminated after %s\n", ttl)
			srv.Serve(ctx)
		},
	}
	rootCmd.AddCommand(cmd)
}

type timerService struct{}

var _ service.Interface = &timerService{} // interface guard

func newTimerService() service.Interface {
	return &timerService{}
}

func (s *timerService) Serve(ctx context.Context) error {
	goxp.Every(ctx, time.Second, func() error {
		if goxp.IsContextDone(ctx) {
			return ctx.Err()
		}

		fmt.Printf("%s\n", time.Now().UTC().Format(time.RFC3339))
		return nil
	})

	<-ctx.Done()
	return ctx.Err()
}
