package fx

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSlice(t *testing.T) {
	s := NewSlice([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	s = s.Filter(func(x int) bool { return x%2 == 0 })
	require.Equal(t, []int{0, 2, 4, 6, 8}, s.Slice())
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

func TestSort(t *testing.T) {
	type args struct {
		s []int
	}
	tests := [...]struct {
		name string
		args args
		want []int
	}{
		{"valid", args{[]int{4, 3, 2, 1}}, []int{1, 2, 3, 4}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.False(t, IsSorted(tt.args.s))

			got := Clone(tt.args.s)
			got = Sort(got)
			require.Equal(t, tt.want, got)
			require.True(t, IsSorted(got))
		})
	}
}

func TestGroupBy(t *testing.T) {
	got := GroupBy([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}, func(e int) int { return e % 3 })
	require.Equal(t, map[int][]int{
		0: {3, 6, 9},
		1: {1, 4, 7},
		2: {2, 5, 8},
	}, got)
}

func TestPartitionBy(t *testing.T) {
	got := PartitionBy([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}, func(e int) int { return e % 3 })
	require.Equal(t, [][]int{
		{1, 4, 7},
		{2, 5, 8},
		{3, 6, 9},
	}, got)
}

func TestUniqBy(t *testing.T) {
	got := UniqBy([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}, func(e int) int { return e % 3 })
	require.Equal(t, []int{1, 2, 3}, got)
}

func TestChunk(t *testing.T) {
	got := Chunk([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, 4)
	require.Equal(t, [][]int{{0, 1, 2, 3}, {4, 5, 6, 7}, {8, 9}}, got)
}

func TestInterleave(t *testing.T) {
	type args struct {
		s [][]int
	}
	tests := [...]struct {
		name string
		args args
		want []int
	}{
		{`valid`, args{[][]int{}}, []int{}},
		{`valid`, args{[][]int{{1, 2}, {3, 4}}}, []int{1, 3, 2, 4}},
		{`valid`, args{[][]int{{1, 2}, {3}}}, []int{1, 3, 2}},
		{`valid`, args{[][]int{{1, 2}, {3, 4, 5}}}, []int{1, 3, 2, 4, 5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Interleave(tt.args.s...)
			require.Equal(t, tt.want, got, `Interleave() failed: got = %+v, want = %+v`, got, tt.want)
		})
	}
}

func TestReverse(t *testing.T) {
	got := Reverse([]int{0, 1, 2, 3, 4, 5, 6})
	require.Equal(t, []int{6, 5, 4, 3, 2, 1, 0}, got)
}

func TestToMap(t *testing.T) {
	got := ToMap([]int{0, 1, 2, 3, 4, 5}, func(e int) (int, int) { return e % 2, e })
	require.Equal(t, map[int]int{0: 4, 1: 5}, got)
}

func TestDrop(t *testing.T) {
	type args struct {
		s []int
		n int
	}
	tests := [...]struct {
		name string
		args args
		want []int
	}{
		{`valid`, args{}, nil},
		{`valid`, args{[]int{0, 1, 2, 3, 4, 5}, 2}, []int{2, 3, 4, 5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Drop(tt.args.s, tt.args.n)
			require.Equal(t, tt.want, got, `Drop() failed: got = %+v, want = %v`, got, tt.want)
		})
	}
}

func TestDropRight(t *testing.T) {
	type args struct {
		s []int
		n int
	}
	tests := [...]struct {
		name string
		args args
		want []int
	}{
		{`valid`, args{}, nil},
		{`valid`, args{[]int{0, 1, 2, 3, 4, 5}, 2}, []int{0, 1, 2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DropRight(tt.args.s, tt.args.n)
			require.Equal(t, tt.want, got, `DropRight() failed: got = %+v, want = %v`, got, tt.want)
		})
	}
}

func TestDropRightWhile(t *testing.T) {
	type args struct {
		s []int
		f func(int) bool
	}
	tests := [...]struct {
		name string
		args args
		want []int
	}{
		{`valid`, args{}, nil},
		{`valid`, args{[]int{0, 1, 2, 3, 4, 5}, func(e int) bool { return e != 2 }}, []int{0, 1, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DropRightWhile(tt.args.s, tt.args.f)
			require.Equal(t, tt.want, got, `DropRightWhile() failed: got = %+v, want = %v`, got, tt.want)
		})
	}
}

func TestDropWhile(t *testing.T) {
	type args struct {
		s []int
		f func(int) bool
	}
	tests := [...]struct {
		name string
		args args
		want []int
	}{
		{`valid`, args{}, nil},
		{`valid`, args{[]int{0, 1, 2, 3, 4, 5}, func(e int) bool { return e != 2 }}, []int{2, 3, 4, 5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DropWhile(tt.args.s, tt.args.f)
			require.Equal(t, tt.want, got, `DropWhile() failed: got = %+v, want = %v`, got, tt.want)
		})
	}
}

func TestReject(t *testing.T) {
	r := Reject([]int{1, 2, 3, 4}, func(v int) bool { return v%2 == 0 })
	require.Equal(t, []int{1, 3}, r)
}
