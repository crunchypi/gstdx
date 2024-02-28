package iox

import (
	"context"
	"io"
)

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

// -----------------------------------------------------------------------------

func PipeV2V[T any](
	ctx context.Context,
	w Writer[T],
	r Reader[T],
) (
	err error,
) {
	for v, ok, err := r.Read(ctx); ok; v, ok, err = r.Read(ctx) {
		if err != nil {
			return err
		}

		err = w.Write(ctx, v)
		if err != nil {
			return err
		}
	}

	return err
}

func PipeB2V[T any](
	ctx context.Context,
	w Writer[T],
	r io.Reader,
	d ...func(io.Reader) Decoder,
) (
	err error,
) {

	xr := NewB2VReader[T](r, d...)
	for v, ok, err := xr.Read(ctx); ok; v, ok, err = xr.Read(ctx) {
		if err != nil {
			return err
		}

		err = w.Write(ctx, v)
		if err != nil {
			return err
		}
	}

	return err
}

// TODO: Naming comes into conflict with io.Pipe, this is more of a "copy"
func PipeV2B[T any](
	ctx context.Context,
	w io.Writer,
	r Reader[T],
	e ...func(io.Writer) Encoder,
) (
	err error,
) {
	ctx.Deadline()
	xw := NewV2BWriter[T](w, e...)
	for v, ok, err := r.Read(ctx); ok; v, ok, err = r.Read(ctx) {
		if err != nil {
			return err
		}

		xw.Write(ctx, v)
	}

	return err
}
