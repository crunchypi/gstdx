package iox

import (
	"io"
)

// -----------------------------------------------------------------------------
// Encoder & Decoder, used for interoperability of io.Reader and io.Writer
// with Reader and Writer of this package.
// -----------------------------------------------------------------------------

// Encoder encodes values into binary form. It is used for converting io.Reader
// and io.Writer into Reader and Writer (of this pkg), see funcs for docs.
// Encoder uses pre-generics semantics for interoperability with pre-generics
// packages. Some commonly used encoders are:
//   - gob.NewEncoder(bytes.NewBuffer(nil))
//   - json.NewEncoder(bytes.NewBuffer(nil))
type Encoder interface {
	Encode(e any) error
}

// EncoderImpl implements Encoder with it's Encode method by deferring to 'Impl'.
// This is for convenience, as you may use functional implementation of Encoder
// without defining a new type (that's done for you here).
type EncoderImpl struct {
	Impl func(e any) error
}

// Encode implements Encoder by deferring to the internal "Impl" func.
// If the internal "Impl" is not set, an io.ErrClosedPipe will be returned.
func (impl EncoderImpl) Encode(e any) error {
	if impl.Impl == nil {
		return io.ErrClosedPipe
	}

	return impl.Impl(e)
}

// Decoder decodes values from binary form. It is used for converting io.Reader
// and io.Writer into Reader and Writer (of this pkg), see funcs for docs.
// Decoder uses pre-generics semantics for interoperability with pre-generics
// packages. Some commonly used encoders are:
//   - gob.NewDecoder(bytes.NewBuffer(nil))
//   - json.NewDecoder(bytes.NewBuffer(nil))
type Decoder interface {
	Decode(e any) error
}

// DecoderImpl implements Decoder with it's Decode method by deferring to 'Impl'.
// This is for convenience, as you may use functional implementation of Decoder
// without defining a new type (that's done for you here).
type DecoderImpl struct {
	Impl func(d any) error
}

// Decode implements Decoder by deferring to the internal "Impl" func.
// If the internal "Impl" is not set, an io.EOF will be returned.
func (impl DecoderImpl) Decode(d any) error {
	if impl.Impl == nil {
		return io.EOF
	}

	return impl.Impl(d)
}

// -----------------------------------------------------------------------------
// Implementation io.Reader, io.Writer, io.ReadWriter and closer variants.
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
