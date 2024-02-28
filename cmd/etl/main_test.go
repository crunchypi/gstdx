package main

import (
	"context"
	"net/http"
	"time"

	"github.com/crunchypi/gstdx/cmd/etl/pkg/etl"
	"github.com/crunchypi/gstdx/iox"
)

func NewGetReader[T any]() iox.Reader[T] {
	resp, _ := http.Get("")
	r := iox.NewB2VReader[T](resp.Body)

	return iox.ReaderImpl[T]{
		Impl: func(ctx context.Context) (v T, ok bool, err error) {
			v, ok, err = r.Read(ctx)

			return
		},
	}
}

func NewPostWriter[T any]() iox.Writer[T] {
	return iox.WriterImpl[T]{
		Impl: func(ctx context.Context, v T) error {
			r := iox.NewV2BReader[T](iox.NewV2VReader[T](v))
			http.Post("", "", r)
			return nil
		},
	}
}

func x() {
	nr := iox.NewV2VReader(1, 2, 3)
	sr := etl.NewSleepVReader(nr, time.Second)
	pr := etl.NewPageReader(sr, 2)
	xr := iox.ReadMapFn[etl.Page, etl.Page](pr)(
		func(page etl.Page) etl.Page {
			return page
		},
	)
}
