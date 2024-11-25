package goxp

import (
	"math"
	"math/rand"
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
		{"abs", func(n float64) float64 { return Abs(n) }},
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

func TestPercent(t *testing.T) {
	type args struct {
		a any
		b any
	}
	tests := [...]struct {
		name string
		args args
		want float64
	}{
		{`valid`, args{0, 0}, 0},
		{`valid`, args{1, 10}, 10},
		{`valid`, args{1.0, 10.0}, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.args.a.(type) {
			case int:
				testPercent(t, tt.args.a.(int), tt.args.b.(int), tt.want)
			case float64:
				testPercent(t, tt.args.a.(float64), tt.args.b.(float64), tt.want)
			default:
				require.Failf(t, "unsupported type", "%v (%T)", tt.args.a, tt.args.a)
			}
		})
	}
}

func testPercent[T RealNumber](t *testing.T, a T, b T, want float64) {
	got := Percent(a, b)
	require.Equal(t, want, got)
}

func FuzzPercent(f *testing.F) {
	f.Add(1.2345, 1.2345)
	f.Fuzz(func(t *testing.T, a, b float64) {
		testPercent(t, a, b, (a*100)/b)
	})
}
