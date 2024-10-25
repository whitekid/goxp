package goxp

// Must panic if error
func Must(err error) {
	if err != nil {
		panic(err)
	}
}
