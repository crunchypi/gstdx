package endpoint

import (
	"context"
)

type Reader[T any] interface {
	Read(
		ctx context.Context,
		namespace string,
	) (
		v T,
		ssc int,
		err error,
	)
}

type Writer[T any] interface {
	Write(
		ctx context.Context,
		namespace string,
		v T,
	) (
		ssc int,
		err error,
	)
}

type ReadWriter[I, O any] interface {
	ReadWrite(
		ctx context.Context,
		namespace string,
		in I,
	) (
		out O,
		ssc int,
		err error,
	)
}
