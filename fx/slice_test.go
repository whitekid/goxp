package fx

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSlice(t *testing.T) {
	s := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	ForEach(s, func(_, e int) { fmt.Printf("%d", e) })
	fmt.Printf("\n")

	s1 := Map(s, func(e int) int { return e * 2 })
	s2 := Filter(s1, func(e int) bool { return e%3 == 0 })
	s3 := Map(s2, func(e int) string { return strconv.FormatInt(int64(e), 10) })
	fmt.Printf("dump: %v\n", s3)

	require.Equal(t, 45, Sum(s))
	require.Equal(t, 9, Max(s))
	require.Equal(t, 0, Min(s))

	require.Equal(t, []int{1, 2, 3, 4, 5, 6, 7}, Distinct([]int{1, 1, 1, 2, 3, 4, 5, 5, 6, 7}))
}

func TestForEach(t *testing.T) {
	r := []string{}
	ForEach([]string{"a", "b", "c", "d"}, func(i int, v string) {
		r = append(r, fmt.Sprintf("%d:%s", i, v))
	})
	require.Equal(t, []string{"0:a", "1:b", "2:c", "3:d"}, r)
}

func TestForEachE(t *testing.T) {
	r := []string{}
	ForEachE([]string{"a", "b", "c", "d"}, func(i int, v string) error {
		r = append(r, fmt.Sprintf("%d:%s", i, v))
		return nil
	})
	require.Equal(t, []string{"0:a", "1:b", "2:c", "3:d"}, r)
}

func TestFilter(t *testing.T) {
	r := Filter([]int{1, 2, 3, 4}, func(v int) bool { return v%2 == 0 })
	require.Equal(t, []int{2, 4}, r)
}

func TestMap(t *testing.T) {
	r := Map([]int{1, 2, 3, 4}, func(v int) string { return strconv.FormatInt(int64(v), 10) })
	require.Equal(t, []string{"1", "2", "3", "4"}, r)
}

func TestReduce(t *testing.T) {
	r := Reduce([]int{1, 2, 3, 4}, func(x, y int) int { return x + y }, 0)
	require.Equal(t, 10, r)
}

func TestTimes(t *testing.T) {
	r := Times(10, func(v int) int { return v })
	require.Equal(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, r)
}

func TestShuffle(t *testing.T) {
	r := Shuffle([]int{1, 2, 3, 4, 5, 6, 7, 8, 9})
	require.NotEqual(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, r)
}

func TestDistinct(t *testing.T) {
	type args struct {
		col []int
	}
	tests := [...]struct {
		name string
		args args
		want []int
	}{
		{`valid`, args{[]int{1, 2, 3, 4, 4}}, []int{1, 2, 3, 4}},
		{`valid`, args{[]int{1, 2, 3, 4}}, []int{1, 2, 3, 4}},
		{`valid`, args{[]int{1, 2, 3, 3, 3, 3, 4}}, []int{1, 2, 3, 4}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Distinct(tt.args.col)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestContains(t *testing.T) {
	require.True(t, Contains([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}, 9))
}

func TestIndex(t *testing.T) {
	require.Equal(t, 5, Index([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}, 6))
}

func TestFind(t *testing.T) {
	v, ok := Find([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}, func(v int) bool { return v == 6 })
	require.True(t, ok)
	require.Equal(t, 6, v)
}

func TestEvery(t *testing.T) {
	require.True(t, Every([]int{1, 2, 3, 4, 5}, []int{2, 4}))
}

func TestSample(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	require.Contains(t, s, Sample(s))
}

func TestSamples(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for _, e := range Samples(s, 5) {
		require.Contains(t, s, e)
	}
}

func TestZip(t *testing.T) {
	r := Zip([]int{1, 2, 3}, []string{"a", "b", "c"})
	require.Equal(t, map[int]string{
		1: "a",
		2: "b",
		3: "c",
	}, r)
}

func TestIntersect(t *testing.T) {
	type args struct {
		cola []int
		colb []int
	}
	tests := [...]struct {
		name string
		args args
		want []int
	}{
		{`valid`, args{[]int{1, 2, 3, 4}, []int{3}}, []int{1, 2, 4}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Interset(tt.args.cola, tt.args.colb)
			require.Equal(t, tt.want, got, `Interset() failed: want=%v, got=%v`, tt.want, got)
		})
	}
}

func TestFlatten(t *testing.T) {
	type args struct {
		cola []int
		colb []int
	}
	tests := [...]struct {
		name string
		args args
		want []int
	}{
		{`valid`, args{[]int{1, 2}, []int{3, 4}}, []int{1, 2, 3, 4}},
		{`valid`, args{[]int{1, 2}, []int{2, 4}}, []int{1, 2, 2, 4}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Flatten(tt.args.cola, tt.args.colb)
			require.Equal(t, tt.want, got, `Flatten() failed: want=%v, got=%v`, tt.want, got)
		})
	}
}

func TestUnion(t *testing.T) {
	type args struct {
		cola []int
		colb []int
	}
	tests := [...]struct {
		name string
		args args
		want []int
	}{
		{`valid`, args{[]int{1, 2}, []int{3, 4}}, []int{1, 2, 3, 4}},
		{`valid`, args{[]int{1, 2}, []int{2, 4}}, []int{1, 2, 4}},
		{`valid`, args{[]int{}, []int{2, 4}}, []int{2, 4}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Union(tt.args.cola, tt.args.colb)
			require.Equal(t, tt.want, got, `Union() failed: want=%v, got=%v`, tt.want, got)
		})
	}
}
