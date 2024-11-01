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
