package iterx

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIter(t *testing.T) {
	it := Of(0, 1, 2, 3, 4)

	for e := range it {
		_ = e
	}
}

func TestAll(t *testing.T) {
	for i, e := range Of(0, 1, 2, 3, 4).All() {
		require.Equal(t, i, e)
	}
}

func TestAppend(t *testing.T) {
	it1 := Of(0, 1, 2, 3, 4)
	it2 := Of(5, 6, 7, 8, 9)

	require.Equal(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, it1.Chain(it2).Collect())
}

func TestChunk(t *testing.T) {
	e := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	type args struct {
		size int
	}
	tests := [...]struct {
		name    string
		args    args
		wantErr bool
	}{
		{`valid`, args{2}, false},
		{`valid`, args{3}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Chunk(slices.Values(e), tt.args.size)

			chunks := 0
			for ss := range c {
				s := slices.Collect(ss)
				require.LessOrEqualf(t, len(s), tt.args.size, "%v", s)

				chunks++
			}
			require.Equal(t, (len(e)+tt.args.size-1)/tt.args.size, chunks)
		})
	}
}

func TestChunk2(t *testing.T) {
	e := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	type args struct {
		size int
	}
	tests := [...]struct {
		name    string
		args    args
		wantErr bool
	}{
		{`valid`, args{2}, false},
		{`valid`, args{3}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Chunk2(slices.All(e), tt.args.size)

			chunks := 0
			for ss := range c {
				s := []int{}
				for v1, v2 := range ss {
					s = append(s, v1)
					require.Equal(t, v1, v2)
				}
				require.LessOrEqualf(t, len(s), tt.args.size, "%v", s)

				chunks++
			}
			require.Equal(t, (len(e)+tt.args.size-1)/tt.args.size, chunks)
		})
	}
}
