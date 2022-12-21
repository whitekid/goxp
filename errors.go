package goxp

import "errors"

func ErrorAs[T error](err error) (T, bool) {
	var ee T
	ok := errors.As(err, &ee)
	return ee, ok
}
