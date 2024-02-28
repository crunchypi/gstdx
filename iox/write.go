package iox

import (
	"bytes"
	"context"
	"encoding/gob"
	"io"
)

type writeCloserImpl struct {
	CImpl func() error
	WImpl func([]byte) (int, error)
}

func (impl writeCloserImpl) Close() (err error) {
	if impl.CImpl == nil {
		return
	}

	return impl.CImpl()
}

func (impl writeCloserImpl) Write(p []byte) (n int, err error) {
	if impl.WImpl == nil {
		err = io.ErrClosedPipe
		return
	}

	return impl.WImpl(p)
}

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

type WriteCloser[T any] interface {
	io.Closer
	Writer[T]
}

// WriteCloserImpl supports a functional implementation of WriteCloser.
type WriteCloserImpl[T any] struct {
	CImpl func() error
	WImpl func(context.Context, T) error
}

// Write defers to WriteCloserImpl.CImpl.
func (impl WriteCloserImpl[T]) Close() error {
	if impl.CImpl == nil {
		return nil
	}

	return impl.CImpl()
}

// Write defers to WriteCloserImpl.WImpl.
func (impl WriteCloserImpl[T]) Write(
	ctx context.Context,
	v T,
) (
	err error,
) {
	if impl.WImpl == nil {
		err = io.ErrClosedPipe
		return
	}

	return impl.WImpl(ctx, v)
}

func NewB2VWriteCloser[T any](
	w WriteCloser[T],
	d ...Decoder,
) (
	_ io.WriteCloser,
) {
	buf := bytes.NewBuffer(nil)
	var dec Decoder
	dec = gob.NewDecoder(buf)

	if len(d) > 0 {
		last := d[len(d)-1]
		if last != nil {
			dec = last
		}
	}

	if w == nil {
		return writeCloserImpl{}
	}

	return writeCloserImpl{
		CImpl: w.Close,
		WImpl: func(p []byte) (n int, err error) {
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

func NewB2VWriter[T any](w Writer[T], e ...Decoder) io.Writer {
	wc := WriteCloserImpl[T]{WImpl: w.Write}
	return NewB2VWriteCloser[T](wc, e...)
}

func NewV2BWriteCloser[T any](
	w io.WriteCloser,
	e ...func(io.Writer) Encoder,
) (
	_ WriteCloser[T],
) {
	// What is i2e, and normalize it.
	buf := bytes.NewBuffer(nil)
	i2e := func(w io.Writer) Encoder { return gob.NewEncoder(w) }
	if len(e) > 0 {
		i2e = e[len(e)-1]
	}

	if w == nil {
		return WriteCloserImpl[T]{}
	}

	enc := i2e(buf)
	return WriteCloserImpl[T]{
		CImpl: w.Close,
		WImpl: func(ctx context.Context, v T) error {
			err := enc.Encode(v)
			if err != nil {
				return err
			}

			_, err = buf.WriteTo(w)
			return err
		},
	}
}

func NewV2BWriter[T any](w io.Writer, e ...func(io.Writer) Encoder) Writer[T] {
	wc := writeCloserImpl{WImpl: w.Write}
	return NewV2BWriteCloser[T](wc, e...)
}

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
