package gin

import "net/http"

// ValidatorFunc ...
type ValidatorFunc func(*Context) error

// Validate ...
func (c *Context) Validate(fns ...ValidatorFunc) error {
	for _, fn := range fns {
		if err := fn(c); err != nil {
			return err
		}
	}

	return nil
}

// ValidateFailed ...
func ValidateFailed(format string, args ...interface{}) error {
	return NewHTTPError(http.StatusBadRequest, format, args...)
}
