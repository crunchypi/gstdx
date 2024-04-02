package iox

import (
	"context"
	"io"
)

// TODO: Consider Trade.
type Swapper[T, U any] interface {
	Swap(ctx context.Context, v T) (r U, err error)
}

type SwapperImpl[T, U any] struct {
	Impl func(context.Context, T) (U, error)
}

func (impl SwapperImpl[T, U]) Swap(ctx context.Context, v T) (r U, err error) {
	if impl.Impl == nil {
		err = io.EOF
		return
	}

	return impl.Impl(ctx, v)
}

func SwapFilterIFn[T, U any](s Swapper[T, U]) func(f func(T) bool) Swapper[T, U] {
	return func(f func(T) bool) Swapper[T, U] {
		return SwapperImpl[T, U]{
			Impl: func(ctx context.Context, v T) (r U, err error) {
				// TODO: Filtering here doesn't make sense.

				r, err = s.Swap(ctx, v)

				return
			},
		}
	}
}

func SwapI1MapFn[T, U, V any](s Swapper[V, U]) func(f func(T) V) Swapper[T, U] {
	return func(f func(T) V) Swapper[T, U] {
		return SwapperImpl[T, U]{
			Impl: func(ctx context.Context, v T) (r U, err error) {
				r, err = s.Swap(ctx, f(v))
				return
			},
		}
	}
}

func SwapI2MapFn[T, V any](s Swapper[T, any]) func(f func(V) T) Swapper[V, any] {
	return func(f func(V) T) Swapper[V, any] {
		return SwapperImpl[V, any]{
			Impl: func(ctx context.Context, v V) (r any, err error) {
				r, err = s.Swap(ctx, f(v))
				return
			},
		}
	}
}
