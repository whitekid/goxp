package fx

import (
	"fmt"
	"slices"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSlice(t *testing.T) {
	got := Of(0, 1, 2, 3, 4, 5, 6, 7, 8, 9).
		Filter(func(x int) bool { return x%2 == 0 }).
		Collect()
	require.Equal(t, []int{0, 2, 4, 6, 8}, got)
}

func TestForEach(t *testing.T) {
	got := []string{}
	Of("a", "b", "c", "d").Each(func(i int, v string) {
		got = append(got, fmt.Sprintf("%d:%s", i, v))
	})
	require.Equal(t, []string{"0:a", "1:b", "2:c", "3:d"}, got)
}

func TestForEachE(t *testing.T) {
	got := []string{}
	Of("a", "b", "c", "d").EachE(func(i int, v string) error {
		got = append(got, fmt.Sprintf("%d:%s", i, v))
		return nil
	})
	require.Equal(t, []string{"0:a", "1:b", "2:c", "3:d"}, got)
}

func TestFilter(t *testing.T) {
	got := Of(1, 2, 3, 4).Filter(func(v int) bool { return v%2 == 0 }).Collect()
	require.Equal(t, []int{2, 4}, got)
}

func TestMap(t *testing.T) {
	got := slices.Collect(Map(Iter(1, 2, 3, 4), func(v int) string { return strconv.FormatInt(int64(v), 10) }))
	require.Equal(t, []string{"1", "2", "3", "4"}, got)
}

func TestReduce(t *testing.T) {
	got := Of(1, 2, 3, 4).Reduce(func(x, y int) int { return x + y })
	require.Equal(t, 10, got)
}

func TestTimes(t *testing.T) {
	got := Times(10, func(v int) int { return v })
	require.Equal(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, got)
}

func TestShuffle(t *testing.T) {
	got := Shuffle(Of(1, 2, 3, 4, 5, 6, 7, 8, 9).Collect())
	require.NotEqual(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, got)
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

func TestEvery(t *testing.T) {
	require.True(t, Every([]int{1, 2, 3, 4, 5}, []int{2, 4}))
}

func TestConcat(t *testing.T) {
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
			got := Of(tt.args.cola...).Concat(Of(tt.args.colb...)).Collect()
			require.Equal(t, tt.want, got, `Concat() failed: want=%v, got=%v`, tt.want, got)
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
			got := slices.Collect(Union(Iter(tt.args.cola...), Iter(tt.args.colb...)))
			require.Equal(t, tt.want, got, `Union() failed: want=%v, got=%v`, tt.want, got)
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
