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

// NewByteReadCloserFn is identical to NewByteReaderFn, except that it returns
// an io.ReadCloser which may close 'r'.
func NewByteReadCloserFn[T any](r ReadCloser[T]) func(f encoderFn) io.ReadCloser {
	return func(f func(io.Writer) Encoder) io.ReadCloser {
		if r == nil {
			return readWriteCloserImpl{}
		}

		return readWriteCloserImpl{
			ImplC: r.Close,
			ImplR: NewByteReaderFn(r)(f).Read,
		}
	}
}

// NewBatchedValueReader returns a reader which batches 'r' into slices with
// the specified 'size'.  If the size is not set (or negative), it will be set
// to a small number. Note that the last slice may contain values when the
// returned reader gives an io.EOF.
func NewBatchedValueReader[T any](r Reader[T], size int) Reader[[]T] {
	if r == nil {
		return ReaderImpl[[]T]{}
	}

	if size <= 0 {
		size = 8
	}

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

			return s, err
		},
	}
}

// NewUnbatchedValueReader returns a reader of T from a reader of []T.
// Note that there is some internal buffering, so you may want to use this
// with caution as an unread buffer may cause value loss.
func NewUnbatchedValueReader[T any](r Reader[[]T]) Reader[T] {
	if r == nil {
		return ReaderImpl[T]{}
	}

	var errCache error
	var buf []T
	return ReaderImpl[T]{
		Impl: func(ctx context.Context) (val T, err error) {
			if len(buf) > 0 {
				val = buf[0]
				buf = buf[1:]
				return
			}

			if errCache != nil {
				err = errCache
				return
			}

			buf, err = r.Read(ctx)

			switch {
			case len(buf) == 0 && err != nil:
				return val, err
			case len(buf) == 0 && err == nil:
				return val, io.EOF
			case len(buf) != 0 && err != nil:
				errCache = err
				err = nil
			case len(buf) != 0 && err == nil:
			}

			val = buf[0]
			buf = buf[1:]
			return
		},
	}
}
