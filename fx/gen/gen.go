package gen

import (
	"maps"
	"slices"
)

// Gen generate next value, return false if no values
type Gen[T any] func() (T, bool)

// Next create new Generator
func Next[T any](next func() func() (T, bool)) Gen[T] { return next() }

// Slice return generator from slice
func Slice[T any](s []T) Gen[T] { return Seq(slices.Values(s)) }

type Gen2[T1, T2 any] func() (T1, T2, bool)

func Next2[T1, T2 any](next func() func() (T1, T2, bool)) Gen2[T1, T2] { return next() }

// Map return generator from map
func Map[K comparable, V any](m map[K]V) Gen2[K, V] { return Seq2(maps.All(m)) }
