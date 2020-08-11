package service

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// SetupSignal return context done when get system terminmaton signal
//
// returned context will be done when os.Interrupt, syscall.SIGTERM are invoked
// and call os.Exit() to exit program
func SetupSignal(ctx context.Context) context.Context {
	c := make(chan os.Signal, 2)

	signal.Notify(c, []os.Signal{os.Interrupt, syscall.SIGTERM}...)

	termCtx, cancel := context.WithCancel(ctx)
	go func() {
		defer close(c)
		sig := <-c

		log.Debugf("got signal: %s", sig)

		cancel()
		os.Exit(1)
	}()

	return termCtx
}
