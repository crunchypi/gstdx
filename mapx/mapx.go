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

// ReduceFn returns a func which reduces 'm' using the given reducer.
// Example:
//
//	m := ReduceFn(map[int]int{1: 2, 2: 3})(
//		func(a, c Pair[int, int]) Pair[int, int] {
//			a.K += c.K
//			a.V += c.V
//			return a
//		},
//	)
//
//	// m is Pair[int, int]{3, 5}
func ReduceFn[K Key, V Val, M ~map[K]V, P Pair[K, V]](m M) func(f func(P, P) P) P {
	return func(f func(acc P, curr P) P) (r P) {
		if f == nil {
			return
		}

		for k, v := range m {
			r = f(r, P{K: k, V: v})
		}

		return r
	}
}

// ReduceFn returns a func which reduces 'm' into keys using the given reducer,
// which operates on keys.
// Example:
//
//	v := ReduceKFn(map[int]int{1: 1, 2: 2})(
//		func(accum, curr int) int {
//			return accum + curr
//		},
//	)
//
//	// v is 3.
func ReduceKFn[K Key, V Val](m map[K]V) func(f func(K, K) K) K {
	return func(f func(K, K) K) (k K) {
		if f == nil {
			return
		}

		return ReduceFn(m)(
			func(acc, curr Pair[K, V]) Pair[K, V] {
				return Pair[K, V]{K: f(acc.K, curr.K)}
			},
		).K
	}
}

// ReduceVFn returns a func which reduces 'm' into values using the given reducer,
// which operates on values.
// Example:
//
//	v := ReduceVFn(map[int]int{1: 1, 2: 2})(
//		func(accum, curr int) int {
//			return accum + curr
//		},
//	)
//
//	// v is 3.
func ReduceVFn[K Key, V Val](m map[K]V) func(f func(V, V) V) V {
	return func(f func(V, V) V) (v V) {
		if f == nil {
			return
		}

		return ReduceFn(m)(
			func(acc, curr Pair[K, V]) Pair[K, V] {
				return Pair[K, V]{V: f(acc.V, curr.V)}
			},
		).V
	}
}

// IntoClone returns a shallow copy of the given map.
func IntoClone[K Key, V Val, M ~map[K]V](m M) M {
	r := make(M, len(m))
	for k, v := range m {
		r[k] = v
	}

	return r
}

// IntoSlice returns a slice of pairs from the given map.
func IntoSlice[K Key, V Val, M ~map[K]V](m M) []Pair[K, V] {
	r := make([]Pair[K, V], 0, len(m))
	for k, v := range m {
		r = append(r, Pair[K, V]{K: k, V: v})
	}

	return r
}

// IntoSliceK returns a slice of keys from the given map.
func IntoSliceK[K Key, V Val, M ~map[K]V](m M) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}

	return r
}

// IntoSliceV returns a slice of vals from the given map.
func IntoSliceV[K Key, V Val, M ~map[K]V](m M) []V {
	r := make([]V, 0, len(m))
	for _, v := range m {
		r = append(r, v)
	}

	return r
}

// IntoChan returns a chan of pairs from the given map. Items are fed into the
// chan from a new goroutine so any use of 'm' should stop after this call.
func IntoChan[K Key, V Val, M ~map[K]V](m M) <-chan Pair[K, V] {
	ch := make(chan Pair[K, V])
	go func() {
		defer close(ch)

		for k, v := range m {
			ch <- Pair[K, V]{K: k, V: v}
		}
	}()

	return ch
}

// IntoChanK returns a chan of keys from the given map. Keys are fed into the
// chan from a new goroutine so any use of 'm' should stop after this call.
func IntoChanK[K Key, V Val, M ~map[K]V](m M) <-chan K {
	ch := make(chan K)
	go func() {
		defer close(ch)

		for k := range m {
			ch <- k
		}
	}()

	return ch
}
