package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"runtime"
	"time"

	"github.com/spf13/cobra"

	"github.com/whitekid/goxp"
	"github.com/whitekid/goxp/cobrax"
	"github.com/whitekid/goxp/log"
	"github.com/whitekid/goxp/slicex"
)

func init() {
	cmd := cobrax.Add(rootCmd, &cobra.Command{
		Use: "worker",
		Run: func(cmd *cobra.Command, args []string) {
			loggers := slicex.Times(runtime.NumCPU(), func(i int) log.Interface {
				return log.Named(fmt.Sprintf("worker %d", i))
			})

			goxp.DoWithWorker(cmd.Context(), 0, func(ctx context.Context, i int) error {
				logger := loggers[i]

				logger.Infof("go for work~ %d", i)
				defer logger.Infof("call it a day %d", i)

				sleepMSec, _ := rand.Int(rand.Reader, big.NewInt(10))

				after := time.NewTimer(time.Duration(sleepMSec.Int64()) * time.Second)
				select {
				case <-ctx.Done():
					return ctx.Err()

				case <-after.C:
					return nil
				}
			})
		},
	}, nil)

	cobrax.Add(cmd, &cobra.Command{
		Use: "every",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			go goxp.Every(ctx, time.Second, false, func(ctx context.Context) {
				fmt.Printf("%s\n", time.Now().Format(time.RFC3339))
			})
			<-ctx.Done()
		},
	}, nil)
}
