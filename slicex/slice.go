package slicex

import (
	"context"
)

func New[T any](vs ...T) []T {
	return vs
}

func Filter[T any](s []T, f func(T) bool) []T {
	if len(s) == 0 || f == nil {
		return []T{}
	}

	r := make([]T, 0, len(s))
	for _, v := range s {
		if f(v) {
			r = append(r, v)
		}
	}

	return r
}

func FilterFn[T any](s []T) func(filter func(T) bool) []T {
	return func(f func(T) bool) []T {
		return Filter(s, f)
	}
}

func Map[T, U any](s []T, f func(T) U) []U {
	if len(s) == 0 || f == nil {
		return []U{}
	}

	r := make([]U, 0, len(s))
	for _, v := range s {
		r = append(r, f(v))
	}

	return r
}

func MapFn[T, U any](s []T) func(mapper func(T) U) []U {
	return func(f func(T) U) []U {
		return Map(s, f)
	}
}

func Reduce[T any](s []T, f func(accumulate, current T) T) (r T) {
	if len(s) == 0 || f == nil {
		return r
	}

	for _, v := range s {
		r = f(r, v)
	}

	return r
}

func ReduceFn[T any](s []T) func(reducer func(accumulate, current T) T) T {
	return func(f func(T, T) T) T {
		return Reduce(s, f)
	}
}

func IntoClone[T any](s []T) []T {
	r := make([]T, 0, len(s))
	for _, v := range s {
		r = append(r, v)
	}

	return r
}

func IntoMapK[K comparable, V any](s []K, f func(K) V) map[K]V {
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

func IntoMapKFn[K comparable, V any](s []K) func(func(K) V) map[K]V {
	return func(f func(K) V) map[K]V {
		return IntoMapK(s, f)
	}
}

func IntoMapV[K comparable, V any](s []V, f func(V) K) map[K]V {
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

func IntoMapVFn[K comparable, V any](s []V) func(func(V) K) map[K]V {
	return func(f func(V) K) map[K]V {
		return IntoMapV(s, f)
	}
}

func IntoChan[T any](ctx context.Context, s []T) <-chan T {
	if len(s) == 0 {
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

		for _, v := range s {
			select {
			case <-ctx.Done():
				return
			case r <- v:
			}
		}
	}()

	return r
}

func IntoGenerator[T any](s []T) func() (v T, ok bool) {
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
