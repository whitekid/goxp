package fx

import (
	"iter"
	"math"
	"math/rand"
	"slices"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAbs(t *testing.T) {
	type args struct {
		n any
	}
	tests := [...]struct {
		name string
		args args
		want any
	}{
		{`valid`, args{1}, 1},
		{`valid`, args{0}, 0},
		{`valid`, args{-1}, 1},
		{`valid`, args{1.2345}, 1.2345},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch want := tt.want.(type) {
			case int:
				testAbs(t, tt.args.n.(int), want)
			case float64:
				testAbs(t, tt.args.n.(float64), want)
			default:
				require.Failf(t, "unsupported type", "%v (%T)", want, want)
			}
		})
	}
}

func testAbs[T RealNumber](t *testing.T, v T, want T) {
	var zero T
	got := Abs(v)
	require.GreaterOrEqual(t, got, zero)
}

func BenchmarkAbs(b *testing.B) {
	n := rand.Float64()

	tests := [...]struct {
		name string
		abs  func(n float64) float64
	}{
		{"absWithFloat", func(n float64) float64 { return absWithFloat(n) }},
		{"abs", func(n float64) float64 { return abs(n) }},
	}
	for _, bb := range tests {
		b.Run(bb.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				bb.abs(n)
			}
		})
	}
}

func FuzzAbs(f *testing.F) {
	f.Add(1.2345)
	f.Fuzz(func(t *testing.T, n float64) {
		testAbs(t, n, math.Abs(n))
		testAbs(t, int(n), int(math.Abs(n)))
	})
}

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
				testSum(t, To[int](Iter(tt.args.items...)), want)
			case float64:
				testSum(t, To[float64](Iter(tt.args.items...)), want)
			default:
				require.Fail(t, "unsupported type", "%v (%T)", want, want)
			}
		})
	}
}

func testSum[T1 Number](t *testing.T, s iter.Seq[T1], want T1) {
	require.Equal(t, want, Sum(s))
}

func FuzzSum(f *testing.F) {
	f.Add(1, 2, 3, 4, 5)
	f.Fuzz(func(t *testing.T, n1, n2, n3, n4, n5 int) {
		testSum(t, Iter(n1, n2, n3, n4, n5), n1+n2+n3+n4+n5)
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
	require.Equal(t, want, SumBy(Iter(items...), sumByWrap))
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
				testScale(t, slices.Collect(To[int](Iter(tt.args.items...))), v)
			case float64:
				testScale(t, slices.Collect(To[float64](Iter(tt.args.items...))), v)
			default:
				require.Fail(t, "unsupported type", "%v (%T)", v, v)
			}
		})
	}
}

func testScale[R Number](t *testing.T, args []R, scale R) {
	got := slices.Collect(Scale(Iter(args...), scale))

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
