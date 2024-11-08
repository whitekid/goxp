package mapx

import (
	"iter"
	"slices"
	"sync"

	"github.com/whitekid/goxp"
)

type SyncMap[K comparable, V any] struct {
	m sync.Map
}

func NewSyncMap[K comparable, V any]() *SyncMap[K, V] {
	return &SyncMap[K, V]{
		m: sync.Map{},
	}
}

func (m *SyncMap[K, V]) Store(k K, v V) { m.m.Store(k, v) }

func (m *SyncMap[K, V]) Get(k K) (V, bool) {
	v, ok := m.m.Load(k)
	return v.(V), ok
}

func (m *SyncMap[K, V]) Delete(k K) { m.m.Delete(k) }

func (m *SyncMap[K, V]) All() iter.Seq2[K, V] {
	items := []*goxp.Tuple2[K, V]{}

	m.m.Range(func(key, value any) bool {
		items = append(items, goxp.T2(key.(K), value.(V)))
		return true
	})

	return func(yield func(K, V) bool) {
		for _, e := range items {
			if !yield(e.V1, e.V2) {
				return
			}
		}
	}
}

func (m *SyncMap[K, V]) Keys() iter.Seq[K] {
	items := []K{}

	m.m.Range(func(key, value any) bool {
		items = append(items, key.(K))
		return true
	})

	return slices.Values(items)
}

func (m *SyncMap[K, V]) Values() iter.Seq[V] {
	items := []V{}

	m.m.Range(func(key, value any) bool {
		items = append(items, value.(V))
		return true
	})

	return slices.Values(items)
}
