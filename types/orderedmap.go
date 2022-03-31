package types

type OrderedMap[K comparable, V any] struct {
	store map[K]V
	keys  []K
}

func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		store: map[K]V{},
		keys:  []K{},
	}
}

func (o *OrderedMap[K, V]) Get(key K) (V, bool) {
	value, exists := o.store[key]
	return value, exists
}

func (o *OrderedMap[K, V]) Set(key K, value V) {
	if _, exists := o.store[key]; !exists {
		o.keys = append(o.keys, key)
	}

	o.store[key] = value
}

func (o *OrderedMap[K, V]) Delete(key K) {
	delete(o.store, key)

	idx := -1

	for i, val := range o.keys {
		if val == key {
			idx = i
			break
		}
	}

	if idx != -1 {
		o.keys = append(o.keys[:idx], o.keys[idx+1:]...)
	}
}

func (o *OrderedMap[K, V]) Len() int  { return len(o.keys) }
func (o *OrderedMap[K, V]) Keys() []K { return o.keys }

func (o *OrderedMap[K, V]) ForEach(each func(int, K, V) bool) {
	for i, k := range o.keys {
		if each(i, k, o.store[k]) == false {
			break
		}
	}
}
