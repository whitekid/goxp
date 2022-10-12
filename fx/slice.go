package fx

type Slice[T any] []T

func S[T any](s []T) Slice[T] { return Slice[T](s) }

func (s Slice[T]) Slice() []T           { return s }
func (s Slice[T]) Each(fx func(int, T)) { ForEach(s, fx) }
func (s Slice[T]) EachE(fx func(int, T) error) error {
	for i, e := range s {
		if err := fx(i, e); err != nil {
			return err
		}
	}
	return nil
}
func (s Slice[T]) Map(fx func(T) T) Slice[T]       { return Map(s, fx) }
func (s Slice[T]) Filter(fx func(T) bool) Slice[T] { return Filter(s, fx) }
func (s Slice[T]) Reduce(fx func(r T, e T) T) T    { return Reduce(s, fx) }

func (s Slice[T]) Shuffle() []T {
	s1 := make([]T, len(s))
	copy(s1, s)
	Shuffle(s1)
	return s1
}

func (s Slice[T]) Sample() T                  { return Sample(s) }
func (s Slice[T]) Samples(count int) Slice[T] { return Samples(s, count) }
