package types

import "github.com/whitekid/goxp/fx"

type Map[K comparable, V any] map[K]V

func NewMap[K comparable, V any](x map[K]V) Map[K, V] { return Map[K, V](x) }

func (m Map[K, V]) Keys() []K              { return fx.Keys(m) }
func (m Map[K, V]) Values() []V            { return fx.Values(m) }
func (m Map[K, V]) Each(fn func(k K, v V)) { fx.ForEachMap(m, fn) }
