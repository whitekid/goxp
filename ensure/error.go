package ensure

// NoError panic if err
func NoError(err error) {
	if err != nil {
		panic(err)
	}
}
