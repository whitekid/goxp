package fx

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKeys(t *testing.T) {
	v := map[int]string{1: "a", 2: "b", 3: "c"}
	keys := Keys(v)
	sort.Ints(keys)
	require.Equal(t, []int{1, 2, 3}, keys)
}

func TestValues(t *testing.T) {
	m := map[int]string{1: "a", 2: "b", 3: "c"}
	v := Values(m)
	sort.Strings(v)
	require.Equal(t, []string{"a", "b", "c"}, v)
}

func TestFilterMap(t *testing.T) {
	r := FilterMap(map[int]string{1: "a", 2: "b", 3: "c"}, func(k int, v string) bool { return k%2 == 0 })
	require.Equal(t, map[int]string{2: "b"}, r)
}

func TestForEachMap(t *testing.T) {
	r := map[string]int{}
	ForEachMap(map[int]string{1: "a", 2: "b", 3: "c"}, func(k int, v string) { r[v] = k })
	require.Equal(t, map[string]int{"a": 1, "b": 2, "c": 3}, r)
}

func TestMapKeys(t *testing.T) {
	r := MapKeys(map[int]string{1: "a", 2: "b", 3: "c"}, func(k int, v string) string { return v })
	require.Equal(t, map[string]string{"a": "a", "b": "b", "c": "c"}, r)
}

func TestMapValues(t *testing.T) {
	r := MapValues(map[int]string{1: "a", 2: "b", 3: "c"}, func(k int, v string) int { return k })
	require.Equal(t, map[int]int{1: 1, 2: 2, 3: 3}, r)
}

func TestMergeMap(t *testing.T) {
	m1 := map[int]string{1: "a", 2: "b"}
	m2 := map[int]string{3: "c", 4: "d"}
	require.Equal(t, map[int]string{1: "a", 2: "b", 3: "c", 4: "d"}, MergeMap(m1, m2))
}

func TestSampleMap(t *testing.T) {
	m := map[int]string{1: "a", 2: "b", 3: "c"}
	k, v := SampleMap(m)
	require.Equal(t, v, m[k])
}
