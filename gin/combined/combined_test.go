package combined

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
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

		assert.Equal(t, test.code, w.Code)
		assert.Contains(t, buffer.String(), fmt.Sprintf("%d", test.code))
		assert.Contains(t, buffer.String(), "2018")
		log.Printf(buffer.String())

		assert.Equal(t, "", buffer.String())
	}
}
