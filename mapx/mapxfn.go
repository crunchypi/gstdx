package mapx

func FilterFn[K Key, V Val, M ~map[K]V](m M) func(rcv func(Pair[K, V]) bool) M {
	return func(filter func(Pair[K, V]) bool) M {
		return Filter(m, filter)
	}
}

func FilterKFn[K Key, V Val, M ~map[K]V](m M) func(rcv func(K) bool) M {
	return func(rcv func(K) bool) M {
		return FilterK(m, rcv)
	}
}

func FilterVFn[K Key, V Val, M ~map[K]V](m M) func(rcv func(V) bool) M {
	return func(rcv func(V) bool) M {
		return FilterV(m, rcv)
	}
}

func MapKFn[K1, K2 Key, V Val](m map[K1]V) func(f func(K1) K2) map[K2]V {
	return func(f func(K1) K2) map[K2]V {
		return MapK(m, f)
	}
}

func MapVFn[K Key, V1, V2 Val](m map[K]V1) func(f func(V1) V2) map[K]V2 {
	return func(f func(V1) V2) map[K]V2 {
		return MapV(m, f)
	}
}

func ReduceKFn[K Key, V Val](m map[K]V) func(f func(K, K) K) K {
	return func(f func(K, K) K) K {
		return ReduceK(m, f)
	}
}

func ReduceVFn[K Key, V Val](m map[K]V) func(f func(V, V) V) V {
	return func(f func(V, V) V) V {
		return ReduceV(m, f)
	}
}
