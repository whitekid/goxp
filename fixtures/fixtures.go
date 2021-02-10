package fixtures

// Teardown teardown function
type Teardown func()

// Chain chaining fixtures
func Chain(tears ...Teardown) Teardown {
	return func() {
		for _, f := range tears {
			f()
		}
	}
}
