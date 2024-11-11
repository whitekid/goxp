package slicex

import (
	"cmp"
	"crypto/rand"
	"iter"
	"math/big"
	mrand "math/rand/v2"
	"slices"

	"github.com/whitekid/goxp/fx/gen"
	"github.com/whitekid/goxp/sets"
)

type Slice[S ~[]E, E any] []E

func (s Slice[S, E]) Slice() []E { return S(s) }

func Of[S ~[]E, E any](e ...E) Slice[S, E] { return e }
func Gen[S ~[]E, E any](next gen.Generator[E]) Slice[S, E] {
	s := []E{}
	for v, ok := next(); ok; v, ok = next() {
		s = append(s, v)
	}
	return s
}

func (s Slice[S, E]) All() iter.Seq2[int, E]                { return slices.All(s) }
func (s Slice[S, E]) AppendSeq(seq iter.Seq[E]) Slice[S, E] { return slices.AppendSeq(s, seq) }
func (s Slice[S, E]) Backward() iter.Seq2[int, E]           { return slices.Backward(s) }
func (s Slice[S, E]) Concat(ss ...S) S                      { return slices.Concat(append([]S{S(s)}, ss...)...) }
func (s Slice[S, E]) Chunk(n int) iter.Seq[S]               { return slices.Chunk[S, E](s.Slice(), n) }
func (s Slice[S, E]) Clip() Slice[S, E]                     { return slices.Clip(s) }
func (s Slice[S, E]) Clone() Slice[S, E]                    { return slices.Clone(s) }
func (s Slice[S, E]) Delete(i, j int) Slice[S, E]           { return slices.Delete(s, i, j) }
func (s Slice[S, E]) Grow(n int) Slice[S, E]                { return slices.Grow(s, n) }
func (s Slice[S, E]) Insert(i int, v ...E) Slice[S, E]      { return slices.Insert(s, i, v...) }
func (s Slice[S, E]) Replace(i, j int, v ...E) Slice[S, E]  { return slices.Replace(s, i, j, v...) }
func (s Slice[S, E]) Reverse() Slice[S, E] {
	s1 := slices.Clone(s)
	slices.Reverse(s1)
	return s1
}
func (s Slice[S, E]) Values() iter.Seq[E] { return slices.Values(s) }

func Repeat[S ~[]E, E any](x S, count int) Slice[S, E] {
	return Slice[S, E](slices.Repeat(x, count))
}

func Map[S ~[]T1, T1 any, T2 any](s S, fx func(T1) T2) []T2 {
	if s == nil {
		return nil
	}

	r := make([]T2, len(s))
	for i := 0; i < len(s); i++ {
		r[i] = fx(s[i])
	}

	return r
}

func Max[E cmp.Ordered](s ...E) E { return slices.Max(s) }
func Min[E cmp.Ordered](s ...E) E { return slices.Min(s) }

func Sample[S ~[]E, E any](s S) E {
	i, err := rand.Int(rand.Reader, big.NewInt(int64(len(s))))
	if err != nil {
		panic(err)
	}

	return s[i.Int64()]
}

func Times[T any](count int, f func(int) T) Slice[[]T, T] {
	r := make([]T, count)

	for i := range r {
		r[i] = f(i)
	}

	return r
}

func To[E any](s []any) []E   { return Map(s, func(x any) E { return x.(E) }) }
func ToPtr[E any](s []E) []*E { return Map(s, func(x E) *E { return &x }) }

func (s Slice[S, E]) Each(fn func(int, E)) { Each(s, fn) }
func Each[S ~[]E, E any](s S, fx func(int, E)) {
	for i, e := range s {
		fx(i, e)
	}
}

func (s Slice[S, E]) Filter(fn func(E) bool) Slice[S, E] { return Filter(s, fn) }
func Filter[S ~[]E, E any](s S, fx func(E) bool) S {
	if s == nil {
		return nil
	}

	r := make(S, 0, len(s))

	Each(s, func(i int, e E) {
		if !fx(e) {
			return
		}
		r = append(r, e)
	})

	return slices.Clip(r)
}

func (s Slice[S, E]) Reduce(fn func(E, E) E) E { return Reduce(s, fn) }
func Reduce[S ~[]E, E any](s S, fx func(E, E) E) (r E) {
	for _, e := range s {
		r = fx(r, e)
	}

	return r
}

func Shuffle[S ~[]E, E any](s S) S {
	if s == nil {
		return nil
	}

	r := slices.Clone(s)
	mrand.Shuffle(len(r), func(i, j int) { r[i], r[j] = r[j], r[i] })

	return r
}

func Uniq[S ~[]E, E comparable](s S) S {
	seen := map[E]struct{}{}
	r := make([]E, 0, len(s))

	for _, e := range s {
		if _, ok := seen[e]; ok {
			continue
		}

		r = append(r, e)
		seen[e] = struct{}{}
	}

	return slices.Clip(r)
}

func (s Slice[S, E]) SortedFunc(cmp func(E, E) int) iter.Seq[E] { return SortedFunc(s.Values(), cmp) }
func SortedFunc[E any](seq iter.Seq[E], cmp func(E, E) int) iter.Seq[E] {
	return slices.Values(slices.SortedFunc(seq, cmp))
}

func Intersect[S ~[]E, E comparable](s1, s2 S) S {
	s := sets.New(s2)

	return Filter(s1, func(e E) bool {
		if s.Contains(e) {
			return false
		}

		s.Set(e)
		return true
	})
}

func Flatten[T any](s [][]T) []T {
	r := []T{}

	for _, e := range s {
		r = append(r, e...)
	}

	return r
}

func GroupBy[S ~[]E, E any, K comparable](s S, key func(E) K) map[K]S {
	r := make(map[K]S, 0)

	for _, e := range s {
		k := key(e)
		if _, ok := r[k]; !ok {
			r[k] = []E{}
		}

		r[k] = append(r[k], e)
	}

	return r
}
