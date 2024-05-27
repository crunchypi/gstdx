package iox

import (
	"bytes"
	"context"
	"encoding/gob"
	"io"
)

// -----------------------------------------------------------------------------
// New Reader iface + impl.
// -----------------------------------------------------------------------------

// Reader reads T, it is intended to act as a generic variant of io.Reader.
type Reader[T any] interface {
	Read(context.Context) (T, error)
}

// ReaderImpl implements Reader with it's Read method by deferring to 'Impl'.
// This is for convenience, as you may use a functional implementation of Reader
// without defining a new type (that's done for you here).
type ReaderImpl[T any] struct {
	Impl func(context.Context) (T, error)
}

// Read implements Reader by deferring to the internal "Impl" func.
// If the internal "Impl" is not set, an io.EOF will be returned.
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

// ReadCloser groups Reader with io.Closer.
type ReadCloser[T any] interface {
	io.Closer
	Reader[T]
}

// ReadCloserImpl implements Reader and io.Closer with it's methods by deferring
// to ImplC (closer) and ImplR (reader). This is for convenience, as you may use
// a functional implementation of the interfaces wihout defining a new type.
type ReadCloserImpl[T any] struct {
	ImplC func() error
	ImplR func(context.Context) (T, error)
}

// Read implements Closer by deferring to the internal "ImplC" func.
// If the internal "ImplC" func is nil, nothing will happen.
func (impl ReadCloserImpl[T]) Close() (err error) {
	if impl.ImplC == nil {
		return
	}

	return impl.ImplC()
}

// Read implements Reader by deferring to the internal "ImplR" func.
// If the internal "ImplR" is not set, an io.EOF will be returned.
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

// NewValueReaderFn creates a new T reader from an io.Reader and Decoder.
// It simply reads bytes from 'r', decodes them, and passes them along to the
// caller. As such, the decoder must match the encoder used to create the bytes.
// If 'r' is nil, an empty Reader is returned; if 'f' is nil, the decoder is set
// to gob.NewDecoder. Example:
//
//	// Used as io.Reader
//	b := bytes.NewBuffer(nil)
//
//	// Using json encoder, so the decoder has to be json in NewValueReaderFn
//	json.NewEncoder(b).Encode("test1")
//	json.NewEncoder(b).Encode("test2")
//
//	r := NewValueReaderFn[string](b)(
//		func(r io.Reader) Decoder {
//			return json.NewDecoder(r)
//		},
//	)
//
//	t.Log(r.Read(context.Background())) // "test1" <nil>
//	t.Log(r.Read(context.Background())) // "test2" <nil>
//	t.Log(r.Read(context.Background())) // "", io.EOF
func NewValueReaderFn[T any](r io.Reader) func(f decoderFn) Reader[T] {
	return func(f func(io.Reader) Decoder) Reader[T] {
		if r == nil {
			return ReaderImpl[T]{}
		}

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

// NewValueReadCloserFn is identical to NewValueReaderFn, except that it returns
// a ReadCloser which may close 'r'.
func NewValueReadCloserFn[T any](r io.ReadCloser) func(f decoderFn) ReadCloser[T] {
	return func(f func(io.Reader) Decoder) ReadCloser[T] {
		if r == nil {
			return ReadCloserImpl[T]{}
		}

		return ReadCloserImpl[T]{
			ImplC: r.Close,
			ImplR: NewValueReaderFn[T](r)(f).Read,
		}
	}
}

// NewByteReaderFn creates an io.Reader from a Reader and Encoder.
// It simply reads values from 'r', encodes them, and passes them along to the
// caller. As such, when decoding values from the returned io.Reader one should
// use a decoder which matches the encoder passed here. If 'r' is nil, an
// empty (not nil) io.Reader is returned; if 'f' is nil, the encoder is set to
// gob.NewEncoder. Example:
//
//	// First encode some values into bytes.
//	b := bytes.NewBuffer(nil)
//	json.NewEncoder(b).Encode("test1")
//	json.NewEncoder(b).Encode("test2")
//
//	// Conversion to a value reader, then back to a byte reader.
//	vr := NewValueReaderFn[string](b)(func(r io.Reader) Decoder { return json.NewDecoder(r) })
//	br := NewByteReaderFn[string](vr)(func(w io.Writer) Encoder { return json.NewEncoder(w) })
//
//	// Instantly pass it to a decoder just so we may log out the values.
//	dec := json.NewDecoder(br)
//	val := ""
//
//	t.Log(dec.Decode(&val), val) // <nil>, "test1"
//	t.Log(dec.Decode(&val), val) // <nil>, "test2"
//	t.Log(dec.Decode(&val), val) // EOF, ""
func NewByteReaderFn[T any](r Reader[T]) func(f encoderFn) io.Reader {
	return func(f func(io.Writer) Encoder) io.Reader {
		if r == nil {
			r = ReaderImpl[T]{}
		}

		b := bytes.NewBuffer(nil)
		e := Encoder(gob.NewEncoder(b))
		if f != nil {
			if _e := f(b); _e != nil {
				e = _e
			}
		}

		return readWriteCloserImpl{
			ImplR: func(p []byte) (n int, err error) {
				v, err := r.Read(context.Background())
				if err != nil {
					return 0, err
				}

				err = e.Encode(v)
				if err != nil {
					return 0, err
				}

				return b.Read(p)
			},
		}
	}
}
