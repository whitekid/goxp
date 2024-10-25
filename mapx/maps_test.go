package mapx

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"
)

func TestFilter(t *testing.T) {
	r := Filter(map[int]string{1: "a", 2: "b", 3: "c"}, func(k int, v string) bool { return k%2 == 0 })
	require.Equal(t, map[int]string{2: "b"}, r)
}

func TestEach(t *testing.T) {
	r := map[string]int{}
	Each(map[int]string{1: "a", 2: "b", 3: "c"}, func(k int, v string) { r[v] = k })
	require.Equal(t, map[string]int{"a": 1, "b": 2, "c": 3}, r)
}

func TestMapKey(t *testing.T) {
	r := MapKey(map[int]string{1: "a", 2: "b", 3: "c"}, func(k int, v string) string { return v })
	require.Equal(t, map[string]string{"a": "a", "b": "b", "c": "c"}, r)
}

func TestMapValue(t *testing.T) {
	r := MapValue(map[int]string{1: "a", 2: "b", 3: "c"}, func(k int, v string) int { return k })
	require.Equal(t, map[int]int{1: 1, 2: 2, 3: 3}, r)
}

func TestSlice(t *testing.T) {
	r := Slice(map[int]string{1: "a", 2: "b", 3: "c"},
		func(k int, v string) string { return fmt.Sprintf("%d:%s", k, v) })
	slices.Sort(r)
	require.Equal(t, []string{"1:a", "2:b", "3:c"}, r)
}

func TestMerge(t *testing.T) {
	m1 := map[int]string{1: "a", 2: "b"}
	m2 := map[int]string{3: "c", 4: "d"}
	require.Equal(t, map[int]string{1: "a", 2: "b", 3: "c", 4: "d"}, Merge(m1, m2))
}

func TestSampleMap(t *testing.T) {
	m := map[int]string{1: "a", 2: "b", 3: "c"}
	k, v := Sample(m)
	require.Equal(t, v, m[k])
}
