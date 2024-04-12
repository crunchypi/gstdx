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
