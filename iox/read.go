package iox

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"io"
)

// readCloserImpl supports a functional implementation for io.ReadCloser.
type readCloserImpl struct {
	CImpl func() error
	RImpl func([]byte) (int, error)
}

// Close defers to readCloserImpl.CImpl.
func (impl readCloserImpl) Close() (err error) {
	if impl.CImpl == nil {
		return
	}

	return impl.CImpl()
}

// TODO: what is the significanse of ok? is it that T is default
// or io.EOF?
type Reader[T any] interface {
	Read(context.Context) (T, bool, error)
}

// ReaderImpl supports a functional implementation of Reader.
type ReaderImpl[T any] struct {
	Impl func(context.Context) (T, bool, error)
}

// Read defers to ReaderImpl.Impl.
func (impl ReaderImpl[T]) Read(
	ctx context.Context,
) (
	r T,
	ok bool,
	err error,
) {
	if impl.Impl == nil {
		err = io.EOF
		return
	}

	return impl.Impl(ctx)
}

type ReadCloser[T any] interface {
	io.Closer
	Reader[T]
}

// ReadCloserImpl supports a functional implementation of ReadCloser.
type ReadCloserImpl[T any] struct {
	CImpl func() error
	RImpl func(context.Context) (T, bool, error)
}

// Close defers to ReadCloserImpl.CImpl.
func (impl ReadCloserImpl[T]) Close() (err error) {
	if impl.CImpl == nil {
		return
	}

	return impl.CImpl()
}

// Read defers to ReadCloserImpl.RImpl.
func (impl ReadCloserImpl[T]) Read(
	ctx context.Context,
) (
	r T,
	ok bool,
	err error,
) {
	if impl.RImpl == nil {
		err = io.EOF
		return
	}

	return impl.RImpl(ctx)
}

// Read defers to readCloser.RImpl.
func (impl readCloserImpl) Read(p []byte) (n int, err error) {
	if impl.RImpl == nil {
		err = io.EOF
		return
	}

	return impl.RImpl(p)
}

func NewV2VReader[T any](vs ...T) Reader[T] {
	i := 0
	return ReaderImpl[T]{
		Impl: func(ctx context.Context) (v T, ok bool, err error) {
			if i > len(vs)-1 {
				return
			}

			i++
			return vs[i-1], true, nil
		},
	}
}

func NewD2VReader[T any](dec Decoder) Reader[T] {
	if dec == nil {
		return ReaderImpl[T]{}
	}

	return ReaderImpl[T]{
		Impl: func(ctx context.Context) (v T, ok bool, err error) {
			err = dec.Decode(&v)
			if errors.Is(err, io.EOF) {
				err = nil
				return
			}

			if err != nil {
				return
			}

			ok = true
			return
		},
	}
}

func NewB2VReadCloser[T any](
	r io.ReadCloser,
	d ...func(io.Reader) Decoder,
) (
	_ ReadCloser[T],
) {
	var dec Decoder
	dec = gob.NewDecoder(r)

	if len(d) > 0 {
		last := d[len(d)-1]
		if last != nil {
			dec = last(r)
		}
	}

	if r == nil {
		return ReadCloserImpl[T]{}
	}

	// TODO handle closed pipe.
	return ReadCloserImpl[T]{
		CImpl: r.Close,
		RImpl: NewD2VReader[T](dec).Read,
	}
}

func NewB2VReader[T any](r io.Reader, d ...func(io.Reader) Decoder) Reader[T] {
	rc := readCloserImpl{RImpl: r.Read}
	return NewB2VReadCloser[T](rc, d...)
}

func NewV2BReadCloser[T any](
	r ReadCloser[T],
	e ...func(io.Writer) Encoder,
) (
	_ io.ReadCloser,
) {
	buf := bytes.NewBuffer(nil)
	i2e := func(w io.Writer) Encoder { return gob.NewEncoder(w) }

	if len(e) > 0 {
		last := e[len(e)-1]
		if last != nil {
			i2e = last
		}
	}

	if r == nil {
		return readCloserImpl{}
	}

	enc := i2e(buf)
	return readCloserImpl{
		CImpl: r.Close,
		RImpl: func(p []byte) (n int, err error) {
			// TODO cont/ok consistency.
			v, cont, err := r.Read(context.Background())
			if !cont {
				return 0, io.EOF
			}
			if err != nil {
				return 0, err
			}

			err = enc.Encode(v)
			if err != nil {
				return 0, err
			}

			return buf.Read(p)
		},
	}
}

func NewV2BReader[T any](r Reader[T], e ...func(io.Writer) Encoder) io.Reader {
	rc := ReadCloserImpl[T]{RImpl: r.Read}
	return NewV2BReadCloser[T](rc, e...)
}

func ReadFilterFn[T any](r Reader[T]) func(filter func(T) bool) Reader[T] {
	return func(filter func(v T) bool) Reader[T] {
		return ReaderImpl[T]{
			Impl: func(ctx context.Context) (v T, ok bool, err error) {
				for v, ok, err := r.Read(ctx); ; v, ok, err = r.Read(ctx) {
					if err != nil || !ok {
						return v, ok, err
					}

					if !filter(v) {
						continue
					}

					return v, ok, err
				}
			},
		}
	}
}

func ReadMapFn[T, U any](r Reader[T]) func(mapper func(T) U) Reader[U] {
	return func(mapper func(T) U) Reader[U] {
		return ReaderImpl[U]{
			Impl: func(ctx context.Context) (vu U, ok bool, err error) {
				var vt T
				vt, ok, err = r.Read(ctx)
				if err != nil || !ok {
					return
				}

				vu = mapper(vt)
				return
			},
		}
	}
}

func ReadReduceFn[T any](r Reader[T]) func(reducer func(T, T) T) (T, error) {
	return func(reducer func(T, T) T) (c T, err error) {
		for v, ok, err := r.Read(nil); ; v, ok, err = r.Read(nil) {
			if err != nil || !ok {
				return c, err
			}

			c = reducer(c, v)
		}
	}
}
