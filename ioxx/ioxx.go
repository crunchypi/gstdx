package ioxx

import "io"

type Encoder interface {
	Encode(e any) error
}

type EncoderImpl struct {
	Impl func(e any) error
}

func (impl EncoderImpl) Encode(e any) error {
	if impl.Impl == nil {
		return nil
	}

	return impl.Impl(e)
}

type Decoder interface {
	Decode(e any) error
}

type DecoderImpl struct {
	Impl func(d any) error
}

func (impl DecoderImpl) Decode(d any) error {
	if impl.Impl == nil {
		return nil
	}

	return impl.Impl(d)
}

// Abbreviations.

type decoderFn = func(io.Reader) Decoder
type encoderFn = func(io.Writer) Encoder

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
