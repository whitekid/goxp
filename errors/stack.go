package errors

import (
	"fmt"
	"runtime"
)

type withStack struct {
	message string
	err     error
	stack   []uintptr
}

func (e *withStack) Error() string  { return e.message }
func (e *withStack) String() string { return e.message }
func (e *withStack) Unwrap() error  { return e.err }

func (e *withStack) Format(f fmt.State, c rune) {
	switch c {
	case 'v':
		if f.Flag('+') {
			e.formatWithStack(f)
		} else {
			fmt.Fprint(f, e.message)
		}
	case 's':
		fmt.Fprint(f, e.message)
	}
}

func (e *withStack) formatWithStack(f fmt.State) {
	fmt.Fprintf(f, "%s\n", e.message)
	
	// Format stack trace efficiently
	for _, pc := range e.stack {
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}
		file, line := fn.FileLine(pc)
		fmt.Fprintf(f, "\t%s\n\t\t%s:%d\n", fn.Name(), file, line)
	}

	// Handle wrapped errors recursively
	if e.err != nil {
		fmt.Fprint(f, "Caused by: ")
		if wrappedErr, ok := e.err.(*withStack); ok {
			// Recursively format wrapped errors
			wrappedErr.formatWithStack(f)
		} else {
			// Use the error's own formatting
			fmt.Fprintf(f, "%+v\n", e.err)
		}
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
	// If stack traces are disabled globally, use standard errors
	if !EnableStackTrace {
		if err == nil {
			return NewWithoutStack(message)
		}
		return WrapWithoutStack(err, message)
	}
	
	stackBuf := make([]uintptr, MaxStackDepth)
	length := runtime.Callers(skip, stackBuf[:])
	
	// Only allocate the exact size needed
	stack := make([]uintptr, length)
	copy(stack, stackBuf[:length])

	return &withStack{
		message: message,
		err:     err,
		stack:   stack,
	}
}

// HasStack returns true if the error or any error in its chain has stack trace information
func HasStack(err error) bool {
	for err != nil {
		if _, ok := err.(*withStack); ok {
			return true
		}
		err = Unwrap(err)
	}
	return false
}

// Cause returns the root cause of the error chain
func Cause(err error) error {
	for {
		underlying := Unwrap(err)
		if underlying == nil {
			return err
		}
		err = underlying
	}
}

// StackTrace returns the stack trace from the first error in the chain that has one
func StackTrace(err error) []uintptr {
	for err != nil {
		if ws, ok := err.(*withStack); ok {
			// Return a copy to prevent modification
			result := make([]uintptr, len(ws.stack))
			copy(result, ws.stack)
			return result
		}
		err = Unwrap(err)
	}
	return nil
}
