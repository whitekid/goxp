package fx

type Dict[K comparable, V any] map[K]V

func FromMap[K comparable, V any](m map[K]V) Dict[K, V] { return Dict[K, V](m) }

func (d Dict[K, V]) Keys() []K                      { return Keys(d) }
func (d Dict[K, V]) Values() []V                    { return Values(d) }
func (d Dict[K, V]) MapValues(fx func(V) V) map[K]V { return MapValues(d, fx) }
