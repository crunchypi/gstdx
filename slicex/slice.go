package slicex

func New[T any](vs ...T) []T {
	if len(vs) == 0 {
		return []T{}
	}

	return vs
}

func Filter[T any, S ~[]T](s S, f func(T) bool) S {
	if len(s) == 0 {
		return []T{}
	}

	if f == nil {
		r := make([]T, len(s))
		copy(r, s)
		return r
	}

	r := make([]T, 0, len(s))
	for _, v := range s {
		if f(v) {
			r = append(r, v)
		}
	}

	return r
}

func FilterFn[T any, S ~[]T](s S) func(filter func(T) bool) S {
	return func(f func(T) bool) S {
		return Filter(s, f)
	}
}

func Map[T, U any, S ~[]T](s S, f func(T) U) []U {
	if len(s) == 0 || f == nil {
		return []U{}
	}

	r := make([]U, 0, len(s))
	for _, v := range s {
		r = append(r, f(v))
	}

	return r
}

func MapFn[T, U any, S ~[]T](s S) func(mapper func(T) U) []U {
	return func(f func(T) U) []U {
		return Map(s, f)
	}
}

func Reduce[T any, S ~[]T](s S, f func(accumulate, current T) T) (r T) {
	if len(s) == 0 || f == nil {
		return r
	}

	for _, v := range s {
		r = f(r, v)
	}

	return r
}

func ReduceFn[T any, S ~[]T](s S) func(reducer func(acc, curr T) T) T {
	return func(f func(T, T) T) T {
		return Reduce(s, f)
	}
}

func IntoClone[T any, S ~[]T](s S) S {
	r := make([]T, 0, len(s))
	for _, v := range s {
		r = append(r, v)
	}

	return r
}

func IntoMapK[K comparable, V any, S ~[]K](s S, f func(K) V) map[K]V {
	if len(s) == 0 {
		return map[K]V{}
	}

	r := make(map[K]V, len(s))
	for _, k := range s {
		var v V
		if f != nil {
			v = f(k)
		}

		r[k] = v
	}

	return r
}

func IntoMapKFn[K comparable, V any, S ~[]K](s S) func(func(K) V) map[K]V {
	return func(f func(K) V) map[K]V {
		return IntoMapK(s, f)
	}
}

func IntoMapV[K comparable, V any, S ~[]V](s S, f func(V) K) map[K]V {
	if len(s) == 0 {
		return map[K]V{}
	}

	r := make(map[K]V, len(s))
	for _, v := range s {
		var k K
		if f != nil {
			k = f(v)
		}

		r[k] = v
	}

	return r
}

func IntoMapVFn[K comparable, V any, S ~[]V](s S) func(func(V) K) map[K]V {
	return func(f func(V) K) map[K]V {
		return IntoMapV(s, f)
	}
}

func IntoChan[T any, S ~[]T](s S) <-chan T {
	if len(s) == 0 {
		r := make(chan T)
		close(r)
		return r
	}

	r := make(chan T)
	go func() {
		defer close(r)

		for _, v := range s {
			r <- v
		}
	}()

	return r
}

func IntoGenerator[T any, S ~[]T](s S) func() (v T, ok bool) {
	if len(s) == 0 {
		return func() (v T, ok bool) { return }
	}

	i := 0
	return func() (v T, ok bool) {
		if i >= len(s) {
			return
		}

		v = s[i]
		ok = true
		i++
		return
	}
}
