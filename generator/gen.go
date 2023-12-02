package generator

import "context"

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

func Filter[T any](g Gen[T], f func(T) bool) Gen[T] {
	if g == nil {
		return New[T]()
	}
	if f == nil {
		return g
	}

	return func() (v T, cont bool) {
		v, cont = g()
		for ; f(v) && cont; v, cont = g() {
		}

		return v, cont
	}
}

func FilterFn[T any](g Gen[T]) func(rcv func(T) bool) Gen[T] {
	return func(f func(v T) bool) Gen[T] {
		return Filter(g, f)
	}
}

func Map[T, U any](g Gen[T], f func(T) U) Gen[U] {
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

func MapFn[T, U any](g Gen[T]) func(rcv func(T) U) Gen[U] {
	return func(rcv func(T) U) Gen[U] {
		return Map(g, rcv)
	}
}

func Reduce[T any](g Gen[T], f func(acc, cur T) T) (r T) {
	if g == nil || f == nil {
		return r
	}

	for v, cont := g(); cont; v, cont = g() {
		r = f(r, v)
	}

	return r
}

func ReduceFn[T any](g Gen[T]) func(rcv func(T, T) T) T {
	return func(rcv func(T, T) T) T {
		return Reduce(g, rcv)
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

func IntoMapK[K comparable, V any](g Gen[K], f func(K) V) map[K]V {
	if g == nil || f == nil {
		return map[K]V{}
	}

	r := make(map[K]V)
	for k, cont := g(); cont; k, cont = g() {
		r[k] = f(k)
	}

	return r
}

func IntoMapKFn[K comparable, V any](g Gen[K]) func(func(K) V) map[K]V {
	return func(rcv func(K) V) map[K]V {
		return IntoMapK(g, rcv)
	}
}

func IntoMapV[K comparable, V any](g Gen[V], f func(V) K) map[K]V {
	if g == nil || f == nil {
		return map[K]V{}
	}

	r := make(map[K]V)
	for v, cont := g(); cont; v, cont = g() {
		r[f(v)] = v
	}

	return r
}

func IntoMapVFn[K comparable, V any](g Gen[V]) func(func(V) K) map[K]V {
	return func(rcv func(V) K) map[K]V {
		return IntoMapV(g, rcv)
	}
}

func IntoChan[T any](ctx context.Context, g Gen[T]) <-chan T {
	if g == nil {
		r := make(chan T)
		close(r)
		return r
	}

	if ctx == nil {
		ctx = context.Background()
	}

	r := make(chan T)
	go func() {
		defer close(r)
		for v, cont := g(); cont; v, cont = g() {
			select {
			case <-ctx.Done():
				return
			case r <- v:
			}
		}

	}()

	return r
}
