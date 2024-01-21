package goxp

import (
	"cmp"
	"slices"
)

func Max[E cmp.Ordered](s ...E) E { return slices.Max(s) }
func Min[E cmp.Ordered](s ...E) E { return slices.Min(s) }
