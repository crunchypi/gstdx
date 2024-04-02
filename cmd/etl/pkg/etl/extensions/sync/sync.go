package sync

import (
	"context"
	"sync"

	"github.com/crunchypi/gstdx/iox"
)

func NewReader[T any](r iox.Reader[T]) iox.Reader[T] {
	if r == nil {
		return iox.ReaderImpl[T]{}
	}

	mx := sync.RWMutex{}
	return iox.ReaderImpl[T]{
		Impl: func(ctx context.Context) (v T, err error) {
			mx.RLock()
			defer mx.RUnlock()

			return r.Read(ctx)
		},
	}
}

func NewWriter[T any](r iox.Writer[T]) iox.Writer[T] {
	if r == nil {
		return iox.WriterImpl[T]{}
	}

	mx := sync.Mutex{}
	return iox.WriterImpl[T]{
		Impl: func(ctx context.Context, v T) error {
			mx.Lock()
			defer mx.Unlock()

			return r.Write(ctx, v)
		},
	}
}
