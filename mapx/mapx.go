package mapx

type Key = comparable // Abbreviation, used a lot.
type Val = any        // Abbreviation, used a lot.

type Pair[K Key, V Val] struct {
	K K
	V V
}

// New returns a map[K]V from the given pairs.
func New[K Key, V Val](kv ...Pair[K, V]) map[K]V {
	r := make(map[K]V, len(kv))
	for _, pair := range kv {
		r[pair.K] = pair.V

	}

	return r
}

// FilterFn returns a func which filters 'm' using a given filter func.
// Example:
//
//	m := FilterFn(map[int]int{1: 1, 2: 2, 3: 3})(
//		func(p Pair[int, int]) bool {
//			return p.K > 1
//		},
//	)
//
//	// m is map[int]int{2:2, 3:3}
func FilterFn[K Key, V Val, M ~map[K]V](m M) func(rcv func(Pair[K, V]) bool) M {
	return func(f func(Pair[K, V]) bool) (r M) {
		if len(m) == 0 {
			return map[K]V{}
		}

		if f == nil {
			r := make(map[K]V, len(m))
			for k, v := range m {
				r[k] = v
			}

			return r
		}

		r = make(map[K]V, len(m))
		for k, v := range m {
			if f(Pair[K, V]{K: k, V: v}) {
				r[k] = v
			}
		}

		return r
	}
}

// FilterKFn returns a func which filters 'm' using a given filter func,
// which bases evaluation on keys.
// Example:
//
//	m := FilterKFn(map[int]int{1: 1, 2: 2, 3: 3})(
//		func(k int) bool {
//			return k > 1
//		},
//	)
//
//	// m is map[int]int{2:2, 3:3}
func FilterKFn[K Key, V Val, M ~map[K]V](m M) func(rcv func(K) bool) M {
	return func(f func(K) bool) M {
		return FilterFn(m)(
			func(pair Pair[K, V]) bool {
				return f == nil || f(pair.K)
			},
		)
	}
}

// FilterVFn returns a func which filters 'm' using a given filter func,
// which bases evaluation on map vals.
// Example:
//
//	m := FilterVFn(map[int]int{1: 1, 2: 2, 3: 3})(
//		func(k int) bool {
//			return k > 1
//		},
//	)
//
//	// m is map[int]int{2:2, 3:3}
func FilterVFn[K Key, V Val, M ~map[K]V](m M) func(rcv func(V) bool) M {
	return func(f func(V) bool) M {
		return FilterFn(m)(
			func(pair Pair[K, V]) bool {
				return f == nil || f(pair.V)
			},
		)
	}
}

// MapFn returns a func which maps 'm' using a given mapper func.
// Example:
//
//	m := MapFn[int, int, int, int](map[int]int{1: 2, 3: 4})(
//		func(p Pair[int, int]) Pair[int, int] {
//			p.K++
//			p.V++
//			return p
//		},
//	)
//
//	// m is map[int]int{2:3, 4:5}
func MapFn[K1 Key, V1 Val, K2 Key, V2 Val](
	m map[K1]V1,
) (
	f func(f func(Pair[K1, V1]) Pair[K2, V2]) map[K2]V2,
) {
	return func(f func(Pair[K1, V1]) Pair[K2, V2]) map[K2]V2 {
		if len(m) == 0 || f == nil {
			return map[K2]V2{}
		}

		r := make(map[K2]V2, len(m))
		for k1, v1 := range m {
			p1 := Pair[K1, V1]{K: k1, V: v1}
			p2 := f(p1)
			r[p2.K] = p2.V
		}

		return r
	}
}

// MapKFn returns a func which maps 'm' using the given filter func,
// which bases conversion on keys.
// Example:
//
//	m := MapKFn[int, int](map[int]int{1: 2, 2: 3})(
//		func(k int) int {
//			return k + 1
//		},
//	)
//
//	// m is map[int]int{2:2, 3:3}
func MapKFn[KI, KO Key, V Val](m map[KI]V) func(f func(KI) KO) map[KO]V {
	return func(f func(KI) KO) map[KO]V {
		if f == nil {
			return map[KO]V{}
		}

		return MapFn[KI, V, KO, V](m)(
			func(p Pair[KI, V]) Pair[KO, V] {
				return Pair[KO, V]{K: f(p.K), V: p.V}
			},
		)
	}
}

// MapVFn returns a func which maps 'm' using the given map func, which bases
// conversion on map values.
// Example:
//
//	m := MapVFn[int, int, int](map[int]int{1: 1, 2: 2})(
//		func(k int) int {
//			return k + 1
//		},
//	)
//
//	// m is map[int]int{1:2, 2:3}
func MapVFn[K Key, VI, VO Val](m map[K]VI) func(f func(VI) VO) map[K]VO {
	return func(f func(VI) VO) map[K]VO {
		if f == nil {
			return map[K]VO{}
		}

		return MapFn[K, VI, K, VO](m)(
			func(p Pair[K, VI]) Pair[K, VO] {
				return Pair[K, VO]{K: p.K, V: f(p.V)}
			},
		)
	}
}
