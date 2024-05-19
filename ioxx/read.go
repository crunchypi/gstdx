package ioxx

import (
	"bytes"
	"context"
	"encoding/gob"
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
// Converters.
// -----------------------------------------------------------------------------

func NewValueReaderFn[T any](r io.Reader) func(f decoderFn) Reader[T] {
	return func(f func(io.Reader) Decoder) Reader[T] {
		// TODO nils.

		var d Decoder = gob.NewDecoder(r)
		if f != nil {
			if _d := f(r); _d != nil {
				d = _d
			}
		}

		return ReaderImpl[T]{
			Impl: func(ctx context.Context) (v T, err error) {
				err = d.Decode(&v)
				return
			},
		}
	}
}

func NewValueReadCloserFn[T any](r io.ReadCloser) func(f decoderFn) ReadCloser[T] {
	return func(f func(io.Reader) Decoder) ReadCloser[T] {
		return ReadCloserImpl[T]{
			ImplC: r.Close,
			ImplR: NewValueReaderFn[T](r)(f).Read,
		}
	}
}

func NewByteReaderFn[T any](r Reader[T]) func(f encoderFn) io.Reader {
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

func NewByteReadCloserFn[T any](r ReadCloser[T]) func(f encoderFn) io.ReadCloser {
	return func(f func(io.Writer) Encoder) io.ReadCloser {
		return readWriteCloserImpl{
			ImplC: r.Close,
			ImplR: NewByteReaderFn(r)(f).Read,
		}
	}
}
