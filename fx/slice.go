package fx

type Slice[T any] []T

func NewSlice[T any](collection []T) Slice[T] { return Slice[T](collection) }

func (s Slice[T]) ForEach(fx func(int, T))         { ForEach(s, fx) }
func (s Slice[T]) Map(fx func(T) T) Slice[T]       { return Map(s, fx) }
func (s Slice[T]) Filter(fx func(T) bool) Slice[T] { return Filter(s, fx) }
func (s Slice[T]) Shuffle() []T                    { return Shuffle(s) }
func (s Slice[T]) Sample() T                       { return Sample(s) }
func (s Slice[T]) Samples(count int) []T           { return Samples(s, count) }
