package fx

import (
	"golang.org/x/exp/constraints"
)

func Sum[T constraints.Integer | constraints.Float](collection []T) (r T) {
	for _, v := range collection {
		r += v
	}
	return r
}

func Max[T constraints.Ordered](col []T) (r T) {
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

func MaxBy[T any](col []T, cmp func(a T, b T) bool) (r T) {
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

func Min[T constraints.Ordered](col []T) (r T) {
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

func MinBy[T any](col []T, cmp func(a T, b T) bool) (r T) {
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
