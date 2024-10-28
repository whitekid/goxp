package iterx

import (
	"cmp"
	"iter"
	"slices"
)

type Seq[E any] iter.Seq[E]

func Of[T any](s ...T) Seq[T] {
	return func(yield func(T) bool) {
		for _, e := range s {
			if !yield(e) {
				return
			}
		}
	}
}

func (s Seq[E]) Iter() iter.Seq[E] { return iter.Seq[E](s) }

func (s Seq[E]) All() iter.Seq2[int, E] { return All(s.Iter()) }
func All[E any](s iter.Seq[E]) iter.Seq2[int, E] {
	return func(yield func(int, E) bool) {
		i := 0
		for v := range s {
			if !yield(i, v) {
				return
			}
			
			i++
		}
	}
}

func (s Seq[E]) Append(seq Seq[E]) Seq[E] {
	return func(yield func(E) bool) {
		for e := range s {
			if !yield(e) {
				return
			}
		}

		for e := range seq {
			if !yield(e) {
				return
			}
		}
	}
}

func (s Seq[E]) Collect() []E { return slices.Collect(s.Iter()) }

func Count[T comparable](s iter.Seq[T], v T) int {
	r := 0

	for e := range s {
		if e == v {
			r++
		}
	}

	return r
}

func (s Seq[E]) Chunk(n int) iter.Seq[iter.Seq[E]] { return Chunk(s.Iter(), n) }

func Chunk[E any](s iter.Seq[E], n int) iter.Seq[iter.Seq[E]] {
	return func(yield func(iter.Seq[E]) bool) {
		next, stop := iter.Pull(s)
		defer stop()

		stopIt := false

		for !stopIt {
			v, ok := next()
			if !ok {
				return
			}

			if !yield(func(yield func(E) bool) {
				if !yield(v) {
					return
				}

				for i := 0; i < n-1; i++ {
					v, ok := next()
					if !ok {
						stopIt = true
						return
					}
					if !yield(v) {
						stopIt = true
						return
					}
				}
			}) {
				return
			}
		}
	}
}

func Chunk2[E1, E2 any](s iter.Seq2[E1, E2], n int) iter.Seq[iter.Seq2[E1, E2]] {
	return func(yield func(iter.Seq2[E1, E2]) bool) {
		next, stop := iter.Pull2(s)
		defer stop()

		stopIt := false

		for !stopIt {
			v1, v2, ok := next()
			if !ok {
				return
			}

			if !yield(func(yield func(E1, E2) bool) {
				if !yield(v1, v2) {
					return
				}

				for i := 0; i < n-1; i++ {
					v1, v2, ok := next()
					if !ok {
						stopIt = true
						return
					}
					if !yield(v1, v2) {
						stopIt = true
						return
					}
				}
			}) {
				return
			}
		}
	}
}

func (s Seq[E]) Concat(ss ...Seq[E]) Seq[E] { return Concat(append([]Seq[E]{s}, ss...)...) }

func Concat[E any](s ...Seq[E]) Seq[E] {
	return func(yield func(E) bool) {
		for _, ss := range s {
			for e := range ss {
				if !yield(e) {
					return
				}
			}
		}
	}
}

func (s Seq[E]) Each(fx func(int, E)) {
	for i, e := range s.All() {
		fx(i, e)
	}
}

func (s Seq[E]) Filter(fx func(E) bool) Seq[E] {
	return func(yield func(E) bool) {
		for e := range s {
			if fx(e) {
				if !yield(e) {
					return
				}
			}
		}
	}
}

func Map[T1, T2 any](s Seq[T1], fx func(T1) T2) Seq[T2] {
	return func(yield func(T2) bool) {
		for e := range s {
			if !yield(fx(e)) {
				return
			}
		}
	}
}

func (s Seq[E]) Reduce(fx func(E, E) E) E {
	var r E

	for e := range s {
		r = fx(r, e)
	}

	return r
}

func Repeat[E any](v E, count int) Seq[E] {
	return func(yield func(E) bool) {
		for i := 0; i < count; i++ {
			if !yield(v) {
				return
			}
		}
	}
}

func Sorted[E cmp.Ordered](seq Seq[E]) []E { return slices.Sorted(seq.Iter()) }
