package sets

import (
	"iter"
	"slices"
)

type Set[S ~[]E, E comparable] struct {
	keys   map[E]struct{}
	values S
}

func New[S ~[]E, E comparable](e ...E) *Set[S, E] {
	s := &Set[S, E]{
		keys: make(map[E]struct{}),
	}

	s.Set(e...)

	return s
}

func (s *Set[S, E]) Slice() (r S) {
	r = make(S, len(s.values))
	copy(r, s.values)
	return r
}

func (s *Set[S, E]) Len() int {
	return len(s.values)
}

func (s *Set[S, E]) Set(elements ...E) {
	for _, e := range elements {
		if _, ok := s.keys[e]; ok {
			continue
		}

		s.keys[e] = struct{}{}
		s.values = append(s.values, e)
	}
}

func (s *Set[S, E]) Remove(elements ...E) {
	for _, e := range elements {
		if _, ok := s.keys[e]; !ok {
			continue
		}

		delete(s.keys, e)
		i := slices.Index(s.values, e)
		s.values = slices.Delete(s.values, i, i+1)
	}
}

func (s *Set[S, E]) Contains(e E) bool {
	_, ok := s.keys[e]
	return ok
}

func (s *Set[S, E]) Each(fx func(int, E)) {
	for i, v := range s.values {
		fx(i, v)
	}
}

func (s *Set[S, E]) All() iter.Seq2[int, E] { return slices.All(s.values) }
func (s *Set[S, E]) Values() iter.Seq[E]    { return slices.Values(s.values) }
