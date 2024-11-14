package errors

import "errors"

// standard errors package wrapper

var ErrUnsupported = errors.ErrUnsupported

func New(message string) error      { return errors.New(message) }
func Is(err, target error) bool     { return errors.Is(err, target) }
func As(err error, target any) bool { return errors.As(err, target) }
func Join(errs ...error) error      { return errors.Join(errs...) }
func Unwrap(err error) error        { return errors.Unwrap(err) }
