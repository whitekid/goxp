package gin

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/whitekid/go-utils/httptest"
)

func TestServer(t *testing.T) {
	r := New()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err := r.RunWithContext(ctx, "127.0.0.1:0")
	assert.NoError(t, err)
}

func TestHandler(t *testing.T) {
	r := New()

	r.GET("/", func(c *Context) { c.String(http.StatusOK, "GET") })
	r.POST("/", func(c *Context) { c.String(http.StatusOK, "POST") })
	r.PUT("/", func(c *Context) { c.String(http.StatusOK, "PUT") })
	r.DELETE("/", func(c *Context) { c.String(http.StatusOK, "DELETE") })
	r.PATCH("/", func(c *Context) { c.String(http.StatusOK, "PATCH") })
	r.OPTIONS("/", func(c *Context) { c.String(http.StatusOK, "OPTIONS") })

	for _, test := range []struct {
		method string
	}{
		{http.MethodGet}, {http.MethodPost}, {http.MethodPut},
	} {
		w := httptest.Request(r, test.method, "/", httptest.Options{})
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, w.Body.String(), test.method)
	}
}

func TestGroup(t *testing.T) {
	r := New()

	v1 := r.Group("v1")
	{
		v1.GET("/hello", func(c *Context) { c.String(http.StatusOK, "Hello v1") })
	}

	w := httptest.Request(r, http.MethodGet, "/v1/hello", httptest.Options{})
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Body.String(), "Hello v1")
}

func TestError(t *testing.T) {
	r := New()

	r.GET("/400", func(c *Context) { c.Error(http.StatusOK, NewHTTPError(400, "Bad Request")) })
	r.GET("/401", func(c *Context) { c.Error(401, errors.New("401")) })

	for _, test := range []struct {
		path string
		code int
	}{
		{"/400", 400}, {"/401", 401},
	} {
		w := httptest.Request(r, http.MethodGet, test.path, httptest.Options{})
		assert.Equal(t, test.code, w.Code)
	}
}

func TestToken(t *testing.T) {
	for _, test := range []struct {
		header string
		token  string
	}{
		{"token 1234", "1234"},
		{"Token  1234", "1234"},
		{"Token \t1234", "1234"},
	} {
		token := authToken(test.header)

		assert.Equal(t, test.token, token)
	}
}
