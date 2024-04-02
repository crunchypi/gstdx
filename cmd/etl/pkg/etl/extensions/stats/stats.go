package stats

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/crunchypi/gstdx/iox"
)

type Stat[T any] struct {
	Tag     string         `json:"tag"`
	Val     T              `json:"val"`
	Err     error          `json:"err"`
	CtxVals map[string]any `json:"ctxVals"`
	Delta   time.Duration  `json:"delta"`
}

func (s Stat[T]) IntoMap() map[string]any {
	return map[string]any{
		"tag":     s.Tag,
		"val":     s.Val,
		"err":     s.Err,
		"ctxVals": s.CtxVals,
		"delta":   s.Delta,
	}
}

func NewReader[T any](
	r iox.Reader[T],
	tag string,
	ctxKeys ...string,
) (
	_ iox.Reader[Stat[T]],
) {
	if r == nil {
		r = iox.ReaderImpl[T]{}
	}

	return iox.ReaderImpl[Stat[T]]{
		Impl: func(ctx context.Context) (s Stat[T], err error) {
			stamp := time.Now()
			v, err := r.Read(ctx)

			if err != nil && errors.Is(err, io.EOF) {
				return
			}

			s.Tag = tag
			s.Val = v
			s.Err = err
			s.CtxVals = make(map[string]any, len(ctxKeys))
			s.Delta = time.Now().Sub(stamp)

			if ctx != nil {
				for _, k := range ctxKeys {
					s.CtxVals[k] = ctx.Value(k)
				}
			}

			return
		},
	}
}

func NewTeeReaderFn[T any](
	r iox.Reader[T],
	tag string,
	ctxKeys ...string,
) (
	rcv func(w iox.Writer[Stat[T]]) iox.Reader[T],
) {
	if r == nil {
		r = iox.ReaderImpl[T]{}
	}

	return func(w iox.Writer[Stat[T]]) iox.Reader[T] {
		if w == nil {
			return r
		}

		sr := NewReader(r, tag, ctxKeys...)
		return iox.ReaderImpl[T]{
			Impl: func(ctx context.Context) (v T, err error) {
				s, err := sr.Read(ctx)
				if err != nil && errors.Is(err, io.EOF) {
					return s.Val, err
				}

				w.Write(ctx, s)
				return s.Val, err
			},
		}
	}
}

func NewWriter[T any](
	w iox.Writer[Stat[T]],
	tag string,
	ctxKeys ...string,
) (
	_ iox.Writer[T],
) {
	stamp := time.Now()
	return iox.WriterImpl[T]{
		Impl: func(ctx context.Context, v T) (err error) {
			newStamp := time.Now()
			s := Stat[T]{}
			s.Tag = tag
			s.Err = err
			s.Val = v
			s.CtxVals = make(map[string]any, len(ctxKeys))
			s.Delta = newStamp.Sub(stamp)

			if ctx != nil {
				for _, k := range ctxKeys {
					s.CtxVals[k] = ctx.Value(k)
				}
			}

			w.Write(ctx, s)
			stamp = newStamp
			return err
		},
	}
}

func NewTeeWriterFn[T any](
	tw iox.Writer[T],
	tag string,
	ctxKeys ...string,
) (
	rcv func(sw iox.Writer[Stat[T]]) iox.Writer[T],
) {
	if tw == nil {
		tw = iox.WriterImpl[T]{}
	}

	return func(sw iox.Writer[Stat[T]]) iox.Writer[T] {
		if sw == nil {
			return tw
		}

		return iox.WriterImpl[T]{
			Impl: func(ctx context.Context, v T) (err error) {
				stamp := time.Now()
				defer func() {
					if err != nil && errors.Is(err, io.EOF) {
						return
					}

					s := Stat[T]{}
					s.Tag = tag
					s.Err = err
					s.Val = v
					s.CtxVals = make(map[string]any, len(ctxKeys))
					s.Delta = time.Now().Sub(stamp)

					if ctx != nil {
						for _, k := range ctxKeys {
							s.CtxVals[k] = ctx.Value(k)
						}
					}

					sw.Write(ctx, s)
				}()

				err = tw.Write(ctx, v)
				return err
			},
		}
	}
}
