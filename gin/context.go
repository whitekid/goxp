package gin

import (
	"io"
	"net/http"
	"regexp"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

// Context is gin.Context wrapper
type Context struct {
	*gin.Context
}

const (
	keyContext = "context"
)

const (
	keyAuthToken   = ".auth.token"
	keyBearerToken = ".baerer.token"
)

var (
	tokenRegex       = regexp.MustCompile(`(?i)(token)\s+(?P<token>.+)\s*`)
	bearerTokenRegex = regexp.MustCompile(`(?i)(bearer)\s+(?P<token>.+)\s*`)
)

// AuthToken return token in header
// Authorization: token <token>
func (c *Context) AuthToken() string {
	value, ok := c.Get(keyAuthToken)
	if !ok {
		token := authToken(c.Request.Header.Get("Authorization"))

		c.Set(keyAuthToken, token)
		return token
	}

	return value.(string)
}

func authToken(value string) (token string) {
	m := tokenRegex.FindStringSubmatchIndex(value)
	if m != nil {
		token = value[m[4]:m[5]]
	}

	return token
}

// BearerToken return bearer token in header
// Authorization: bearer <token>
func (c *Context) BearerToken() string {
	value, ok := c.Get(keyBearerToken)
	if !ok {
		token := bearerToken(c.Request.Header.Get("Authorization"))

		c.Set(keyBearerToken, token)
		return token
	}

	return value.(string)
}

func bearerToken(value string) (token string) {
	m := bearerTokenRegex.FindStringSubmatchIndex(value)
	if m != nil {
		token = value[m[4]:m[5]]
	}

	return token
}

// Error return http errors
func (c *Context) Error(err error) error {
	resp := newErrorResponse(err)

	c.AbortWithStatusJSON(resp.Code, resp)

	return err
}

// FromReader write response from io.Reader
// It similar DataFromReader() but not required content-length
func (c *Context) FromReader(r io.Reader) {
	c.Stream(func(w io.Writer) bool {
		io.Copy(w, r)

		return false
	})
}

// Template ...
func (c *Context) Template(code int, template string, context map[string]interface{}) {
	tpl, err := pongo2.FromString(template)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if err := tpl.ExecuteWriter(context, c.Writer); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
}
