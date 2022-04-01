package fx

type Set[T comparable] struct {
	keys   map[T]struct{}
	values []T
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{}
}

func (s *Set[T]) Slice() []T { return s.values }
func (s *Set[T]) Len() int   { return len(s.values) }

func (s *Set[T]) Append(elements ...T) {
	for _, e := range elements {
		if _, ok := s.keys[e]; ok {
			continue
		}

		s.values = append(s.values, e)
	}
}

func (s *Set[T]) ForEach(fx func(int, T)) { ForEach(s.values, fx) }
