package types

type Map[K comparable, V any] map[K]V

func NewMap[K comparable, V any](x map[K]V) Map[K, V] { return Map[K, V](x) }

func (m Map[K, V]) GetMap() map[K]V { return map[K]V(m) }

func (m Map[K, V]) Each(fx func(k K, v V)) {
	for k, v := range m {
		fx(k, v)
	}
}
