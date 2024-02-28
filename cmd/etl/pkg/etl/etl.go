package etl

import (
	"context"
	"math/rand"
	"time"

	"github.com/crunchypi/gstdx/iox"
)

// -----------------------------------------------------------------------------
// Pagination.
// -----------------------------------------------------------------------------

type Page struct {
	Skip  int
	Limit int
}

func newPageReader(n int, limit int) iox.Reader[Page] {
	var skip int

	return iox.ReaderImpl[Page]{
		Impl: func(ctx context.Context) (p Page, ok bool, err error) {
			if skip >= n {
				return
			}

			p.Skip = skip
			p.Limit = limit

			if skip+limit > n {
				p.Limit -= (skip + limit) - n
			}

			skip += limit
			ok = true
			return
		},
	}
}

func NewPageReader(r iox.Reader[int], limit int) iox.Reader[Page] {
	var n int
	var pages iox.Reader[Page]

	return iox.ReaderImpl[Page]{
		Impl: func(ctx context.Context) (p Page, ok bool, err error) {
			if pages == nil {
				n, ok, err = r.Read(ctx)
				if err != nil || !ok {
					return
				}

				pages = newPageReader(n, limit)
			}

			p, ok, err = pages.Read(ctx)
			if !ok {
				n, ok, err = r.Read(ctx)
				if err != nil || !ok {
					return
				}

				pages = newPageReader(n, limit)
				p, ok, err = pages.Read(ctx)
			}

			return p, ok, err
		},
	}
}

// Paged Buf Reader

// -----------------------------------------------------------------------------
// Sleep.
// -----------------------------------------------------------------------------

func newSleepFuzzReader[T any](r iox.Reader[T]) iox.Reader[T] {
	return iox.ReaderImpl[T]{
		Impl: func(ctx context.Context) (v T, ok bool, err error) {
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
	return iox.ReaderImpl[T]{
		Impl: func(ctx context.Context) (v T, ok bool, err error) {
			time.Sleep(d)
			return r.Read(ctx)
		},
	}
}

func NewSleepVReader[T any](r iox.Reader[T], d time.Duration) iox.Reader[T] {
	return iox.ReaderImpl[T]{
		Impl: func(ctx context.Context) (v T, ok bool, err error) {
			stamp := time.Now()

			v, ok, err = r.Read(ctx)
			if !ok || err != nil {
				return
			}

			time.Sleep(d - time.Now().Sub(stamp))
			return
		},
	}
}

func newSleepFuzzWriter[T any](w iox.Writer[T]) iox.Writer[T] {
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
	return iox.WriterImpl[T]{
		Impl: func(ctx context.Context, v T) error {
			time.Sleep(d)
			return w.Write(ctx, v)
		},
	}
}

func NewSleepVWriter[T any](w iox.Writer[T], d time.Duration) iox.Writer[T] {
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

// -----------------------------------------------------------------------------
// Ctx.
// -----------------------------------------------------------------------------

func NewCtxCancelledReader[T any](r iox.Reader[T]) iox.Reader[T] {
	return iox.ReaderImpl[T]{
		Impl: func(ctx context.Context) (v T, ok bool, err error) {
			if ctx != nil {
				select {
				case <-ctx.Done():
				default:
				}
			}

			v, ok, err = r.Read(ctx)
			return
		},
	}
}

// -----------------------------------------------------------------------------
// Stats.
// -----------------------------------------------------------------------------

type ReadStat[T any] struct {
	Tag   string
	Err   error
	Val   T
	Delta time.Duration
}

func NewStatsReader[T any](
	r iox.Reader[T],
	tag string,
) (
	rcv func(w iox.Writer[ReadStat[T]]) iox.Reader[T],
) {
	return func(w iox.Writer[ReadStat[T]]) iox.Reader[T] {
		return iox.ReaderImpl[T]{
			Impl: func(ctx context.Context) (v T, ok bool, err error) {
				stamp := time.Now()

				s := ReadStat[T]{}
				defer func() {
					s.Tag = tag
					s.Err = err
					s.Val = v
					s.Delta = time.Now().Sub(stamp)
					w.Write(ctx, s)
				}()

				v, ok, err = r.Read(ctx)
				return
			},
		}
	}
}

// -----------------------------------------------------------------------------
// Batch.
// -----------------------------------------------------------------------------

func NewBatchedVReader[T any](
	r iox.Reader[T],
	size int,
) (
	_ iox.Reader[T],
) {
	type result struct {
		v   T
		ok  bool
		err error
	}

	bufs := make([]result, size)
	bufp := iox.NewV2VReader[result]()

	return iox.ReaderImpl[T]{
		Impl: func(ctx context.Context) (v T, ok bool, err error) {
			_v, _ok, _ := bufp.Read(ctx)

			if !_ok {
				for i := 0; i < size; i++ {
					bufs[i].v, bufs[i].ok, bufs[i].err = r.Read(ctx)
				}

				bufp = iox.NewV2VReader[result](bufs...)
				_v, _ok, _ = bufp.Read(ctx)
			}

			v, ok, err = _v.v, _v.ok, _v.err
			return
		},
	}
}

func NewBatchedSReader[T any](
	r iox.Reader[T],
	size int,
) (
	_ iox.Reader[[]T],
) {
	return iox.ReaderImpl[[]T]{
		Impl: func(ctx context.Context) (s []T, ok bool, err error) {
			s = make([]T, 0, size)
			var v T
			for i := 0; i < size; i++ {
				v, ok, err = r.Read(ctx)
				if !ok || err != nil {
					break

				}

				s = append(s, v)
			}

			return s, len(s) != 0, err
		},
	}
}

func NewBatchedVWriter[T any](w iox.Writer[T], size int) iox.Writer[T] {
	buf := make([]T, 0, size)
	return iox.WriterImpl[T]{
		Impl: func(ctx context.Context, v T) (err error) {
			buf = append(buf, v)

			if len(buf) >= size {
				for i := 0; i < size; i++ {
					_v := buf[0]
					buf = buf[1:]
					err = w.Write(ctx, _v)
					if err != nil {
						return
					}
				}
			}

			return err
		},
	}
}

func NewBatchedSWriter[T any](w iox.Writer[[]T], size int) iox.Writer[T] {
	buf := make([]T, 0, size)
	return iox.WriterImpl[T]{
		Impl: func(ctx context.Context, v T) (err error) {
			buf = append(buf, v)

			if len(buf) >= size {
				err = w.Write(ctx, buf)
				buf = make([]T, 0, size)
			}

			return err
		},
	}
}

// Multi stats.
// Merged reader
