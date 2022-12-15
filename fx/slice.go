package fx

import (
	"math/rand"
	"time"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

var rnd *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func Chunk[S ~[]T, T any](s S, size int) []S {
	num := len(s) / size
	if len(s)%size != 0 {
		num++
	}

	r := make([]S, 0, num)
	for i := 0; i < num; i++ {
		last := (i + 1) * size
		if last > len(s) {
			last = len(s)
		}

		r = append(r, s[i*size:last])
	}
	return r
}

func Count[S ~[]T, T comparable](s S, v T) (r int) {
	for _, e := range s {
		if e == v {
			r++
		}
	}
	return
}

func CountBy[S ~[]T, T any](s S, f func(T) bool) (r int) {
	for _, e := range s {
		if f(e) {
			r++
		}
	}
	return
}

func CountValues[S ~[]T, T comparable](s S) map[T]int {
	r := make(map[T]int)

	for _, e := range s {
		r[e]++
	}

	return r
}

func CountValuesBy[S ~[]T, T any, U comparable](s S, f func(T) U) map[U]int {
	r := make(map[U]int)

	for _, e := range s {
		r[f(e)]++
	}

	return r
}

func Drop[S ~[]T, T any](s S, n int) S      { return Delete(s, 0, n) }
func DropRight[S ~[]T, T any](s S, n int) S { return Delete(s, len(s)-n, len(s)) }

func DropRightWhile[S ~[]T, T any](s S, f func(T) bool) S {
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

	return Delete(s, length-i-1, length)
}

func DropWhile[S ~[]T, T any](s S, f func(T) bool) S {
	if s == nil {
		return nil
	}

	i := 0
	for ; i < len(s); i++ {
		if !f(s[i]) {
			break
		}
	}

	return Delete(s, 0, i)
}

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

	r := make(S, len(s))

	j := 0
	ForEach(s, func(i int, e T) {
		if !fx(e) {
			return
		}
		r[j] = e
		j++
	})

	return r[:j]
}

