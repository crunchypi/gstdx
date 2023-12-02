package syncx

import "sync"

type Protected[T any] struct {
	mx sync.RWMutex
	v  T
}

func NewProtected[T any](v T) *Protected[T] {
	return &Protected[T]{v: v}
}

func (p *Protected[T]) Put(v T) {
	p.mx.Lock()
	defer p.mx.Unlock()

	p.v = v
}

func (p *Protected[T]) Get(f func(v T)) {
	if f == nil {
		return
	}

	p.mx.RLock()
	defer p.mx.RUnlock()

	f(p.v)
}

func (p *Protected[T]) Mod(f func(v T) T) {
	if f == nil {
		return
	}

	p.mx.Lock()
	defer p.mx.Unlock()

	p.v = f(p.v)
}
