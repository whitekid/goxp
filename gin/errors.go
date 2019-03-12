package gin

import "fmt"

// HTTPError ...
type HTTPError interface {
	error
	Code() int
}

// NewHTTPError ...
func NewHTTPError(code int, format string, args ...interface{}) HTTPError {
	return &httpError{
		code,
		fmt.Errorf(format, args...),
	}
}

type httpError struct {
	code int
	error
}

func (e httpError) Code() int {
	return e.code
}
