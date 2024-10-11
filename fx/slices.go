package fx

import (
	"iter"
	"math/rand"
	"slices"
	"time"

	"github.com/whitekid/goxp/sets"
)

var rnd *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func Count[S ~[]E, E comparable](s S, v E) (r int) {
	for _, e := range s {
		if e == v {
			r++
		}
	}
	return
}

func CountBy[S ~[]E, E any](s S, f func(E) bool) (r int) {
	for _, e := range s {
		if f(e) {
			r++
		}
	}
	return
}

func CountValues[S ~[]E, E comparable](s S) map[E]int {
	r := make(map[E]int)

	for _, e := range s {
		r[e]++
	}

	return r
}

func CountValuesBy[S ~[]E, E any, U comparable](s S, f func(E) U) map[U]int {
	r := make(map[U]int)

	for _, e := range s {
		r[f(e)]++
	}

	return r
}

func Drop[S ~[]E, E any](s S, n int) S      { return slices.Delete(s, 0, n) }
func DropRight[S ~[]E, E any](s S, n int) S { return slices.Delete(s, len(s)-n, len(s)) }

func DropRightWhile[S ~[]E, E any](s S, f func(E) bool) S {
	if s == nil {
		return nil
	}

	length := len(s)
	i := length - 1
	for ; i >= 0; i-- {
		if !f(s[i]) {
			break
		}
	}

	return slices.Delete(s, length-i-1, length)
}

func DropWhile[S ~[]E, E any](s S, f func(E) bool) S {
	if s == nil {
		return nil
	}

	i := 0
	for ; i < len(s); i++ {
		if !f(s[i]) {
			break
		}
	}

	return slices.Delete(s, 0, i)
}

// Each iteraterate slice and apply function
func Each[E any](s iter.Seq[E], fx func(int, E)) {
	i := 0
	for e := range s {
		fx(i, e)
		i++
	}
}

// EachE stop foreach if error
func EachE[E any](s iter.Seq[E], fx func(int, E) error) error {
	i := 0
	for e := range s {
		if err := fx(i, e); err != nil {
			return err
		}
		i++
	}
	return nil
}

// Filter return filtered slice
func Filter[E any](s iter.Seq[E], fx func(E) bool) iter.Seq[E] {
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

func Interleave[S ~[]E, E any](ss ...S) S {
	max := 0
	length := 0

	for _, s := range ss {
		length += len(s)
		if len(s) > max {
			max = len(s)
		}
	}

	r := make(S, 0, length)
	for i := 0; i < max; i++ {
		for j := range ss {
			if i > len(ss[j])-1 {
				continue
			}
			r = append(r, ss[j][i])
		}
	}

	return r
}

func GroupBy[S ~[]E, E any, K comparable](s S, f func(E) K) map[K]S {
	r := make(map[K]S)

	for _, x := range s {
		k := f(x)
		r[k] = append(r[k], x)
	}

	return r
}

func PartitionBy[S ~[]E, E any, K comparable](s S, f func(E) K) []S {
	seen := make(map[K]int)
	r := make([]S, 0)

	for _, x := range s {
		k := f(x)
		idx, ok := seen[k]
		if !ok {
			r = append(r, make(S, 0))
			idx = len(r) - 1
			seen[k] = idx
		}
		r[idx] = append(r[idx], x)
	}

	return r
}

// Map map element and return new type
func Map[T1 any, T2 any](s iter.Seq[T1], fx func(T1) T2) iter.Seq[T2] {
	return func(yield func(T2) bool) {
		for e := range s {
			if !yield(fx(e)) {
				return
			}
		}
	}
}

func Reduce[E any](s iter.Seq[E], fx func(E, E) E) (r E) {
	for e := range s {
		r = fx(r, e)
	}

	return r
}

func Reverse[S ~[]E, E any](s S) S {
	if s == nil {
		return nil
	}

	length := len(s)
	r := make(S, length)
	copy(r, s)

	for i := 0; i < length/2; i++ {
		j := length - i - 1
		r[i], r[j] = r[j], r[i]
	}

	return r
}

// Times repeat count times
func Times[T any](count int, f func(int) T) []T {
	r := make([]T, count)

	for i := range r {
		r[i] = f(i)
	}

	return r
}

func ToMap[S ~[]E, E any, K comparable, V any](s S, f func(E) (K, V)) map[K]V {
	r := make(map[K]V)
	for _, e := range s {
		k, v := f(e)
		r[k] = v
	}
	return r
}

// Shuffle return shuffled slice
func Shuffle[S ~[]E, E any](s S) S {
	if s == nil {
		return nil
	}

	r := make(S, len(s))
	copy(r, s)

	rnd.Shuffle(len(r), func(i, j int) { r[i], r[j] = r[j], r[i] })

	return r
}

// Distinct return distinct slice
func Distinct[S ~[]E, E comparable](s S) S {
	if s == nil {
		return nil
	}

	set := sets.New(s)
	return set.Slice()
}
func Uniq[S ~[]E, E comparable](s S) S { return Distinct(s) }

// Every return true if y is subset x
func Every[S ~[]E, E comparable](s1, s2 S) bool {
	set := sets.New(s1)

	for _, e := range s2 {
		if !set.Has(e) {
			return false
		}
	}

	return true
}

func Flatten[E any](cols ...iter.Seq[E]) iter.Seq[E] {
	return func(yield func(E) bool) {
		for _, c := range cols {
			for e := range c {
				if !yield(e) {
					return
				}
			}
		}
	}
}

func Union[E comparable](ss ...iter.Seq[E]) (r iter.Seq[E]) {
	return func(yield func(E) bool) {
		keys := make(map[E]struct{})

		for _, s := range ss {
			for e := range s {
				if _, ok := keys[e]; ok {
					continue
				}

				if !yield(e) {
					return
				}
				keys[e] = struct{}{}
			}
		}
	}
}

func UniqBy[S ~[]E, E any, U comparable](s S, f func(E) U) []E {
	r := []E{}
	seen := map[U]struct{}{}

	for _, e := range s {
		k := f(e)

		if _, ok := seen[k]; ok {
			continue
		}
		r = append(r, e)
		seen[k] = struct{}{}
	}

	return r
}

type Slice[E any] iter.Seq[E]

func Of[E any](ss ...E) Slice[E]      { return Slice[E](slices.Values(ss)) }
func Iter[E any](ss ...E) iter.Seq[E] { return slices.Values(ss) }

func (s Slice[E]) Seq() iter.Seq[E]         { return iter.Seq[E](s) }
func (s Slice[E]) Reduce(fx func(E, E) E) E { return Reduce(s.Seq(), fx) }

func (s Slice[E]) Concat(ss ...Slice[E]) Slice[E] {
	return func(yield func(E) bool) {
		for e := range s {
			if !yield(e) {
				return
			}
		}

		for _, s1 := range ss {
			for e := range s1 {
				if !yield(e) {
					return
				}

			}
		}
	}
}

func (s Slice[E]) Filter(fx func(E) bool) Slice[E]   { return Slice[E](Filter(s.Seq(), fx)) }
func (s Slice[E]) EachE(fx func(int, E) error) error { return EachE(s.Seq(), fx) }
func (s Slice[E]) Each(fx func(int, E))              { Each(s.Seq(), fx) }
func (s Slice[E]) Collect() []E                      { return slices.Collect(s.Seq()) }
