package sets

import "slices"

type Set[S ~[]E, E comparable] struct {
	keys   map[E]struct{}
	values S
}

func New[E comparable](lists ...[]E) *Set[[]E, E] {
	s := &Set[[]E, E]{
		keys: make(map[E]struct{}),
	}

	for i := range lists {
		for j := range lists[i] {
			s.Set(lists[i][j])
		}
	}

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
		s.values = append(s.values[0:i], s.values[i+1:len(s.values)]...)
	}
}

func (s *Set[S, E]) Contains(e E) (ok bool) {
	_, ok = s.keys[e]
	return
}

func (s *Set[S, E]) Each(fx func(int, E)) {
	for i, v := range s.values {
		fx(i, v)
	}
}
