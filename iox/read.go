package iox

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"io"
)

// -----------------------------------------------------------------------------
// New Reader iface + impl.
// -----------------------------------------------------------------------------

type Reader[T any] interface {
	Read(context.Context) (T, error)
}

// ReaderImpl supports a functional implementation of Reader.
type ReaderImpl[T any] struct {
	Impl func(context.Context) (T, error)
}

// Read defers to ReaderImpl.Impl.
func (impl ReaderImpl[T]) Read(ctx context.Context) (r T, err error) {
	if impl.Impl == nil {
		err = io.EOF
		return
	}

	return impl.Impl(ctx)
}

// -----------------------------------------------------------------------------
// New ReadCloser iface + impl.
// -----------------------------------------------------------------------------

type ReadCloser[T any] interface {
	io.Closer
	Reader[T]
}

// ReadCloserImpl supports a functional implementation of ReadCloser.
type ReadCloserImpl[T any] struct {
	ImplC func() error
	ImplR func(context.Context) (T, error)
}

// Close defers to ReadCloserImpl.CImpl.
func (impl ReadCloserImpl[T]) Close() (err error) {
	if impl.ImplC == nil {
		return
	}

	return impl.ImplC()
}

// Read defers to ReadCloserImpl.RImpl.
func (impl ReadCloserImpl[T]) Read(ctx context.Context) (r T, err error) {
	if impl.ImplR == nil {
		err = io.EOF
		return
	}

	return impl.ImplR(ctx)
}

// -----------------------------------------------------------------------------
// Factory funcs.
// -----------------------------------------------------------------------------

func NewV2VReader[T any](vs ...T) Reader[T] {
	i := 0
	return ReaderImpl[T]{
		Impl: func(ctx context.Context) (v T, err error) {
			if i > len(vs)-1 {
				return v, io.EOF
			}

			i++
			return vs[i-1], nil
		},
	}
}

func NewD2VReader[T any](dec Decoder) Reader[T] {
	if dec == nil {
		return ReaderImpl[T]{}
	}

	return ReaderImpl[T]{
		Impl: func(ctx context.Context) (v T, err error) {
			err = dec.Decode(&v)
			return
		},
	}
}

func NewB2VReaderFn[T any](r io.Reader) func(f func(io.Reader) Decoder) Reader[T] {
	return func(f func(io.Reader) Decoder) Reader[T] {
		// TODO nils.

		var d Decoder = gob.NewDecoder(r)
		if f != nil {
			if _d := f(r); _d != nil {
				d = _d
			}
		}

		return ReaderImpl[T]{Impl: NewD2VReader[T](d).Read}
	}
}

func NewB2VReadCloserFn[T any](r io.ReadCloser) func(f func(io.Reader) Decoder) ReadCloser[T] {
	return func(f func(io.Reader) Decoder) ReadCloser[T] {
		return ReadCloserImpl[T]{
			ImplC: r.Close,
			ImplR: NewB2VReaderFn[T](r)(f).Read,
		}
	}
}

func NewV2BReaderFn[T any](r Reader[T]) func(f func(io.Writer) Encoder) io.Reader {
	// TODO nils like below:
	if r == nil {
		r = ReaderImpl[T]{}
	}

	return func(f func(io.Writer) Encoder) io.Reader {
		buf := bytes.NewBuffer(nil)
		enc := func(w io.Writer) Encoder { return gob.NewEncoder(w) }(buf)

		if f != nil {
			if _e := f(buf); _e != nil {
				enc = _e
			}
		}

		return readWriteCloserImpl{
			ImplR: func(p []byte) (n int, err error) {
				v, err := r.Read(context.Background())
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
}

func NewV2BReadCloserFn[T any](r ReadCloser[T]) func(f func(io.Writer) Encoder) io.ReadCloser {
	return func(f func(io.Writer) Encoder) io.ReadCloser {
		return readWriteCloserImpl{
			ImplC: r.Close,
			ImplR: NewV2BReaderFn(r)(f).Read,
		}
	}
}

// -----------------------------------------------------------------------------
// Batching.
// -----------------------------------------------------------------------------

func NewBatchedVReader[T any](r Reader[T], size int) Reader[T] {
	type result struct {
		v   T
		err error
	}

	rw := NewV2VReadWriter[result]()

	return ReaderImpl[T]{
		Impl: func(ctx context.Context) (v T, err error) {
			_v, _err := rw.Read(ctx)
			if errors.Is(_err, io.EOF) {
				for i := 0; i < size; i++ {
					res := result{}
					res.v, res.err = r.Read(ctx)
					rw.Write(ctx, res)
				}

				_v, _ = rw.Read(ctx)
			}

			v, err = _v.v, _v.err
			return
		},
	}
}

// The above implementation should be refactored from this later.
//func NewBatchedVReader[T any](r Reader[T], size int) Reader[T] {
//	type result struct {
//		v   T
//		ok  bool
//		err error
//	}
//
//	bufs := make([]result, size)
//	bufp := NewV2VReader[result]()
//
//	return ReaderImpl[T]{
//		Impl: func(ctx context.Context) (v T, ok bool, err error) {
//			_v, _ok, _ := bufp.Read(ctx)
//
//			if !_ok {
//				for i := 0; i < size; i++ {
//					bufs[i].v, bufs[i].ok, bufs[i].err = r.Read(ctx)
//				}
//
//				bufp = NewV2VReader[result](bufs...)
//				_v, _ok, _ = bufp.Read(ctx)
//			}
//
//			v, ok, err = _v.v, _v.ok, _v.err
//			return
//		},
//	}
//}

func NewBatchedSReader[T any](r Reader[T], size int) Reader[[]T] {
	return ReaderImpl[[]T]{
		Impl: func(ctx context.Context) (s []T, err error) {
			s = make([]T, 0, size)
			var v T
			for i := 0; i < size; i++ {
				v, err = r.Read(ctx)
				if err != nil {
					break

				}

				s = append(s, v)
			}

			if len(s) == 0 {
				err = io.EOF
			}

			return s, err
		},
	}
}

// -----------------------------------------------------------------------------
// Functional.
// -----------------------------------------------------------------------------

func ReadFilterFn[T any](r Reader[T]) func(filter func(T) bool) Reader[T] {
	return func(filter func(v T) bool) Reader[T] {
		return ReaderImpl[T]{
			Impl: func(ctx context.Context) (v T, err error) {
				for v, err := r.Read(ctx); ; v, err = r.Read(ctx) {
					if err != nil {
						return v, err
					}

					if !filter(v) {
						continue
					}

					return v, err
				}
			},
		}
	}
}

func ReadMapFn[T, U any](r Reader[T]) func(mapper func(T) U) Reader[U] {
	return func(mapper func(T) U) Reader[U] {
		return ReaderImpl[U]{
			Impl: func(ctx context.Context) (vu U, err error) {
				var vt T
				vt, err = r.Read(ctx)
				if err != nil {
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
		for v, err := r.Read(nil); !errors.Is(err, io.EOF); v, err = r.Read(nil) {
			if err != nil {
				return c, err
			}

			c = reducer(c, v)
		}

		return
	}
}
