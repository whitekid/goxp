package service

import (
	"context"
	"os"
	"os/signal"

	"github.com/whitekid/goxp/log"
)

// SetupSignal return context done when get system terminmaton signal
//
// returned context will be done when signals are invoked
// if signals is not given, default signals are os.Interrupt, os.Kill
//
// Special signals:
//
//	os.Interrupt: if get one more os.Interrupt, call os.Exit(1)
//	os.Kill: call os.Exit(1)
func SetupSignal(ctx context.Context, signals ...os.Signal) context.Context {
	c := make(chan os.Signal, 10)
	if len(signals) == 0 {
		signals = []os.Signal{os.Interrupt, os.Kill}
	}
	signal.Notify(c, signals...)
	log.Debugf("setup signals for %s", signals)

	termCtx, cancel := context.WithCancel(ctx)
	go func() {
		defer close(c)

		try := 0

	exit:
		for {
			select {
			case sig := <-c:
				log.Debugf("get signal: %s, trying to cancel context", sig)
				// go signal, done context
				cancel()
				switch sig {
				case os.Interrupt:
					// get interrupt one more, exiting...
					if try > 0 {
						log.Debugf("interrupt one more time. kill anyway...")
						os.Exit(1)
					}
					try++

				case os.Kill:
					// get kill, exiting
					log.Debug("get kill signal. exiting.")
					os.Exit(1)
				}
			case <-ctx.Done():
				log.Debugf("parent context done. exiting...")
				break exit
			}

		}
	}()

	return termCtx
}
