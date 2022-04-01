package fx

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func dump[T any](collection []T) {
	fmt.Printf("dump: %v\n", collection)
}

func TestFX(t *testing.T) {
	s := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	ForEach(s, func(_, e int) { fmt.Printf("%d", e) })
	fmt.Printf("\n")

	s1 := Map(s, func(e int) int { return e * 2 })
	s2 := Filter(s1, func(e int) bool { return e%3 == 0 })
	s3 := Map(s2, func(e int) string { return strconv.FormatInt(int64(e), 10) })

	require.Equal(t, 45, Sum(s))
	require.Equal(t, 9, Max(s))
	require.Equal(t, 0, Min(s))

	require.Equal(t, []int{1, 2, 3, 4, 5, 6, 7}, Distinct([]int{1, 1, 1, 2, 3, 4, 5, 5, 6, 7}))

	dump(s3)
}

func TestFXObject(t *testing.T) {
	s := NewSlice(Times(10, func(i int) int { return i }))

	s1 := s.Map(func(e int) int { return e*2 + 1 }).
		Filter(func(e int) bool { return e%3 == 0 })

	s2 := Map(s1, func(e int) string { return strconv.FormatInt(int64(e), 10) })
	s3 := Reduce(s2, func(agg string, e string) string {
		if agg == "" {
			return e
		}
		return agg + "," + e
	})

	require.Equal(t, "3,9,15", s3)
}

func TestTernary(t *testing.T) {
	type args struct {
		value int
	}
	tests := [...]struct {
		name string
		args args
		want string
	}{
		{"even", args{10}, "even"},
		{"odd", args{11}, "odd"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Ternary(tt.args.value%2 == 0, "even", "odd")
			require.Equal(t, tt.want, got)
		})
	}
}

func TestIf(t *testing.T) {
	require.Equal(t, "true", If(func() bool { return true }, "true").Else("false"))
	require.Equal(t, "false", If(func() bool { return false }, "true").Else("false"))
}
