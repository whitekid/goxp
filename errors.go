package goxp

// Must panic if error
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

// Must2 panic if error
func Must2[T any](err error, t T) T {
	if err != nil {
		panic(err)
	}

	return t
}
