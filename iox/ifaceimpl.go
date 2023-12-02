package iox

import "context"

type ReaderImpl[T any] struct {
	Impl func(context.Context) (T, error)
}

func (impl ReaderImpl[T]) Read(
	ctx context.Context,
) (
	r T,
	err error,
) {
	return impl.Impl(ctx)
}

type WriterImpl[T any] struct {
	Impl func(context.Context, T) error
}

func (impl WriterImpl[T]) Write(
	ctx context.Context,
	v T,
) (
	err error,
) {
	return impl.Write(ctx, v)
}

type ReadWriterImpl[T, U any] struct {
	Impl func(context.Context, T) (U, error)
}

func (impl ReadWriterImpl[T, U]) ReadWrite(
	ctx context.Context,
	v T,
) (
	r U,
	err error,
) {
	return impl.Impl(ctx, v)
}
