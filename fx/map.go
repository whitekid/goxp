package fx

type Dict[K comparable, V any] map[K]V

func FromMap[K comparable, V any](m map[K]V) Dict[K, V] { return Dict[K, V](m) }

func (d Dict[K, V]) Keys() []K   { return Keys(d) }
func (d Dict[K, V]) Values() []V { return Values(d) }

// FilterMap ...
func FilterMap[K comparable, V any](m map[K]V, fx func(K, V) bool) map[K]V {
	r := make(map[K]V)

	for k, v := range m {
		if fx(k, v) {
			r[k] = v
		}
	}

	return r
}
