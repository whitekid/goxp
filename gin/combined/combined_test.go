package combined

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type requestOpts struct {
}

func doRequest(r http.Handler, method, path string, opts requestOpts) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestGinCombinedLog(t *testing.T) {
	buffer := bytes.Buffer{}
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(New(&buffer))

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	for _, test := range []struct {
		path string
		code int
	}{
		{"/", http.StatusOK},
	} {
		buffer.Reset()
		w := doRequest(router, http.MethodGet, test.path, requestOpts{})

		require.Equal(t, test.code, w.Code)
		require.Contains(t, buffer.String(), fmt.Sprintf("%d", test.code))
		require.Contains(t, buffer.String(), "2019")
		log.Printf(buffer.String())
	}
}
