package fx

import "iter"

// To convert []any to []T
func To[T any](s iter.Seq[any]) iter.Seq[T] { return Map(s, func(x any) T { return x.(T) }) }
