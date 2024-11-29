package errors

import (
	"fmt"
	"runtime"
)

type withStack struct {
	messsage string
	err      error
	stack    []uintptr
}

func (e *withStack) Error() string  { return e.messsage }
func (e *withStack) String() string { return e.messsage }
func (e *withStack) Unwrap() error  { return e.err }

func (e *withStack) Format(f fmt.State, c rune) {
	switch c {
	case 'v':
		if f.Flag('+') {
			fmt.Fprintf(f, "%s\n", e.messsage)
			for _, pc := range e.stack {
				fn := runtime.FuncForPC(pc)
				if fn == nil {
					continue
				}
				file, line := fn.FileLine(pc)
				fmt.Fprintf(f, "\t%s\n\t\t%s:%d\n", fn.Name(), file, line)
			}

			if e.err != nil {
				if wrappedErr, ok := e.err.(*withStack); ok {
					fmt.Fprintf(f, "Caused by:\n%+v", wrappedErr)
				} else {
					fmt.Fprintf(f, "Caused by: %+v\n", e.err)
				}
			}
		} else {
			fmt.Fprintf(f, "%v", e.messsage)
		}
	case 's':
		fmt.Fprintf(f, "%s", e.messsage)
	}
}

func Errorf(err error, format string, args ...any) error {
	return wrap(err, fmt.Sprintf(format, args...), 3)
}

func Wrapf(err error, format string, args ...any) error {
	return wrap(err, fmt.Sprintf(format, args...), 3)
}

func Wrap(err error, message string) error {
	return wrap(err, message, 3)
}

func wrap(err error, message string, skip int) error {
	stackBuf := make([]uintptr, 32)
	length := runtime.Callers(skip, stackBuf[:])
	stackBuf = stackBuf[:length]

	return &withStack{
		messsage: message,
		err:      err,
		stack:    stackBuf,
	}
}
