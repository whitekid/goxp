package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEquals(t *testing.T) {
	for _, test := range []struct {
		s1    Strings
		s2    Strings
		equal bool
	}{
		{
			s1:    Strings{"a", "b", "c"},
			s2:    Strings{"a", "b", "c"},
			equal: true,
		},
		{
			s1:    Strings{"a", "b"},
			s2:    Strings{"a", "b", "c"},
			equal: false,
		},
		{
			s1:    Strings{"a", "b", "c"},
			s2:    Strings{"a", "c", "b"},
			equal: false,
		},
	} {
		assert.Equal(t, test.equal, test.s1.Equals(test.s2))
	}
}
