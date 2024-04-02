package sleep

import (
	"context"
	"math/rand"
	"time"

	"github.com/crunchypi/gstdx/iox"
)

func newSleepFuzzReader[T any](r iox.Reader[T]) iox.Reader[T] {
	if r == nil {
		r = iox.ReaderImpl[T]{}
	}

	return iox.ReaderImpl[T]{
		Impl: func(ctx context.Context) (v T, err error) {
			l := 0
			h := 1000

			sl := time.Duration(rand.Intn(h-l) + l)
			sl *= time.Millisecond

			time.Sleep(sl)
			return r.Read(ctx)
		},
	}
}

func NewSleepRReader[T any](r iox.Reader[T], d time.Duration) iox.Reader[T] {
	if r == nil {
		r = iox.ReaderImpl[T]{}
	}

	return iox.ReaderImpl[T]{
		Impl: func(ctx context.Context) (v T, err error) {
			time.Sleep(d)
			return r.Read(ctx)
		},
	}
}

func NewSleepVReader[T any](r iox.Reader[T], d time.Duration) iox.Reader[T] {
	if r == nil {
		r = iox.ReaderImpl[T]{}
	}

	return iox.ReaderImpl[T]{
		Impl: func(ctx context.Context) (v T, err error) {
			stamp := time.Now()

			v, err = r.Read(ctx)
			if err != nil {
				return
			}

			time.Sleep(d - time.Now().Sub(stamp))
			return
		},
	}
}

func newSleepFuzzWriter[T any](w iox.Writer[T]) iox.Writer[T] {
	if w == nil {
		w = iox.WriterImpl[T]{}
	}

	return iox.WriterImpl[T]{
		Impl: func(ctx context.Context, v T) error {

			l := 0
			h := 1000

			sl := time.Duration(rand.Intn(h-l) + l)
			sl *= time.Millisecond

			time.Sleep(sl)
			return w.Write(ctx, v)
		},
	}
}

func NewSleepRWriter[T any](w iox.Writer[T], d time.Duration) iox.Writer[T] {
	if w == nil {
		w = iox.WriterImpl[T]{}
	}

	return iox.WriterImpl[T]{
		Impl: func(ctx context.Context, v T) error {
			time.Sleep(d)
			return w.Write(ctx, v)
		},
	}
}

func NewSleepVWriter[T any](w iox.Writer[T], d time.Duration) iox.Writer[T] {
	if w == nil {
		w = iox.WriterImpl[T]{}
	}

	stamp := time.Now()
	return iox.WriterImpl[T]{
		Impl: func(ctx context.Context, v T) (err error) {
			stampTmp := time.Now()
			time.Sleep(d - stampTmp.Sub(stamp))

			err = w.Write(ctx, v)
			stamp = time.Now()
			return
		},
	}
}
