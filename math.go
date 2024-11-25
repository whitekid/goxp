package goxp

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

func Abs[T RealNumber](n T) T { return T(math.Abs(float64(n))) }

func Min[T RealNumber](a, b T) T {
	if a > b {
		return b
	}

	return a
}

func Max[T RealNumber](a, b T) T {
	if a > b {
		return a
	}

	return b
}

func Percent[T1 RealNumber, T2 RealNumber](num T1, total T2) float64 {
	if total == 0 {
		return 0
	}

	return float64(num) * 100 / float64(total)
}
