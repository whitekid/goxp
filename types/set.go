package types

import "golang.org/x/exp/slices"

type Set[T comparable] struct {
	keys   map[T]struct{}
	values []T
}

func NewSet[T comparable](lists ...[]T) *Set[T] {
	s := &Set[T]{
		keys: make(map[T]struct{}),
	}

	for i := range lists {
		for j := range lists[i] {
			s.Append(lists[i][j])
		}
	}

	return s
}

func (s *Set[T]) Slice() (r []T) {
	r = make([]T, len(s.values))
	copy(r, s.values)
	return r
}

func (s *Set[T]) Len() int {
	return len(s.values)
}

func (s *Set[T]) Append(elements ...T) {
	for _, e := range elements {
		if _, ok := s.keys[e]; ok {
			continue
		}

		s.keys[e] = struct{}{}
		s.values = append(s.values, e)
	}
}

func (s *Set[T]) Remove(elements ...T) {
	for _, e := range elements {
		if _, ok := s.keys[e]; !ok {
			continue
		}

		delete(s.keys, e)
		i := slices.Index(s.values, e)
		s.values = append(s.values[0:i], s.values[i+1:len(s.values)]...)
	}
}

func (s *Set[T]) Has(e T) (ok bool) {
	_, ok = s.keys[e]
	return
}

func (s *Set[T]) Each(fx func(int, T)) {
	for i, v := range s.values {
		fx(i, v)
	}
}