func Interleave[S ~[]T, T any](ss ...S) S {
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

func GroupBy[S ~[]T, T any, K comparable](s S, f func(T) K) map[K]S {
	r := make(map[K]S)

	for _, x := range s {
		k := f(x)
		r[k] = append(r[k], x)
	}

	return r
}

func PartitionBy[S ~[]T, T any, K comparable](s S, f func(T) K) []S {
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
	ForEach(s, func(i int, e T1) { r[i] = fx(e) })

	return r
}

func Reduce[S ~[]T1, T1 any, T2 any](s S, fx func(r T2, e T1) T2, initVal T2) T2 {
	for _, e := range s {
		initVal = fx(initVal, e)
	}

	return initVal
}

func Reject[S ~[]T, T any](s S, f func(T) bool) S {
	if s == nil {
		return nil
	}

	r := make(S, len(s))

	j := 0
	ForEach(s, func(i int, e T) {
		if f(e) {
			return
		}
		r[j] = e
		j++
	})

	return r[:j]
}

func Reverse[S ~[]T, T any](s S) S {
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

func ToMap[S ~[]T, T any, K comparable, V any](s S, f func(T) (K, V)) map[K]V {
	r := make(map[K]V)
	for _, e := range s {
		k, v := f(e)
		r[k] = v
	}
	return r
}

// Shuffle return shuffled slice
func Shuffle[S ~[]T, T any](s S) S {
	if s == nil {
		return nil
	}

	r := make(S, len(s))
	copy(r, s)

	rnd.Shuffle(len(r), func(i, j int) { r[i], r[j] = r[j], r[i] })

	return r
}

// Distinct return distinct slice
func Distinct[S ~[]T, T comparable](s S) S {
	if s == nil {
		return nil
	}

	set := NewSet(s)
	return set.Slice()
}
func Uniq[S ~[]T, T comparable](s S) S { return Distinct(s) }

func Find[S ~[]T, T any](s S, fx func(T) bool) (T, bool) {
	for _, e := range s {
		if fx(e) {
			return e, true
		}
	}

	var r T
	return r, false
}

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
	r = make(map[K]V, len(keys))

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

func UniqBy[S ~[]T, T any, U comparable](s S, f func(T) U) []T {
	r := []T{}
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

type Slice[S ~[]T, T any] []T

func NewSlice[T any](ss ...[]T) Slice[[]T, T] { return Flatten(ss...) }

func (s Slice[S, T]) Clip() Slice[S, T]                         { return Clip(s) }
func (s Slice[S, T]) Clone() Slice[S, T]                        { return Clone(s) }
func (s Slice[S, T]) ContainsFunc(f func(T) bool) bool          { return ContainsFunc(s, f) }
func (s Slice[S, T]) CompactFunc(f func(T, T) bool) Slice[S, T] { return CompactFunc(s, f) }
func (s Slice[S, T]) Delete(i, j int) Slice[S, T]               { return Delete(s, i, j) }
func (s Slice[S, T]) EqualFunc(s1 S, eq func(T, T) bool) bool   { return EqualFunc(s, s1, eq) }
func (s Slice[S, T]) Filter(fx func(T) bool) Slice[S, T]        { return Filter(s, fx) }
func (s Slice[S, T]) Find(fx func(T) bool) (T, bool)            { return Find(s, fx) }
func (s Slice[S, T]) Flatten(ss ...[]T) Slice[S, T]             { return Flatten(ss...) }
func (s Slice[S, T]) ForEach(fx func(int, T))                   { ForEach(s, fx) }
func (s Slice[S, T]) ForEachE(fx func(int, T) error) error      { return ForEachE(s, fx) }
func (s Slice[S, T]) Grow(n int) Slice[S, T]                    { return Grow(s, n) }
func (s Slice[S, T]) IndexFunc(f func(T) bool) int              { return IndexFunc(s, f) }
func (s Slice[S, T]) Insert(i int, v ...T) Slice[S, T]          { return Insert(s, i, v...) }
func (s Slice[S, T]) MaxBy(cmp func(a T, b T) bool) T           { return MaxBy(s, cmp) }
func (s Slice[S, T]) MinBy(cmp func(a T, b T) bool) T           { return MinBy(s, cmp) }
func (s Slice[S, T]) Replace(i, j int, v ...T) Slice[S, T]      { return Replace(s, i, j, v...) }
func (s Slice[S, T]) Sample() T                                 { return Sample(s) }
func (s Slice[S, T]) Samples(count int) Slice[S, T]             { return Samples(s, count) }
func (s Slice[S, T]) Shuffle() Slice[S, T]                      { return Shuffle(s) }
func (s Slice[S, T]) Slice() S                                  { return S(s) }

func (s Slice[S, T]) SortFunc(less func(a, b T) bool) Slice[S, T] {
	r := make(S, len(s))
	copy(r, s)
	SortFunc(r, less)
	return Slice[S, T](r)
}

func (s Slice[S, T]) SortStableFunc(less func(a, b T) bool) S {
	r := make(S, len(s))
	copy(r, s)
	SortStableFunc(s, less)
	return r
}

// SliceC Slice for comparable
type SliceC[S ~[]T, T comparable] []T

func NewSliceC[T comparable](ss ...[]T) SliceC[[]T, T] { return Flatten(ss...) }

func (s SliceC[S, T]) Contains(e T) bool          { return Contains(s, e) }
func (s SliceC[S, T]) Compact() SliceC[S, T]      { return Compact(s) }
func (s SliceC[S, T]) Distinct() SliceC[S, T]     { return SliceC[S, T](Distinct(S(s))) }
func (s SliceC[S, T]) Equal(s1 S) bool            { return Equal(s, s1) }
func (s SliceC[S, T]) Every(sub S) bool           { return Every(S(s), sub) }
func (s SliceC[S, T]) Index(e T) int              { return Index(s, e) }
func (s SliceC[S, T]) Interset(s2 S) SliceC[S, T] { return SliceC[S, T](Interset(S(s), s2)) }
func (s SliceC[S, T]) Slice() S                   { return S(s) }
func (s SliceC[S, T]) Union(ss ...S) SliceC[S, T] {
	return SliceC[S, T](Union(append([]S{S(s)}, ss...)...))
}

// aliases to golang.x/exp/slices

type Ordered = constraints.Ordered

func Clip[S ~[]T, T any](s S) S                           { return slices.Clip(s) }
func Clone[S ~[]T, T any](s S) S                          { return slices.Clone(s) }
func Compact[S ~[]T, T comparable](s S) S                 { return slices.Compact(s) }
func CompactFunc[S ~[]T, T any](s S, f func(T, T) bool) S { return slices.CompactFunc(s, f) }
func Compare[T Ordered](s1, s2 []T) int                   { return slices.Compare(s1, s2) }
func CompareFunc[T1, T2 any](s1 []T1, s2 []T2, cmp func(T1, T2) int) int {
	return slices.CompareFunc(s1, s2, cmp)
}
func Contains[S ~[]T, T comparable](s S, e T) bool         { return slices.Contains(s, e) }
func ContainsFunc[S ~[]T, T any](s S, f func(T) bool) bool { return slices.ContainsFunc(s, f) }
func Delete[S ~[]T, T any](s S, i, j int) S                { return slices.Delete(s, i, j) }
func Equal[T comparable](s1, s2 []T) bool                  { return slices.Equal(s1, s2) }
func EqualFunc[T1, T2 any](s1 []T1, s2 []T2, eq func(T1, T2) bool) bool {
	return slices.EqualFunc(s1, s2, eq)
}
func Grow[S ~[]T, T any](s S, n int) S                 { return slices.Grow(s, n) }
func Index[S ~[]T, T comparable](s S, e T) int         { return slices.Index(s, e) }
func IndexFunc[S ~[]T, T any](s S, f func(T) bool) int { return slices.IndexFunc(s, f) }
func Insert[S ~[]T, T any](s S, i int, v ...T) S       { return slices.Insert(s, i, v...) }
func Replace[S ~[]T, T any](s S, i, j int, v ...T) S   { return slices.Replace(s, i, j, v...) }

// sort
func BinarySearch[T Ordered](x []T, target T) (int, bool) { return slices.BinarySearch(x, target) }
func BinarySearchFunc[T any](x []T, target T, cmp func(T, T) int) (int, bool) {
	return slices.BinarySearchFunc(x, target, cmp)
}
func IsSorted[T Ordered](s []T) bool                         { return slices.IsSorted(s) }
func IsSortedFunc[T any](s []T, less func(a, b T) bool) bool { return slices.IsSortedFunc(s, less) }

func Sort[S ~[]T, T Ordered](s S) S {
	r := make(S, len(s))
	copy(r, s)
	slices.Sort(r)
	return r
}

func SortFunc[S ~[]T, T any](s S, less func(a, b T) bool) S {
	r := make(S, len(s))
	copy(r, s)
	slices.SortFunc(r, less)
	return r
}

func SortStableFunc[S ~[]T, T any](s S, less func(a, b T) bool) S {
	r := make(S, len(s))
	copy(r, s)
	slices.SortStableFunc(r, less)
	return r
}
