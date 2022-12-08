package fx

import (
	"math/rand"
	"time"

	"golang.org/x/exp/slices"
)

var rnd *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

// ForEach iteraterate slice and apply function
func ForEach[S ~[]T, T any](s S, fx func(int, T)) {
	for i, e := range s {
		fx(i, e)
	}
}

// ForEachE stop foreach if error
func ForEachE[S ~[]T, T any](s S, fx func(int, T) error) error {
	for i, e := range s {
		if err := fx(i, e); err != nil {
			return err
		}
	}
	return nil
}

// Filter return filtered slice
func Filter[S ~[]T, T any](s S, fx func(T) bool) S {
	if s == nil {
		return nil
	}

	result := make([]T, len(s))

	j := 0
	ForEach(s, func(i int, e T) {
		if !fx(e) {
			return
		}
		result[j] = e
		j++
	})

	return result[:j]
}

// Map map element and return new type
func Map[S ~[]T, T any, R any](s S, fx func(T) R) []R {
	if s == nil {
		return nil
	}

	r := make([]R, len(s))
	ForEach(s, func(i int, e T) { r[i] = fx(e) })

	return r
}

func Reduce[S ~[]T, T any, R any](s S, fx func(r R, e T) R, initVal R) R {
	for _, e := range s {
		initVal = fx(initVal, e)
	}

	return initVal
}

// Times repeat count times
func Times[T any](count int, f func(int) T) []T {
	r := make([]T, count)

	for i := range r {
		r[i] = f(i)
	}

	return r
}

// Shuffle return shuffled slice
func Shuffle[S ~[]T, T any](s S) S {
	if s == nil {
		return nil
	}

	sf := make(S, len(s))
	copy(sf, s)

	rnd.Shuffle(len(sf), func(i, j int) { sf[i], sf[j] = sf[j], sf[i] })

	return sf
}

// Distinct return distinct slice
func Distinct[S ~[]T, T comparable](s S) S {
	if s == nil {
		return nil
	}

	set := NewSet(s)
	return set.Slice()
}

// Contains return true if e in slice
func Contains[S ~[]T, T comparable](s S, e T) bool         { return slices.Contains(s, e) }
func ContainsFunc[S ~[]T, T any](s S, f func(T) bool) bool { return slices.ContainsFunc(s, f) }
func Index[S ~[]T, T comparable](s S, e T) int             { return slices.Index(s, e) }
func IndexFunc[S ~[]T, T any](s S, f func(T) bool) int     { return slices.IndexFunc(s, f) }

func Find[S ~[]T, T any](s S, fx func(T) bool) (T, bool) {
	for _, e := range s {
		if fx(e) {
			return e, true
		}
	}

	var result T
	return result, false
}

func Insert[S ~[]T, T any](s S, i int, v ...T) S { return slices.Insert(s, i, v...) }
func Delete[S ~[]T, T any](s S, i, j int) S      { return slices.Delete(s, i, j) }
func Equal[E comparable](s1, s2 []E) bool        { return slices.Equal(s1, s2) }
func EqualFunc[E1, E2 any](s1 []E1, s2 []E2, eq func(E1, E2) bool) bool {
	return slices.EqualFunc(s1, s2, eq)
}
func Replace[S ~[]T, T any](s S, i, j int, v ...T) S      { return slices.Replace(s, i, j, v...) }
func Clone[S ~[]T, T any](s S) S                          { return slices.Clone(s) }
func Compact[S ~[]T, T comparable](s S) S                 { return slices.Compact(s) }
func CompactFunc[S ~[]T, T any](s S, f func(T, T) bool) S { return slices.CompactFunc(s, f) }

// Every return true if y is subset x
func Every[S ~[]T, T comparable](s1, s2 S) bool {
	set := NewSet(s1)

	for _, e := range s2 {
		if !set.Has(e) {
			return false
		}
	}

	return true
}

func Sample[S ~[]T, T any](s S) T { return s[rnd.Intn(len(s))] }
func Samples[S ~[]T, T any](s S, count int) S {
	if s == nil {
		return nil
	}

	return Map(make(S, count), func(e T) T { return Sample(s) })
}

// Zip zip slice pair to mapping
// (key1, key2, key3), (values1, value2, values3) --> (key1: value1), (key2: value2), (key3: value3)
func Zip[K comparable, V any](keys []K, values []V) (r map[K]V) {
	r = make(map[K]V)

	ForEach(keys, func(i int, k K) { r[k] = values[i] })

	return r
}

func Interset[S ~[]T, T comparable](s1, s2 S) S {
	s := NewSet(s2)

	return Filter(s1, func(e T) bool {
		if s.Has(e) {
			return false
		}

		s.Append(e)
		return true
	})
}

func Flatten[S ~[]T, T any](cols ...S) S {
	length := 0

	for i := range cols {
		length += len(cols[i])
	}

	r := make(S, 0, length)
	for i := range cols {
		r = append(r, cols[i]...)
	}

	return r
}

func Union[S ~[]T, T comparable](ss ...S) (r S) {
	r = make(S, 0)
	s := NewSet[T]()

	ForEach(Flatten(ss...), func(i int, e T) {
		if s.Has(e) {
			return
		}

		s.Append(e)
		r = append(r, e)
	})

	return r
}

type Slice[S ~[]T, T any] []T

func NewSlice[T any](ss ...[]T) Slice[[]T, T] { return Flatten(ss...) }

func (s Slice[S, T]) Slice() S                             { return S(s) }
func (s Slice[S, T]) ForEach(fx func(int, T))              { ForEach(s, fx) }
func (s Slice[S, T]) ForEachE(fx func(int, T) error) error { return ForEachE(s, fx) }
func (s Slice[S, T]) Filter(fx func(T) bool) Slice[S, T]   { return Filter(s, fx) }
func (s Slice[S, T]) Shuffle() Slice[S, T]                 { return Shuffle(s) }
func (s Slice[S, T]) Find(fx func(T) bool) (T, bool)       { return Find(s, fx) }
func (s Slice[S, T]) Sample() T                            { return Sample(s) }
func (s Slice[S, T]) Samples(count int) Slice[S, T]        { return Samples(s, count) }
func (s Slice[S, T]) Flatten(ss ...[]T) Slice[S, T]        { return Flatten(ss...) }
func (s Slice[S, T]) MaxBy(cmp func(a T, b T) bool) T      { return MaxBy(s, cmp) }
func (s Slice[S, T]) MinBy(cmp func(a T, b T) bool) T      { return MinBy(s, cmp) }

// SliceC Slice for comparable
type SliceC[S ~[]T, T comparable] []T

func NewSliceC[T comparable](ss ...[]T) SliceC[[]T, T] { return Flatten(ss...) }

func (s SliceC[S, T]) Slice() S                   { return S(s) }
func (s SliceC[S, T]) Distinct() S                { return Distinct(S(s)) }
func (s SliceC[S, T]) Contains(e T) bool          { return Contains(s, e) }
func (s SliceC[S, T]) Index(e T) int              { return Index(s, e) }
func (s SliceC[S, T]) Every(sub S) bool           { return Every(S(s), sub) }
func (s SliceC[S, T]) Interset(s2 S) SliceC[S, T] { return SliceC[S, T](Interset(S(s), s2)) }
func (s SliceC[S, T]) Union(ss ...S) SliceC[S, T] {
	r := Union(append([]S{S(s)}, ss...)...)
	return SliceC[S, T](r)
}
