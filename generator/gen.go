package generator

type Gen[T any] func() (v T, cont bool)

func New[T any](vs ...T) Gen[T] {
	if len(vs) == 0 {
		return func() (v T, cont bool) { return }
	}

	i := 0
	return func() (v T, cont bool) {
		if i >= len(vs) {
			return
		}

		i++
		return vs[i-1], true
	}
}

func FilterFn[T any](g Gen[T]) func(rcv func(T) bool) Gen[T] {
	return func(f func(v T) bool) Gen[T] {
		if g == nil {
			return New[T]()
		}
		if f == nil {
			return g
		}

		return func() (v T, cont bool) {
			v, cont = g()
			for ; cont && !f(v); v, cont = g() {
			}

			return v, cont
		}
	}
}

func MapFn[T, U any](g Gen[T]) func(rcv func(T) U) Gen[U] {
	return func(f func(T) U) Gen[U] {
		if g == nil || f == nil {
			return New[U]()
		}

		return func() (vu U, cont bool) {
			vt, cont := g()
			if !cont {
				return vu, cont
			}

			vu = f(vt)
			return vu, cont
		}
	}
}

func ReduceFn[T any](g Gen[T]) func(rcv func(T, T) T) T {
	return func(f func(acc T, curr T) T) (r T) {
		if g == nil || f == nil {
			return r
		}

		for v, cont := g(); cont; v, cont = g() {
			r = f(r, v)
		}

		return r
	}
}

func IntoS[T any](g Gen[T], size ...int) []T {
	if g == nil {
		return []T{}
	}

	l := 8
	if len(size) > 0 {
		l = size[0]
	}

	r := make([]T, 0, l)
	for v, cont := g(); cont; v, cont = g() {
		r = append(r, v)
	}

	return r
}

func IntoMapKFn[K comparable, V any](g Gen[K]) func(func(K) V) map[K]V {
	return func(f func(K) V) map[K]V {
		if g == nil || f == nil {

			return map[K]V{}
		}

		r := make(map[K]V)
		for k, cont := g(); cont; k, cont = g() {
			r[k] = f(k)
		}

		return r
	}
}

func IntoMapVFn[K comparable, V any](g Gen[V]) func(func(V) K) map[K]V {
	return func(f func(V) K) map[K]V {
		if g == nil || f == nil {
			return map[K]V{}
		}

		r := make(map[K]V)
		for v, cont := g(); cont; v, cont = g() {
			r[f(v)] = v
		}

		return r
	}
}

func IntoChan[T any](g Gen[T]) <-chan T {
	if g == nil {
		r := make(chan T)
		close(r)
		return r
	}

	r := make(chan T)
	go func() {
		defer close(r)
		for v, cont := g(); cont; v, cont = g() {
			r <- v
		}

	}()

	return r
}
