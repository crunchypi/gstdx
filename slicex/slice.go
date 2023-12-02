package slicex

import (
	"cmp"
	"context"
	"sort"
)

func New[T any](vs ...T) []T {
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

func Reduce[T any, S ~[]T](s S, f func(accumulate, current T) T) (r T) {
	if len(s) == 0 || f == nil {
		return r
	}

	for _, v := range s {
		r = f(r, v)
	}

	return r
}

func Contains[T comparable, S ~[]T](s S, v T) int {
	r := 0
	for _, elm := range s {
		if elm == v {
			r++
		}
	}

	return r
}

func Equals[T comparable, S ~[]T](s1, s2 S) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			return false
		}
	}

	return true
}

func Sort[T cmp.Ordered, S ~[]T](s S) {
	sort.Slice(s, func(i, j int) bool { return s[i] < s[j] })
}

func Sorted[T cmp.Ordered, S ~[]T](s S) S {
	r := make([]T, len(s))
	copy(r, s)
	Sort(r)
	return r
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

func IntoChan[T any, S ~[]T](ctx context.Context, s S) <-chan T {
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
