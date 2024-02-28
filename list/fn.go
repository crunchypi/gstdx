package list

func Filter[T any](l *List[T], f func(T) bool) *List[T] {
	r := New[T]()
	l.Iter(
		func(i int, v T) bool {
			if f(v) {
				r.Put(r.Len(), v)
			}

			return true
		},
	)

	return r
}

func Map[T, U any](l *List[T], f func(T) U) *List[U] {
	r := New[U]()
	l.Iter(
		func(i int, v T) bool {
			r.Put(r.Len(), f(v))
			return true
		},
	)

	return r
}

func Reduce[T any](l *List[T], f func(acc, cur T) T) (r T) {
	l.Iter(
		func(i int, v T) bool {
			r = f(r, v)
			return true
		},
	)

	return r
}

func IntoSlice[T any](l *List[T]) []T {
	r := make([]T, 0, l.Len())
	l.Iter(
		func(i int, v T) bool {
			r = append(r, v)
			return true
		},
	)

	return r
}

func IntoMapK[K comparable, V any](l *List[K], f func(K) V) map[K]V {
	r := make(map[K]V, l.Len())
	l.Iter(
		func(i int, k K) bool {
			r[k] = f(k)
			return true
		},
	)

	return r
}

func IntoMapV[K comparable, V any](l *List[V], f func(V) K) map[K]V {
	r := make(map[K]V, l.Len())
	l.Iter(
		func(i int, v V) bool {
			r[f(v)] = v
			return true
		},
	)

	return r
}

func IntoChan[T any](l *List[T]) <-chan T {
	r := make(chan T)
	go func() {
		defer close(r)

		l.Iter(
			func(i int, v T) bool {
				r <- v
				return true
			},
		)
	}()

	return r
}

func IntoGen[T any](l *List[T]) func() (T, bool) {
	i := 0
	return func() (v T, ok bool) {
		i++
		return l.Get(i - 1)
	}
}
