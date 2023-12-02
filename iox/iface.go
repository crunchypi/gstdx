package iox

import "context"

type Reader[T any] interface {
	Read(context.Context) (T, error)
}

type Writer[T any] interface {
	Write(context.Context, T) error
}

type ReadWriter[T, U any] interface {
	ReadWrite(context.Context, T) (U, error)
}
