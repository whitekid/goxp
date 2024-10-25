package service

import (
	"context"
	"net"
)

// Interface the service interface
type Interface interface {
	// Run service and block until stop
	Serve(ctx context.Context) error
}

// AddrGetter for get listen address
type AddrGetter interface {
	GetAddr() *net.TCPAddr
}
