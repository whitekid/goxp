package slicex

import (
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

func Scale[S ~[]T, T Number](s S, f T) S {
	return Map(s, func(v T) T { return v * f })
}

func Sum[S ~[]T, T Number](s S) (r T) {
	for _, v := range s {
		r += v
	}
	return r
}

func SumBy[S ~[]T1, T1 any, T2 Number](s S, f func(T1) T2) (r T2) {
	for _, v := range s {
		r += f(v)
	}
	return r
}

func MaxBy[S ~[]T, T any](value S, cmp func(a T, b T) bool) (r T) {
	if len(value) == 0 {
		return
	}

	r = value[0]

	for i := 1; i < len(value); i++ {
		if cmp(value[i], r) {
			r = value[i]
		}
	}

	return
}

func MinBy[S ~[]T, T any](value S, cmp func(a T, b T) bool) (r T) {
	if len(value) == 0 {
		return
	}

	r = value[0]

	for i := 1; i < len(value); i++ {
		if cmp(value[i], r) {
			r = value[i]
		}
	}

	return
}
