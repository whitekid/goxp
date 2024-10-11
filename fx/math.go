package fx

import (
	"iter"
	"math"

	"golang.org/x/exp/constraints"
)

type RealNumber interface {
	constraints.Integer | constraints.Float
}

type Number interface {
	RealNumber | constraints.Complex
}

func Abs[T RealNumber](n T) T { return absWithFloat(n) }

func abs[T RealNumber](n T) T {
	if n > 0 {
		return n
	}

	return -n
}

func absWithFloat[T RealNumber](n T) T {
	return T(math.Abs(float64(n)))
}

func Scale[T Number](s iter.Seq[T], f T) iter.Seq[T] { return Map(s, func(v T) T { return v * f }) }

func Sum[T Number](s iter.Seq[T]) (r T) {
	for v := range s {
		r += v
	}
	return r
}

func SumBy[T1 any, T2 Number](s iter.Seq[T1], f func(T1) T2) (r T2) {
	for v := range s {
		r += f(v)
	}
	return r
}

func MaxBy[T any](s iter.Seq[T], cmp func(T, T) bool) (r T) {
	for e := range s {
		if cmp(e, r) {
			r = e
		}
	}

	return
}

func MinBy[T any](value iter.Seq[T], cmp func(T, T) bool) (r T) {
	for e := range value {
		if !cmp(e, r) {
			r = e
		}
	}

	return
}
