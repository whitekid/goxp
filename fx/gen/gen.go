package gen

import (
	"maps"
	"slices"
)

// Generator generate next value, return false if no values
type Generator[T any] func() (T, bool)

func Next[T any](next func() func() (T, bool)) Generator[T] { return next() }

// Slice return generator from slice
func Slice[T any](s []T) Generator[T] { return Seq(slices.Values(s)) }

type Next2[T1, T2 any] func() (T1, T2, bool)

// Map return generator from map
func Map[K comparable, V any](m map[K]V) Next2[K, V] { return Seq2(maps.All(m)) }
