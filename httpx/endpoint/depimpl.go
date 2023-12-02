package endpoint

import (
	"context"
	"errors"
	"net/http"
)

type ReaderImpl[T any] struct {
	Impl func(
		ctx context.Context,
		namespace string,
	) (
		v T,
		ssc int,
		err error,
	)
}

func (impl ReaderImpl[T]) Read(
	ctx context.Context,
	namespace string,
) (
	v T,
	ssc int,
	err error,
) {
	if impl.Impl == nil {
		ssc = http.StatusNotImplemented
		err = errors.New("endpoint.ReaderImpl: used without impl")
		return
	}

	return impl.Impl(ctx, namespace)
}

type WriterImpl[T any] struct {
	Impl func(
		ctx context.Context,
		namespace string,
		v T,
	) (
		ssc int,
		err error,
	)
}

func (impl WriterImpl[T]) Write(
	ctx context.Context,
	namespace string,
	v T,
) (
	ssc int,
	err error,
) {
	if impl.Impl == nil {
		ssc = http.StatusNotImplemented
		err = errors.New("endpoint.WriterImpl: used without impl")
		return
	}

	return impl.Impl(ctx, namespace, v)
}

type ReadWriterImpl[I, O any] struct {
	Impl func(
		ctx context.Context,
		namespace string,
		in I,
	) (
		out O,
		ssc int,
		err error,
	)
}

func (impl ReadWriterImpl[I, O]) ReadWrite(
	ctx context.Context,
	namespace string,
	in I,
) (
	out O,
	ssc int,
	err error,
) {
	if impl.Impl == nil {
		ssc = http.StatusNotImplemented
		err = errors.New("endpoint.ReadWriterImpl: used without impl")
		return
	}

	return impl.Impl(ctx, namespace, in)
}
