package chanx

func New[T any](vs ...T) <-chan T {
	if len(vs) == 0 {
		r := make(chan T)
		close(r)
		return r
	}

	r := make(chan T)
	go func() {
		defer close(r)

		for _, v := range vs {
			r <- v
		}
	}()

	return r
}

func FilterFn[T any](ch <-chan T) func(func(T) bool) <-chan T {
	return func(f func(T) bool) <-chan T {
		if ch == nil {
			r := make(chan T)
			close(r)
			return r
		}

		if f == nil {
			return ch
		}

		r := make(chan T)
		go func() {
			defer close(r)

			for v := range ch {
				if f(v) {
					r <- v
				}
			}
		}()

		return r
	}
}

func MapFn[T, U any](ch <-chan T) func(func(T) U) <-chan U {
	return func(f func(T) U) <-chan U {
		if ch == nil || f == nil {
			r := make(chan U)
			close(r)
			return r
		}

		r := make(chan U)
		go func() {
			defer close(r)

			for v := range ch {
				r <- f(v)
			}
		}()

		return r
	}
}

func ReduceFn[T any](ch <-chan T) func(func(acc, curr T) T) (r T) {
	return func(f func(acc T, curr T) T) (r T) {
		if ch == nil || f == nil {
			return r
		}

		for v := range ch {
			r = f(r, v)
		}

		return r
	}
}

func IntoSlice[T any](ch <-chan T, size ...int) []T {
	if ch == nil {
		return []T{}
	}

	l := 8
	if len(size) > 0 {
		for _, v := range size {
			l += v
		}
	}

	r := make([]T, 0, l)
	for v := range ch {
		r = append(r, v)
	}

	return r
}

func IntoMapKFn[K comparable, V any](ch <-chan K) func(f func(K) V) map[K]V {
	return func(f func(K) V) map[K]V {
		if ch == nil || f == nil {
			return map[K]V{}
		}

		r := make(map[K]V, 8)
		for k := range ch {
			r[k] = f(k)
		}

		return r
	}
}

func IntoMapVFn[K comparable, V any](ch <-chan V) func(f func(V) K) map[K]V {
	return func(f func(V) K) map[K]V {
		if ch == nil || f == nil {
			return map[K]V{}
		}

		r := make(map[K]V, 8)
		for v := range ch {
			r[f(v)] = v
		}

		return r
	}
}

func IntoGen[T any](ch <-chan T) func() (v T, cont bool) {
	if ch == nil {
		return func() (v T, cont bool) { return }
	}

	return func() (v T, cont bool) {
		v, cont = <-ch
		return v, cont
	}
}
