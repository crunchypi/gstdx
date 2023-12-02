package slicex

func FilterFn[T any, S ~[]T](s S) func(filter func(T) bool) S {
	return func(f func(T) bool) S {
		return Filter(s, f)
	}
}

func MapFn[T, U any, S ~[]T](s S) func(mapper func(T) U) []U {
	return func(f func(T) U) []U {
		return Map(s, f)
	}
}

func ReduceFn[T any, S ~[]T](s S) func(reducer func(accumulate, current T) T) T {
	return func(f func(T, T) T) T {
		return Reduce(s, f)
	}
}

func IntoMapKFn[K comparable, V any, S ~[]K](s S) func(func(K) V) map[K]V {
	return func(f func(K) V) map[K]V {
		return IntoMapK(s, f)
	}
}

func IntoMapVFn[K comparable, V any, S ~[]V](s S) func(func(V) K) map[K]V {
	return func(f func(V) K) map[K]V {
		return IntoMapV(s, f)
	}
}
