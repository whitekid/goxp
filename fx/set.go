package fx

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
		i := Index(s.values, e)
		s.values = append(s.values[0:i], s.values[i+1:len(s.values)]...)
	}
}

func (s *Set[T]) Has(e T) bool {
	_, ok := s.keys[e]
	return ok
}

func (s *Set[T]) ForEach(fx func(int, T)) { ForEach(s.values, fx) }
