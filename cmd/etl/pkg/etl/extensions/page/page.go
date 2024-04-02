package page

import (
	"context"
	"errors"
	"io"

	"github.com/crunchypi/gstdx/iox"
)

type Page struct {
	Skip  int
	Limit int
}

func newPageReader(n int, limit int) iox.Reader[Page] {
	var skip int

	return iox.ReaderImpl[Page]{
		Impl: func(ctx context.Context) (p Page, err error) {
			if skip >= n {
				return p, io.EOF
			}

			p.Skip = skip
			p.Limit = limit

			if skip+limit > n {
				p.Limit -= (skip + limit) - n
			}

			skip += limit
			return
		},
	}
}

func NewPageReader(r iox.Reader[int], limit int) iox.Reader[Page] {
	var n int
	var pages iox.Reader[Page]

	return iox.ReaderImpl[Page]{
		Impl: func(ctx context.Context) (p Page, err error) {
			// First.
			if pages == nil {
				n, err = r.Read(ctx)
				if err != nil {
					return
				}

				pages = newPageReader(n, limit)
			}

			// Next, and reset.
			p, err = pages.Read(ctx)
			if errors.Is(err, io.EOF) {
				n, err = r.Read(ctx)
				if err != nil {
					return
				}

				// Retry.
				pages = newPageReader(n, limit)
				p, err = pages.Read(ctx)
			}

			return p, err
		},
	}
}
