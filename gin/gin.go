package gin

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/whitekid/go-utils/gin/combined"
	"github.com/whitekid/go-utils/logging"
)

var (
	log = logging.New()
)

// Engine ...
type Engine struct {
	*gin.Engine
}

// New ...
func New() *Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(combined.New(nil))

	return &Engine{router}
}

// HandlerFunc ...
type HandlerFunc func(c *Context)

func wrap(handlers ...HandlerFunc) (funcs []gin.HandlerFunc) {
	funcs = make([]gin.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		funcs[i] = func(c *gin.Context) {
			ctx := &Context{c}
			handler(ctx)
		}
	}

	return
}

// GET ...
func (g *Engine) GET(relativePath string, handlers ...HandlerFunc) gin.IRoutes {
	return g.Engine.Handle("GET", relativePath, wrap(handlers...)...)
}

// POST ...
func (g *Engine) POST(relativePath string, handlers ...HandlerFunc) gin.IRoutes {
	return g.Engine.Handle("POST", relativePath, wrap(handlers...)...)
}

// PUT ...
func (g *Engine) PUT(relativePath string, handlers ...HandlerFunc) gin.IRoutes {
	return g.Engine.Handle("PUT", relativePath, wrap(handlers...)...)
}

// DELETE ...
func (g *Engine) DELETE(relativePath string, handlers ...HandlerFunc) gin.IRoutes {
	return g.Engine.Handle("DELETE", relativePath, wrap(handlers...)...)
}

// PATCH ...
func (g *Engine) PATCH(relativePath string, handlers ...HandlerFunc) gin.IRoutes {
	return g.Engine.Handle("PATCH", relativePath, wrap(handlers...)...)
}

// OPTIONS ...
func (g *Engine) OPTIONS(relativePath string, handlers ...HandlerFunc) gin.IRoutes {
	return g.Engine.Handle("OPTIONS", relativePath, wrap(handlers...)...)
}

// Group ...
func (g *Engine) Group(relativePath string, handlers ...HandlerFunc) *RouterGroup {
	return &RouterGroup{g.Engine.Group(relativePath, wrap(handlers...)...)}
}

// RouterGroup ...
type RouterGroup struct {
	*gin.RouterGroup
}

// GET ...
func (g *RouterGroup) GET(relativePath string, handlers ...HandlerFunc) gin.IRoutes {
	return g.RouterGroup.Handle("GET", relativePath, wrap(handlers...)...)
}

// POST ...
func (g *RouterGroup) POST(relativePath string, handlers ...HandlerFunc) gin.IRoutes {
	return g.RouterGroup.Handle("POST", relativePath, wrap(handlers...)...)
}

// PUT ...
func (g *RouterGroup) PUT(relativePath string, handlers ...HandlerFunc) gin.IRoutes {
	return g.RouterGroup.Handle("PUT", relativePath, wrap(handlers...)...)
}

// DELETE ...
func (g *RouterGroup) DELETE(relativePath string, handlers ...HandlerFunc) gin.IRoutes {
	return g.RouterGroup.Handle("DELETE", relativePath, wrap(handlers...)...)
}

// PATCH ...
func (g *RouterGroup) PATCH(relativePath string, handlers ...HandlerFunc) gin.IRoutes {
	return g.RouterGroup.Handle("PATCH", relativePath, wrap(handlers...)...)
}

// OPTIONS ...
func (g *RouterGroup) OPTIONS(relativePath string, handlers ...HandlerFunc) gin.IRoutes {
	return g.RouterGroup.Handle("OPTIONS", relativePath, wrap(handlers...)...)
}

// RunWithContext ...
func (g *Engine) RunWithContext(ctx context.Context, address string) error {
	log.Infof("Listening TCP %s", address)
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

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
