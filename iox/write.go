package iox

import (
	"bytes"
	"context"
	"encoding/gob"
	"io"
)

// -----------------------------------------------------------------------------
// New Writer iface + impl.
// -----------------------------------------------------------------------------

type Writer[T any] interface {
	Write(context.Context, T) error
}

// WriterImpl supports a functional implementation of Writer.
type WriterImpl[T any] struct {
	Impl func(context.Context, T) error
}

// Write defers to WriterImpl.Impl.
func (impl WriterImpl[T]) Write(
	ctx context.Context,
	v T,
) (
	err error,
) {
	if impl.Impl == nil {
		err = io.ErrClosedPipe
		return
	}

	return impl.Impl(ctx, v)
}

// -----------------------------------------------------------------------------
// New WriteCloser iface + impl.
// -----------------------------------------------------------------------------

type WriteCloser[T any] interface {
	io.Closer
	Writer[T]
}

// WriteCloserImpl supports a functional implementation of WriteCloser.
type WriteCloserImpl[T any] struct {
	ImplC func() error
	ImplW func(context.Context, T) error
}

// Write defers to WriteCloserImpl.CImpl.
func (impl WriteCloserImpl[T]) Close() error {
	if impl.ImplC == nil {
		return nil
	}

	return impl.ImplC()
}

// Write defers to WriteCloserImpl.WImpl.
func (impl WriteCloserImpl[T]) Write(
	ctx context.Context,
	v T,
) (
	err error,
) {
	if impl.ImplW == nil {
		err = io.ErrClosedPipe
		return
	}

	return impl.ImplW(ctx, v)
}

// -----------------------------------------------------------------------------
// Factory funcs.
// -----------------------------------------------------------------------------

func NewB2VWriterFn[T any](w Writer[T]) func(d Decoder) io.Writer {
	return func(d Decoder) io.Writer {
		buf := bytes.NewBuffer(nil)
		var dec Decoder = gob.NewDecoder(buf)

		if d != nil {
			dec = d
		}

		if w == nil {
			return readWriteCloserImpl{}
		}

		return readWriteCloserImpl{
			ImplW: func(p []byte) (n int, err error) {
				n, err = buf.Write(p)
				if err != nil {
					return
				}

				var v T
				err = dec.Decode(&v)

				if err != nil {
					return
				}

				err = w.Write(nil, v)
				if err != nil {
					return
				}

				return
			},
		}
	}
}

func NewB2VWriteCloserFn[T any](w WriteCloser[T]) func(d Decoder) io.WriteCloser {
	return func(d Decoder) io.WriteCloser {
		return readWriteCloserImpl{
			ImplC: w.Close,
			ImplW: NewB2VWriterFn[T](w)(d).Write,
		}
	}
}

func NewV2BWriterFn[T any](w io.Writer) func(f func(io.Writer) Encoder) Writer[T] {
	return func(f func(io.Writer) Encoder) Writer[T] {
		buf := bytes.NewBuffer(nil)
		enc := func(w io.Writer) Encoder { return gob.NewEncoder(w) }(buf)

		if f != nil {
			if _e := f(buf); _e != nil {
				enc = _e
			}
		}

		// TODO check nils.

		return WriterImpl[T]{
			Impl: func(ctx context.Context, v T) error {
				err := enc.Encode(v)
				if err != nil {
					return err
				}

				_, err = buf.WriteTo(w)
				return err
			},
		}
	}
}

func NewV2BWriteCloserFn[T any](w io.WriteCloser) func(f func(io.Writer) Encoder) WriteCloser[T] {
	return func(f func(io.Writer) Encoder) WriteCloser[T] {
		return WriteCloserImpl[T]{
			ImplC: w.Close,
			ImplW: NewV2BWriterFn[T](w)(f).Write,
		}
	}
}

// -----------------------------------------------------------------------------
// Batching.
// -----------------------------------------------------------------------------

func NewBatchedVWriter[T any](w Writer[T], size int) Writer[T] {
	buf := make([]T, 0, size)
	return WriterImpl[T]{
		Impl: func(ctx context.Context, v T) (err error) {
			buf = append(buf, v)

			if len(buf) >= size {
				for i := 0; i < size; i++ {
					_v := buf[0]
					buf = buf[1:]
					err = w.Write(ctx, _v)
					if err != nil {
						return
					}
				}
			}

			return err
		},
	}
}

func NewBatchedSWriter[T any](w Writer[[]T], size int) Writer[T] {
	buf := make([]T, 0, size)
	return WriterImpl[T]{
		Impl: func(ctx context.Context, v T) (err error) {
			buf = append(buf, v)

			if len(buf) >= size {
				err = w.Write(ctx, buf)
				buf = make([]T, 0, size)
			}

			return err
		},
	}
}

// -----------------------------------------------------------------------------
// Functional.
// -----------------------------------------------------------------------------

func WriteFilterFn[T any](w Writer[T]) func(filter func(T) bool) Writer[T] {
	return func(filter func(v T) bool) Writer[T] {
		return WriterImpl[T]{
			Impl: func(ctx context.Context, v T) error {
				if !filter(v) {
					return nil
				}

				return w.Write(ctx, v)
			},
		}
	}
}

func WriteMapFn[T, U any](w Writer[U]) func(mapper func(T) U) Writer[T] {
	return func(mapper func(v T) U) Writer[T] {
		return WriterImpl[T]{
			Impl: func(ctx context.Context, vt T) error {
				return w.Write(ctx, mapper(vt))
			},
		}
	}
}
