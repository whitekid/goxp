package goxp

import "sync"

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
func (p *Pool[T]) Get() T  { return p.Pool.Get().(T) }
