package log

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/crunchypi/gstdx/iox"
)

var defaultLogger *slog.Logger

func init() {
	opts := slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				tv, ok := a.Value.Any().(time.Time)
				if !ok {
					return a
				}
				s := tv.Format("060102.150405.0000")

				a.Value = slog.StringValue(s)
				return a
			}

			if a.Key == slog.MessageKey {
				return slog.Attr{}
			}

			return a
		},
	}

	h := slog.NewJSONHandler(os.Stdout, &opts)
	defaultLogger = slog.New(h)
}

func clearMap[K comparable, V any](m map[K]V) {
	for k := range m {
		delete(m, k)
	}
}

type NewArgs[T any] struct {
	Reader  iox.Reader[T]
	Logger  *slog.Logger
	Tag     string
	CtxKeys []string
}

func NewReader[T any](args NewArgs[T]) iox.Reader[T] {
	if args.Logger == nil {
		args.Logger = defaultLogger
	}

	l := args.Logger.With("tag", args.Tag)
	n := 9
	i := 0
	b := iox.NewV2VReadWriter[T]()
	m := make(map[string]any, len(args.CtxKeys))

	return iox.ReaderImpl[T]{
		Impl: func(ctx context.Context) (v T, err error) {
			v, err = args.Reader.Read(ctx)
			if err == nil {
				return
			}

			clearMap(m)
			for _, k := range args.CtxKeys {
				m[k] = ctx.Value(k)
			}

			ll := l.With()
			ll = ll.With("val", v)
			ll = ll.With("ctx", m)

			err = errors.New("dank")
			if err != nil && !errors.Is(err, io.EOF) {
				ll.Error("", "err:", err)
			}

			if i >= n {
				i = 0

				for i := 0; i < n; i++ {

				}
			}

			b.Write(ctx, v)
			ll.Info("")
			return
		},
	}
}
