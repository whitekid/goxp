package fx

import (
	"cmp"
	"iter"
	"math/rand"
	"slices"
	"time"

	"github.com/whitekid/goxp/sets"
	"github.com/whitekid/goxp/slicex"
)

var rnd *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func Chunk[S ~[]E, E any](s S, size int) iter.Seq[S] {
	return slices.Chunk(s, size)
}

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
func Each[S ~[]E, E any](s S, fx func(int, E)) {
	for i, e := range s {
		fx(i, e)
	}
}

// EachE stop foreach if error
func EachE[S ~[]E, E any](s S, fx func(int, E) error) error {
	for i, e := range s {
		if err := fx(i, e); err != nil {
			return err
		}
	}
	return nil
}

// Filter return filtered slice
func Filter[S ~[]E, E any](s S, fx func(E) bool) S {
	if s == nil {
		return nil
	}

	r := make(S, len(s))

	j := 0
	Each(s, func(i int, e E) {
		if !fx(e) {
			return
		}
		r[j] = e
		j++
	})

	return r[:j]
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
func Map[S ~[]T1, T1 any, T2 any](s S, fx func(T1) T2) []T2 {
	if s == nil {
		return nil
	}

	r := make([]T2, len(s))
	Each(s, func(i int, e T1) { r[i] = fx(e) })

	return r
}

func Reduce[S ~[]E, E any](s S, fx func(E, E) E) (r E) {
	for _, e := range s {
		r = fx(r, e)
	}

	return r
}

func Reject[S ~[]E, E any](s S, f func(E) bool) S {
	if s == nil {
		return nil
	}

	r := make(S, len(s))

	j := 0
	Each(s, func(i int, e E) {
		if f(e) {
			return
		}
		r[j] = e
		j++
	})

	return r[:j]
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
func Times[T any](count int, f func(int) T) []T { return slicex.Times(count, f).Slice() }

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

func Find[S ~[]E, E any](s S, fx func(E) bool) (E, bool) {
	for _, e := range s {
		if fx(e) {
			return e, true
		}
	}

	var r E
	return r, false
}

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

func Sample[S ~[]E, E any](s S) E { return s[rnd.Intn(len(s))] }
func Samples[S ~[]E, E any](s S, count int) S {
	if s == nil {
		return nil
	}

	return Map(make(S, count), func(e E) E { return Sample(s) })
}

// Zip zip slice pair to mapping
// (key1, key2, key3), (values1, value2, values3) --> (key1: value1), (key2: value2), (key3: value3)
func Zip[K comparable, V any](keys []K, values []V) (r map[K]V) {
	r = make(map[K]V, len(keys))

	Each(keys, func(i int, k K) { r[k] = values[i] })

	return r
}

func Intersect[S ~[]E, E comparable](s1, s2 S) S {
	s := sets.New(s2)

	return Filter(s1, func(e E) bool {
		if s.Has(e) {
			return false
		}

		s.Append(e)
		return true
	})
}

func Concat[S ~[]E, E any](cols ...S) S { return slices.Concat(cols...) }

func Union[S ~[]E, E comparable](ss ...S) (r S) {
	r = make(S, 0)
	s := sets.New[E]()

	Each(Concat(ss...), func(i int, e E) {
		if s.Has(e) {
			return
		}

		s.Append(e)
		r = append(r, e)
	})

	return r
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

type Slice[S ~[]E, E any] []E

func NewSlice[E any](ss ...[]E) Slice[[]E, E] { return Concat(ss...) }
func Of[E any](ss ...E) Slice[[]E, E]         { return NewSlice(ss) }

func (s Slice[S, E]) Clip() Slice[S, E]                         { return slices.Clip(s) }
func (s Slice[S, E]) Clone() Slice[S, E]                        { return slices.Clone(s) }
func (s Slice[S, E]) ContainsFunc(f func(E) bool) bool          { return slices.ContainsFunc(s, f) }
func (s Slice[S, E]) CompactFunc(f func(E, E) bool) Slice[S, E] { return slices.CompactFunc(s, f) }
func (s Slice[S, E]) Delete(i, j int) Slice[S, E]               { return slices.Delete(s, i, j) }
func (s Slice[S, E]) Sample() E                                 { return Sample(s) }
func (s Slice[S, E]) Replace(i, j int, v ...E) Slice[S, E]      { return slices.Replace(s, i, j, v...) }
func (s Slice[S, E]) Reduce(fx func(E, E) E) E                  { return Reduce(s, fx) }
func (s Slice[S, E]) MinBy(cmp func(a E, b E) bool) E           { return MinBy(s, cmp) }
func (s Slice[S, E]) MaxBy(cmp func(a E, b E) bool) E           { return MaxBy(s, cmp) }
func (s Slice[S, E]) Insert(i int, v ...E) Slice[S, E]          { return slices.Insert(s, i, v...) }
func (s Slice[S, E]) IndexFunc(f func(E) bool) int              { return slices.IndexFunc(s, f) }
func (s Slice[S, E]) Grow(n int) Slice[S, E]                    { return slices.Grow(s, n) }
func (s Slice[S, E]) Flatten(ss ...[]E) Slice[S, E]             { return Concat(ss...) }
func (s Slice[S, E]) Find(fx func(E) bool) (E, bool)            { return Find(s, fx) }
func (s Slice[S, E]) Filter(fx func(E) bool) Slice[S, E]        { return Filter(s, fx) }
func (s Slice[S, E]) EqualFunc(s1 S, eq func(E, E) bool) bool   { return slices.EqualFunc(s, s1, eq) }
func (s Slice[S, E]) EachE(fx func(int, E) error) error         { return EachE(s, fx) }
func (s Slice[S, E]) Each(fx func(int, E))                      { Each(s, fx) }
func (s Slice[S, E]) Samples(count int) Slice[S, E]             { return Samples(s, count) }
func (s Slice[S, E]) Shuffle() Slice[S, E]                      { return Shuffle(s) }
func (s Slice[S, E]) Slice() S                                  { return S(s) }

func (s Slice[S, E]) SortFunc(less func(a, b E) int) Slice[S, E] {
	r := make(S, len(s))
	copy(r, s)
	SortFunc(r, less)
	return Slice[S, E](r)
}

func (s Slice[S, E]) SortStableFunc(less func(a, b E) int) S {
	r := make(S, len(s))
	copy(r, s)
	SortStableFunc(s, less)
	return r
}

// SliceC Slice for comparable
type SliceC[S ~[]E, E comparable] []E

func NewSliceC[E comparable](ss ...[]E) SliceC[[]E, E] { return Concat(ss...) }
func OfC[S ~[]E, E comparable](ss ...E) SliceC[[]E, E] { return NewSliceC(ss) }

func (s SliceC[S, E]) Contains(e E) bool           { return slices.Contains(s, e) }
func (s SliceC[S, E]) Compact() SliceC[S, E]       { return slices.Compact(s) }
func (s SliceC[S, E]) Distinct() SliceC[S, E]      { return SliceC[S, E](Distinct(S(s))) }
func (s SliceC[S, E]) Equal(s1 S) bool             { return slices.Equal(S(s), s1) }
func (s SliceC[S, E]) Every(sub S) bool            { return Every(S(s), sub) }
func (s SliceC[S, E]) Index(e E) int               { return slices.Index(s, e) }
func (s SliceC[S, E]) Intersect(s2 S) SliceC[S, E] { return SliceC[S, E](Intersect(S(s), s2)) }
func (s SliceC[S, E]) Slice() S                    { return S(s) }
func (s SliceC[S, E]) Union(ss ...S) SliceC[S, E] {
	return SliceC[S, E](Union(append([]S{S(s)}, ss...)...))
}

// Sort sort pure function version
func Sort[S ~[]T, T cmp.Ordered](s S) S {
	r := make(S, len(s))
	copy(r, s)
	slices.Sort(r)
	return r
}

// SortFunc sort pure function version
func SortFunc[S ~[]T, T any](s S, cmp func(a, b T) int) S {
	r := make(S, len(s))
	copy(r, s)
	slices.SortFunc(r, cmp)
	return r
}

// SortStableFunc sort pure function version
func SortStableFunc[S ~[]T, T any](s S, cmp func(a, b T) int) S {
	r := make(S, len(s))
	copy(r, s)
	slices.SortStableFunc(r, cmp)
	return r
}
