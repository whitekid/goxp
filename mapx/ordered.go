package mapx

import (
	"iter"
	"slices"
)

// Ordered
type Ordered[Map ~map[K]V, K comparable, V any] struct {
	store map[K]V
	keys  []K
}

func NewOrdered[K comparable, V any]() *Ordered[map[K]V, K, V] {
	return &Ordered[map[K]V, K, V]{
		store: map[K]V{},
		keys:  []K{},
	}
}

func (o *Ordered[Map, K, V]) Get(key K) (V, bool) {
	value, exists := o.store[key]
	return value, exists
}

func (o *Ordered[Map, K, V]) Set(key K, value V) {
	if _, exists := o.store[key]; !exists {
		o.keys = append(o.keys, key)
	}

	o.store[key] = value
}

func (o *Ordered[Map, K, V]) Delete(key K) {
	delete(o.store, key)

	idx := slices.Index(o.keys, key)
	if idx != -1 {
		o.keys = slices.Delete(o.keys, idx, idx+1)
	}
}

func (o *Ordered[Map, K, V]) Len() int          { return len(o.keys) }
func (o *Ordered[Map, K, V]) Keys() iter.Seq[K] { return slices.Values(o.keys) }
func (o *Ordered[Map, K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, k := range o.keys {
			if !yield(o.store[k]) {
				return
			}
		}
	}
}

func (o *Ordered[Map, K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, k := range o.keys {
			if !yield(k, o.store[k]) {
				return
			}
		}
	}
}

func (o *Ordered[Map, K, V]) Each(each func(int, K, V) bool) {
	for i, k := range o.keys {
		if !each(i, k, o.store[k]) {
			break
		}
	}
}
