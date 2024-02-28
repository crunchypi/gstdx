package iox

/*
type ReadWriter[T, U any] interface {
	ReadWrite(context.Context, T) (U, bool, error)
}

type ReadWriterImpl[T, U any] struct {
	Impl func(context.Context, T) (U, bool, error)
}

func (impl ReadWriterImpl[T, U]) ReadWrite(
	ctx context.Context,
	v T,
) (
	r U,
	ok bool,
	err error,
) {
	if impl.Impl == nil {
		err = errors.New("ios: used ReadWriterImpl without impl")
		return
	}

	return impl.Impl(ctx, v)
}

*/
