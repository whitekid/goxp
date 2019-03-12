package combined

import (
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

// New create new combined logger
func New(writer io.Writer) gin.HandlerFunc {
	if writer == nil {
		writer = gin.DefaultWriter
	}

	return func(c *gin.Context) {
		// process request
		c.Next()

		clientIP := c.ClientIP()
		ident := "-"
		remoteUser := c.Request.Header.Get("REMOTE_USER")
		serverProcessedTIme := time.Now().Format("02/Jan/2006:15:04:05 -0700")

		path := ""
		if c.Request.URL.RawQuery != "" {
			path = c.Request.URL.Path + "?" + c.Request.URL.RawQuery
		} else {
			path = c.Request.URL.Path
		}
		request := fmt.Sprintf("%s %s %s", c.Request.Method, path, c.Request.Proto)
		statusCode := c.Writer.Status()
		size := c.Writer.Size()
		referrer := c.Request.Header.Get("Referer")
		userAgent := c.Request.Header.Get("User-Agent")
		fmt.Fprintf(writer, `%s %s %s [%s] "%s" %d %d "%s" "%s"`+"\n",
			clientIP, ident, remoteUser, serverProcessedTIme, request,
			statusCode, size, referrer, userAgent)
	}
}
