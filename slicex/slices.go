package slicex

import (
	"cmp"
	"iter"
	"slices"
)

type Slice[S ~[]E, E any] []E

func (s Slice[S, E]) Slice() []E { return []E(s) }

func Of[E any](e ...E) Slice[[]E, E] { return e }

func (s Slice[S, E]) All() iter.Seq2[int, E]                { return slices.All(s) }
func (s Slice[S, E]) AppendSeq(seq iter.Seq[E]) Slice[S, E] { return slices.AppendSeq(s, seq) }
func (s Slice[S, E]) Backward() iter.Seq2[int, E]           { return slices.Backward(s) }
func (s Slice[S, E]) Chunk(n int) iter.Seq[[]E]             { return slices.Chunk(s.Slice(), n) }
func (s Slice[S, E]) Clip() Slice[S, E]                     { return slices.Clip(s) }
func (s Slice[S, E]) Clone() Slice[S, E]                    { return slices.Clone(s) }
func (s Slice[S, E]) Delete(i, j int) Slice[S, E]           { return slices.Delete(s, i, j) }
func (s Slice[S, E]) Grow(n int) Slice[S, E]                { return slices.Grow(s, n) }
func (s Slice[S, E]) Insert(i int, v ...E) Slice[S, E]      { return slices.Insert(s, i, v...) }
func (s Slice[S, E]) Replace(i, j int, v ...E) Slice[S, E]  { return slices.Replace(s, i, j, v...) }
func (s Slice[S, E]) Reverse() Slice[S, E]                  { slices.Reverse(s); return s }
func (s Slice[S, E]) Values() iter.Seq[E]                   { return slices.Values(s) }

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

func Times[T any](count int, f func(int) T) Slice[[]T, T] {
	r := make([]T, count)

	for i := range r {
		r[i] = f(i)
	}

	return r
}

func To[S ~[]E, E any](s []any) S   { return Map(s, func(x any) E { return x.(E) }) }
func ToPtr[S ~[]E, E any](s S) []*E { return Map(s, func(x E) *E { return &x }) }
