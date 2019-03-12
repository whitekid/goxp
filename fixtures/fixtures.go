package fixtures

// Fixture ...
type Fixture interface {
	Teardown()
}

type dummyFixture struct {
}

func (d *dummyFixture) Teardown() {

}

// Chain fixture chaining..
func Chain(fixtures ...Fixture) Fixture {
	return &chainFixture{
		fixtures: fixtures,
	}
}

type chainFixture struct {
	fixtures []Fixture
}

func (c *chainFixture) Teardown() {
	for _, f := range c.fixtures {
		f.Teardown()
	}
}
