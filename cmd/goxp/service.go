package main

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/whitekid/goxp"
	"github.com/whitekid/goxp/cobrax"
	"github.com/whitekid/goxp/services"
)

func init() {
	cobrax.Add(rootCmd, &cobra.Command{
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
	}, nil)
}

type timerService struct{}

var _ services.Interface = (*timerService)(nil)

func newTimerService() services.Interface {
	return &timerService{}
}

func (s *timerService) Serve(ctx context.Context) error {
	go goxp.Every(ctx, time.Second, false, func(ctx context.Context) {
		if goxp.IsContextDone(ctx) {
			return
		}

		fmt.Printf("%s\n", time.Now().UTC().Format(time.RFC3339))
	})

	<-ctx.Done()
	return ctx.Err()
}
