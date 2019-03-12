package gin

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/whitekid/go-utils/gin/combined"
	"github.com/whitekid/go-utils/log"
)

type Engine struct {
	*gin.Engine

	Addr *net.TCPAddr // real listening address
}

func New() *Engine {
	r := combined.NewRouter()

	return &Engine{Engine: r}
}

type HandlerFunc func(c *Context)

func Handler(fn HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := &Context{c}
		fn(ctx)
	}
}

// RunWithContext run http handler in separate goroutine(non blocking) with context to cancel
//
// Usage:
//	ctx, cancel := context.WithCancel(context.Background())
//  err := e.RunWithContext(ctx, "127.0.0.1:8080")
//	...
//
// to stop server
// 	cancel()
//
func (g *Engine) RunWithContext(ctx context.Context, address string) error {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		log.Debugf("Failed to listen: %s", err)
		return err
	}

	// when port is set to 0, port will be selected by random
	// so extract port information from listen address.
	g.Addr = ln.Addr().(*net.TCPAddr)
	log.Infof("Listening TCP %s", g.Addr.String())

	server := http.Server{
		Addr:    address,
		Handler: g,
	}

	go func() {
		if err := server.Serve(ln); err != http.ErrServerClosed {
			log.Infof("Listen %s", address)
		}
	}()

	go func() {
		<-ctx.Done()

		ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(ctxShutdown); err != nil {
			log.Error("Shutdown error: %s", err)
		}
	}()

	return nil
}

// type aliases
type IRoutes gin.IRoutes
type H gin.H

// RouteHandler sub router handler...
type RouteHandler interface {
	SetupRoute(r IRoutes)
}

// WrapGET nomatch에서 python 버전으로 proxy하는데, Get을 응용하는 HEAD, OPTIONS등의 request가
// python 버전으로 proxy되는 것을 방지.
//
// NOTE: proxy 기능이 사라지면 삭재
func WrapGET(r gin.IRoutes, path string, h gin.HandlerFunc) {
	for _, method := range []string{http.MethodGet, http.MethodHead, http.MethodOptions} {
		r.Handle(method, path, h)
	}
}
