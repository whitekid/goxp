package fx

import (
	"iter"
	"slices"

	"github.com/whitekid/goxp/fx/gen"
)

// Seq iterator type
type Seq[T any] iter.Seq[T]

// Of return seq from slice
func Of[T any](v ...T) Seq[T] { return Seq[T](slices.Values(v)) }

// Gen return seq from Generator
func Gen[T any](next gen.Generator[T]) Seq[T] {
	return func(yield func(T) bool) {
		for {
			v, ok := next()
			if !ok {
				return
			}

			if !yield(v) {
				return
			}
		}
	}
}

func (s Seq[T]) All() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		i := 0
		for v := range s {
			if !yield(i, v) {
				return
			}
			i++
		}
	}
}

func (s Seq[T]) Append(s1 Seq[T]) Seq[T] {
	return func(yield func(T) bool) {
		for v := range s {
			if !yield(v) {
				return
			}
		}

		for v := range s1 {
			if !yield(v) {
				return
			}
		}
	}
}

func (s Seq[T]) Collect() []T { return slices.Collect(iter.Seq[T](s)) }

func (s Seq[T]) Each(fn func(int, T)) {
	for i, e := range s.All() {
		fn(i, e)
	}
}

func (s Seq[T]) Filter(fn func(T) bool) Seq[T] {
	return func(yield func(T) bool) {
		for e := range s {
			if !fn(e) {
				continue
			}

			if !yield(e) {
				return
			}
		}
	}
}

func (s Seq[T]) Iter() iter.Seq[T] { return iter.Seq[T](s) }

func Map[T1, T2 any](s Seq[T1], fn func(T1) T2) Seq[T2] {
	return func(yield func(T2) bool) {
		for e := range s {
			if !yield(fn(e)) {
				return
			}
		}
	}
}
