package iox

import (
	"context"
	"io"
)

// -----------------------------------------------------------------------------
// Original io.ReadWriter iface impl.
// -----------------------------------------------------------------------------

type readWriteCloserImpl struct {
	ImplC func() error
	ImplR func([]byte) (int, error)
	ImplW func([]byte) (int, error)
}

func (impl readWriteCloserImpl) Close() (err error) {
	if impl.ImplC == nil {
		return
	}

	return impl.ImplC()
}

func (impl readWriteCloserImpl) Read(p []byte) (n int, err error) {
	if impl.ImplR == nil {
		err = io.EOF
		return
	}

	return impl.ImplR(p)
}

func (impl readWriteCloserImpl) Write(p []byte) (n int, err error) {
	if impl.ImplW == nil {
		err = io.ErrClosedPipe
		return
	}

	return impl.ImplW(p)
}

// -----------------------------------------------------------------------------
// New ReadWriter iface + impl.
// -----------------------------------------------------------------------------

type ReadWriter[T, U any] interface {
	Reader[T]
	Writer[U]
}

type ReadWriterImpl[T, U any] struct {
	ImplR func(context.Context) (T, error)
	ImplW func(context.Context, U) error
}

func (impl ReadWriterImpl[T, U]) Read(ctx context.Context) (r T, err error) {
	if impl.ImplR == nil {
		err = io.EOF
		return
	}

	return impl.ImplR(ctx)
}

func (impl ReadWriterImpl[T, U]) Write(ctx context.Context, v U) (err error) {
	if impl.ImplW == nil {
		err = io.ErrClosedPipe
		return
	}

	return impl.ImplW(ctx, v)
}

// -----------------------------------------------------------------------------
// New ReadWriteCloser iface + impl.
// -----------------------------------------------------------------------------

type ReadWriteCloser[T, U any] interface {
	io.Closer
	Reader[T]
	Writer[U]
}

type ReadWriteCloserImpl[T, U any] struct {
	ImplC func() error
	ImplR func(context.Context) (T, error)
	ImplW func(context.Context, U) error
}

func (impl ReadWriteCloserImpl[T, U]) Close() (err error) {
	if impl.ImplC == nil {
		return
	}

	return impl.ImplC()
}

func (impl ReadWriteCloserImpl[T, U]) Read(ctx context.Context) (r T, err error) {
	if impl.ImplR == nil {
		err = io.EOF
		return
	}

	return impl.ImplR(ctx)
}

func (impl ReadWriteCloserImpl[T, U]) Write(ctx context.Context, v U) (err error) {
	if impl.ImplW == nil {
		err = io.ErrClosedPipe
		return
	}

	return impl.ImplW(ctx, v)
}

// -----------------------------------------------------------------------------
// Factory funcs.
// -----------------------------------------------------------------------------

func NewV2VReadWriter[T any](vs ...T) ReadWriter[T, T] {
	buf := make([]T, len(vs))
	copy(buf, vs)
	return ReadWriteCloserImpl[T, T]{
		ImplR: func(ctx context.Context) (v T, err error) {
			if len(buf) == 0 {
				return v, io.EOF
			}

			v = buf[0]
			buf = buf[1:]

			return
		},
		ImplW: func(ctx context.Context, v T) (err error) {
			buf = append(buf, v)
			return
		},
	}
}

func NewT2UReadWriterFn[T, U any](
	rw io.ReadWriter,
) (
	_ func(d func(io.Reader) Decoder, e func(io.Writer) Encoder) (_ ReadWriter[T, U]),
) {
	return func(d func(io.Reader) Decoder, e func(io.Writer) Encoder) (_ ReadWriter[T, U]) {
		return ReadWriterImpl[T, U]{
			ImplR: NewB2VReaderFn[T](rw)(d).Read,
			ImplW: NewV2BWriterFn[U](rw)(e).Write,
		}
	}
}

func NewT2UReadWriteCloserFn[T, U any](
	rw io.ReadWriteCloser,
) (
	_ func(d func(io.Reader) Decoder, e func(io.Writer) Encoder) (_ ReadWriter[T, U]),
) {
	return func(d func(io.Reader) Decoder, e func(io.Writer) Encoder) (_ ReadWriter[T, U]) {
		return ReadWriteCloserImpl[T, U]{
			ImplC: rw.Close,
			ImplR: NewB2VReaderFn[T](rw)(d).Read,
			ImplW: NewV2BWriterFn[U](rw)(e).Write,
		}
	}
}
