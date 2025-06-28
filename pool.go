package goxp

import "sync"

// Pool sync.Pool with generics
type Pool[T any] struct {
	sync.Pool
}

func NewPool[T any](New func() T) *Pool[T] {
	pool := &Pool[T]{
		Pool: sync.Pool{
			New: func() any { return New() },
		},
	}

	return pool
}

func (p *Pool[T]) Put(x T) { p.Pool.Put(x) }
func (p *Pool[T]) Get() T {
	v := p.Pool.Get()
	if v == nil {
		var zero T
		return zero
	}
	result, ok := v.(T)
	if !ok {
		var zero T
		return zero
	}
	return result
}
