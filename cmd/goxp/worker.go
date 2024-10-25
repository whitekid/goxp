package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"runtime"
	"time"

	"github.com/spf13/cobra"
	"github.com/whitekid/goxp"
	"github.com/whitekid/goxp/log"
	"github.com/whitekid/goxp/slicex"
)

func init() {
	cmd := &cobra.Command{
		Use: "worker",
		Run: func(cmd *cobra.Command, args []string) {
			loggers := slicex.Times(runtime.NumCPU(), func(i int) log.Interface {
				return log.Named(fmt.Sprintf("worker %d", i))
			})

			ctx := cmd.Context()
			goxp.DoWithWorker(cmd.Context(), 0, func(i int) error {
				logger := loggers[i]

				logger.Infof("go for work~ %d", i)
				defer logger.Infof("call it a day %d", i)

				sleepMSec, _ := rand.Int(rand.Reader, big.NewInt(10))

				after := time.NewTimer(time.Duration(sleepMSec.Int64()) * time.Second)
				select {
				case <-ctx.Done():
					if !after.Stop() {
						go func() { <-after.C }()
					}

				case <-after.C:
				}

				return nil
			})
		},
	}

	cmd.AddCommand(&cobra.Command{
		Use: "every",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			go goxp.Every(ctx, time.Second, false, func() {
				fmt.Printf("%s\n", time.Now().Format(time.RFC3339))
			})
			<-ctx.Done()
		},
	})

	rootCmd.AddCommand(cmd)
}
