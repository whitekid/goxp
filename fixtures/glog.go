package fixtures

import "flag"

func GlogVerbose() Fixture {
	flag.Lookup("logtostderr").Value.Set("true")
	flag.Lookup("v").Value.Set("4")

	return &dummyFixture{}
}
