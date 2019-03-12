package gin

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Context is gin.Context wrapper
type Context struct {
	*gin.Context
}

// Value ...
type Value struct {
	v *string
}

// Int return int
func (v *Value) Int(defaultValue int) int {
	if v.v == nil {
		return defaultValue
	}

	intVal, err := strconv.Atoi(*v.v)
	if err == nil {
		return intVal
	}

	return defaultValue
}

// Bool return value as bool
func (v *Value) Bool(defaultValue bool) bool {
	if v.v == nil {
		return defaultValue
	}

	boolValue := strings.ToLower(*v.v)
	return boolValue == "true" || boolValue == "1"
}

// Int64 return query value as int64
func (v *Value) Int64(defaultValue int64) int64 {
	if v.v == nil {
		return defaultValue
	}

	intVal, err := strconv.ParseInt(*v.v, 10, 64)
	if err == nil {
		return intVal
	}

	return defaultValue
}

// QueryV get query Value
func (c *Context) QueryV(key string) *Value {
	v, ok := c.GetQuery(key)
	if ok {
		return &Value{&v}
	}
	return &Value{nil}
}

// Error return http errors
func (c *Context) Error(code int, err error) {
	var message string

	if err != nil {
		message = err.Error()
		if v, ok := err.(HTTPError); ok {
			code = v.Code()
		}
	}

	log.Errorf("Error: %d %s", code, message)
	c.AbortWithStatusJSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}

const (
	keyAuthToken = ".auth.token"
)

var (
	tokenRegex = regexp.MustCompile(`(?i)(token)\s+(?P<token>\w+)`)
)

// AuthToken ...
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
