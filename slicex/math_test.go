package slicex

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSum(t *testing.T) {
	type args struct {
		items []any
	}
	tests := [...]struct {
		name string
		args args
		want any
	}{
		{"valid", args{[]any{1, 2, 3}}, 6},
		{"valid", args{[]any{1.2, 2.3, 3.5}}, 7.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch want := tt.want.(type) {
			case int:
				testSum(t, To[int](tt.args.items), want)
			case float64:
				testSum(t, To[float64](tt.args.items), want)
			default:
				require.Fail(t, "unsupported type", "%v (%T)", want, want)
			}
		})
	}
}

func testSum[S ~[]T1, T1 Number](t *testing.T, items S, want T1) {
	require.Equal(t, want, Sum(items))
}

func FuzzSum(f *testing.F) {
	f.Add(1, 2, 3, 4, 5)
	f.Fuzz(func(t *testing.T, n1, n2, n3, n4, n5 int) {
		testSum(t, []int{n1, n2, n3, n4, n5}, n1+n2+n3+n4+n5)
	})
}

func TestSumBy(t *testing.T) {
	type args struct {
		items []string
		sumBy func(s string) any
	}
	tests := [...]struct {
		name string
		args args
		want any
	}{
		{"valid", args{
			[]string{"1", "2", "3"}, func(s string) any {
				v, _ := strconv.Atoi(s)
				return int(v)
			}}, 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch want := tt.want.(type) {
			case int:
				testSumBy(t, tt.args.items, want, func(x string) int {
					v, _ := strconv.Atoi(x)
					return v
				})
			default:
				require.Fail(t, "unsupported type", "%v (%T)", want, want)
			}
		})
	}
}

func testSumBy[S ~[]T1, T1 any, T2 Number](t *testing.T, items S, want T2, sumByWrap func(T1) T2) {
	require.Equal(t, want, SumBy(items, sumByWrap))
}

func FuzzSumBy(f *testing.F) {
	f.Add("1", "2", "3", "4", "5")
	f.Fuzz(func(t *testing.T, n1, n2, n3, n4, n5 string) {
		items := []string{n1, n2, n3, n4, n5}
		want := 0
		for _, x := range items {
			v, _ := strconv.Atoi(x)
			want += int(v)

		}

		testSumBy(t, items, want, func(x string) int {
			v, _ := strconv.Atoi(x)
			return v
		})
	})
}

func TestScale(t *testing.T) {
	type args struct {
		items []any
		scale any
	}
	tests := [...]struct {
		name string
		args args
	}{
		{"valid int", args{[]any{1, 2, 3, 4}, 2}},
		{"valid float", args{[]any{1.1, 2.1, 3.1, 4.1}, 2.0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.args.scale.(type) {
			case int:
				testScale(t, To[int](tt.args.items), v)
			case float64:
				testScale(t, To[float64](tt.args.items), v)
			default:
				require.Fail(t, "unsupported type", "%v (%T)", v, v)
			}
		})
	}
}

func testScale[S ~[]R, R Number](t *testing.T, args S, scale R) {
	got := Scale(args, scale)

	for i := 0; i < len(got); i++ {
		require.Equal(t, args[i]*scale, got[i])
	}
}

func FuzzScale(f *testing.F) {
	f.Add(1, 2, 3, 4, 5)
	f.Fuzz(func(t *testing.T, n1, n2, n3, n4, scale int) {
		testScale(t, []int{n1, n2, n3, n4}, scale)
	})
}
