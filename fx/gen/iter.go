package gen

import "iter"

// Seq return generator from seq
func Seq[T any](s iter.Seq[T]) Gen[T] {
	next, stop := iter.Pull(s)

	return func() (T, bool) {
		v, ok := next()
		if !ok {
			defer stop()
		}

		return v, ok
	}
}

// Seq returns iter.Seq from generator
func (next Gen[T]) Seq() iter.Seq[T] {
	return func(yield func(T) bool) {
		for v, ok := next(); ok; v, ok = next() {
			if !yield(v) {
				return
			}
		}
	}
}

func Seq2[T, T2 any](s iter.Seq2[T, T2]) Gen2[T, T2] {
	next, stop := iter.Pull2(s)

	return func() (T, T2, bool) {
		v1, v2, ok := next()
		if !ok {
			defer stop()
		}

		return v1, v2, ok
	}
}
