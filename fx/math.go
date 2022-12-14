package fx

import (
	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Integer | constraints.Float | constraints.Complex
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

func Max[S ~[]T, T Ordered](col S) (r T) {
	if len(col) == 0 {
		return
	}

	r = col[0]

	for i := 1; i < len(col); i++ {
		if col[i] > r {
			r = col[i]
		}
	}

	return
}

func MaxBy[S ~[]T, T any](col S, cmp func(a T, b T) bool) (r T) {
	if len(col) == 0 {
		return
	}

	r = col[0]

	for i := 1; i < len(col); i++ {
		if cmp(col[i], r) {
			r = col[i]
		}
	}

	return
}

func Min[S ~[]T, T Ordered](col S) (r T) {
	if len(col) == 0 {
		return
	}

	r = col[0]

	for i := 1; i < len(col); i++ {
		if col[i] < r {
			r = col[i]
		}
	}

	return
}

func MinBy[S ~[]T, T any](col S, cmp func(a T, b T) bool) (r T) {
	if len(col) == 0 {
		return
	}

	r = col[0]

	for i := 1; i < len(col); i++ {
		if cmp(col[i], r) {
			r = col[i]
		}
	}

	return
}
