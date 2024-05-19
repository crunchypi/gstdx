package ioxx

/*
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
// Converters.
// -----------------------------------------------------------------------------

func NewByteWriterFn[T any](w Writer[T]) func(d Decoder) io.Writer {
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

func NewByteWriteCloserFn[T any](w WriteCloser[T]) func(d Decoder) io.WriteCloser {
	return func(d Decoder) io.WriteCloser {
		return readWriteCloserImpl{
			ImplC: w.Close,
			ImplW: NewByteWriterFn(w)(d).Write,
		}
	}
}

func NewValueWriterFn[T any](w io.Writer) func(f encoderFn) Writer[T] {
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

func NewValueWriteCloserFn[T any](w io.WriteCloser) func(f encoderFn) WriteCloser[T] {
	return func(f func(io.Writer) Encoder) WriteCloser[T] {
		return WriteCloserImpl[T]{
			ImplC: w.Close,
			ImplW: NewValueWriterFn[T](w)(f).Write,
		}
	}
}
*/
