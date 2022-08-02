package types

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEquals(t *testing.T) {
	for _, test := range [...]struct {
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
		require.Equal(t, test.equal, test.s1.Equals(test.s2))
	}
}

func TestContains(t *testing.T) {
	s := NewStrings([]string{"a", "b", "c"})

	for _, test := range [...]struct {
		e        string
		contains bool
	}{
		{"a", true},
		{"x", false},
	} {
		require.Equal(t, test.contains, s.Contains(test.e))

		slice := s.Slice()
		require.Equal(t, test.contains, Strings(slice).Contains(test.e))
	}
}

func TestAppend(t *testing.T) {
	s := NewStrings([]string{"a", "b", "c"}).Append("X")

	require.Equal(t, []string{"a", "b", "c", "X"}, s.Slice())
}

func TestRemoveAt(t *testing.T) {
	s := NewStrings([]string{"a", "b", "c"}).RemoveAt(1)

	require.Equal(t, []string{"a", "c"}, s.Slice())
}

func TestRemove(t *testing.T) {
	s := NewStrings([]string{"a", "b", "c"}).Remove("b")

	require.Equal(t, []string{"a", "c"}, s.Slice())
}

func TestReader(t *testing.T) {
	type args struct {
		s   []string
		sep string
	}
	tests := [...]struct {
		name string
		args args
		want string
	}{
		{"default", args{s: []string{"a", "b", "c"}}, "abc"},
		{"with sep", args{s: []string{"a", "b", "c"}, sep: ","}, "a,b,c"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewStrings(tt.args.s)
			b, _ := io.ReadAll(s.Reader(tt.args.sep))
			require.Equal(t, tt.want, string(b))
		})
	}
}

func TestJoin(t *testing.T) {
	s := NewStrings([]string{"b", "c", "a"})

	require.Equal(t, "b,c,a", s.Join(","))
}

func TestSort(t *testing.T) {
	s := Strings{"b", "c", "a"}
	s.Sort()

	require.Equal(t, []string{"a", "b", "c"}, s.Slice())
}

func TestMap(t *testing.T) {
	got := NewStrings([]string{"bx", "c", "a"}).Map(strings.ToUpper)

	require.Equal(t, []string{"BX", "C", "A"}, got.Slice())

}
