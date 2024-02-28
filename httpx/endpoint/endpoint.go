package endpoint

import (
	"net/http"

	"github.com/crunchypi/gstdx/iox"
)

func EndNotImplemented() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}

type EndReaderArgs[T any] struct {
	URLParamNamespace string
	Reader[T]
}

func (args EndReaderArgs[_]) Ok() (ok bool) {
	ok = true
	ok = ok && args.Reader != nil

	return ok
}

func EndReader[T any](args EndReaderArgs[T]) http.HandlerFunc {
	if !args.Ok() {
		return EndNotImplemented()
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Optional url param namespace.
		ns := HFTryParseURLParam(r, args.URLParamNamespace)
		if args.URLParamNamespace != "" && ns == "" {
			s := "endpoint.EndReader: expected namespace"
			http.Error(w, s, http.StatusInternalServerError)
			return
		}

		// Forward.
		v, ssc, err := args.Read(r.Context(), ns)

		// Respond.
		HFRespond(w, v, ssc, err)
	}
}

type EndWriterArgs[T any] struct {
	URLParamNamespace string
	Writer[T]
}

func (args EndWriterArgs[_]) Ok() (ok bool) {
	ok = true
	ok = ok && args.Writer != nil

	return ok
}

func EndWriter[T any](args EndWriterArgs[T]) http.HandlerFunc {
	if !args.Ok() {
		return EndNotImplemented()
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Optional url param namespace.
		ns := HFTryParseURLParam(r, args.URLParamNamespace)
		if args.URLParamNamespace != "" && ns == "" {
			s := "endpoint.EndWriter: expected namespace"
			http.Error(w, s, http.StatusInternalServerError)
			return
		}

		// Request body.
		v, cont := HFUnmarshal[T](w, r)
		if !cont {
			return
		}

		// Forward.
		ssc, err := args.Write(r.Context(), ns, v)

		// Respond.
		HFRespond(w, nil, ssc, err)
	}
}

type EndReadWriterArgs[I, O any] struct {
	URLParamNamespace string
	ReadWriter[I, O]
}

func (args EndReadWriterArgs[_, _]) Ok() (ok bool) {
	ok = true
	ok = ok && args.ReadWriter != nil

	return ok
}

func EndReadWriter[I, O any](args EndReadWriterArgs[I, O]) http.HandlerFunc {
	if !args.Ok() {
		return EndNotImplemented()
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Optional url param namespace.
		ns := HFTryParseURLParam(r, args.URLParamNamespace)
		if args.URLParamNamespace != "" && ns == "" {
			s := "endpoint.EndReadWriter: expected namespace"
			http.Error(w, s, http.StatusInternalServerError)
			return
		}

		// Request body.
		in, cont := HFUnmarshal[I](w, r)
		if !cont {
			return
		}

		// Forward.
		out, ssc, err := args.ReadWrite(r.Context(), "", in)

		// Respond.
		HFRespond(w, out, ssc, err)
	}
}

type EndXArgs[I, O any] struct {
	R iox.Reader[O]
	W iox.Writer[I]
}

func EndX[I, O any](args EndXArgs[I, O]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		a := iox.NewB2VReader[I](r.Body)
		for v, cont, err := a.Read(ctx); cont; v, cont, err = a.Read(ctx) {
			if err != nil {
				// TODO
			}

			args.W.Write(ctx, v)
		}

		b := iox.NewWriter[O](w)
		for v, cont, err := args.R.Read(ctx); cont; v, cont, err = args.R.Read(ctx) {
			if err != nil {
				// TODO
			}

			b.Write(ctx, v)
		}

	}
}
