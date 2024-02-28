package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/crunchypi/gstdx/generator"
	"github.com/crunchypi/gstdx/iox"
)

func MidMethodWrap(s string, f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != s {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		f(w, r)
	}
}

type EndSaveResourceArgs[T any] struct {
	Writer iox.Writer[T]
}

func EndSaveResource[T any](args EndSaveResourceArgs[T]) http.HandlerFunc {
	dec := func(r io.Reader) iox.Decoder { return json.NewDecoder(r) }

	return func(w http.ResponseWriter, r *http.Request) {
		r.Body.Close()

		err := iox.PipeB2V[T](r.Context(), args.Writer, r.Body, dec)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

type EndLoadResourceArgs[T any] struct {
	Reader iox.Reader[T]
}

func EndLoadResource[T any](args EndLoadResourceArgs[T]) http.HandlerFunc {
	enc := func(w io.Writer) iox.Encoder { return json.NewEncoder(w) }
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		err := iox.PipeV2B(r.Context(), w, args.Reader, enc)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func main() {
	type T string

	g := generator.New[T]("1", "2", "3")

	routes := map[string]http.HandlerFunc{
		"/": EndSaveResource(
			EndSaveResourceArgs[T]{
				Writer: iox.WriterImpl[T]{
					Impl: func(_ context.Context, v T) error {
						fmt.Println(v)
						return nil
					},
				},
			},
		),
		"/g": EndLoadResource(
			EndLoadResourceArgs[T]{
				Reader: iox.ReaderImpl[T]{
					Impl: func(_ context.Context) (v T, ok bool, err error) {
						v, ok = g()
						return
					},
				},
			},
		),
	}

	for k, v := range routes {
		http.Handle(k, MidMethodWrap("POST", v))
	}

	http.ListenAndServe(":8080", nil)
}
