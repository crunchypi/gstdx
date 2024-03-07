package mapx

type Key = comparable
type Val = any

type Pair[K Key, V Val] struct {
	K K
	V V
}

func New[K Key, V Val](kv ...Pair[K, V]) map[K]V {
	r := make(map[K]V, len(kv))
	for _, pair := range kv {
		r[pair.K] = pair.V
	}

	return r
}

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

func FilterKFn[K Key, V Val, M ~map[K]V](m M) func(rcv func(K) bool) M {
	return func(f func(K) bool) M {
		return FilterFn(m)(
			func(pair Pair[K, V]) bool {
				return f(pair.K)
			},
		)
	}
}

func FilterVFn[K Key, V Val, M ~map[K]V](m M) func(rcv func(V) bool) M {
	return func(f func(V) bool) M {
		return FilterFn(m)(
			func(pair Pair[K, V]) bool {
				return f(pair.V)
			},
		)
	}
}

// TODO: Consider [K1, V1, K2, V2], it's clearer to reason with mentally.

func MapFn[K1, K2 Key, V1, V2 Val](m map[K1]V1) func(f func(Pair[K1, V1]) Pair[K2, V2]) map[K2]V2 {
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

func MapKFn[KI, KO Key, V Val](m map[KI]V) func(f func(KI) KO) map[KO]V {
	return func(f func(KI) KO) map[KO]V {
		return MapFn[KI, KO, V, V](m)(
			func(p Pair[KI, V]) Pair[KO, V] {
				return Pair[KO, V]{K: f(p.K), V: p.V}
			},
		)
	}
}

func MapVFn[K Key, VI, VO Val](m map[K]VI) func(f func(VI) VO) map[K]VO {
	return func(f func(VI) VO) map[K]VO {
		return MapFn[K, K, VI, VO](m)(
			func(p Pair[K, VI]) Pair[K, VO] {
				return Pair[K, VO]{K: p.K, V: f(p.V)}
			},
		)
	}
}

func ReduceFn[K Key, V Val, M ~map[K]V, P Pair[K, V]](m M) func(f func(P, P) P) P {
	return func(f func(acc P, curr P) P) (r P) {
		for k, v := range m {
			r = f(r, P{K: k, V: v})
		}

		return r
	}
}

func ReduceKFn[K Key, V Val](m map[K]V) func(f func(K, K) K) K {
	return func(f func(K, K) K) K {
		return ReduceFn(m)(
			func(acc, curr Pair[K, V]) Pair[K, V] {
				return Pair[K, V]{K: f(acc.K, curr.K)}
			},
		).K
	}
}

func ReduceVFn[K Key, V Val](m map[K]V) func(f func(V, V) V) V {
	return func(f func(V, V) V) V {
		return ReduceFn(m)(
			func(acc, curr Pair[K, V]) Pair[K, V] {
				return Pair[K, V]{V: f(acc.V, curr.V)}
			},
		).V
	}
}

func IntoClone[K Key, V Val, M ~map[K]V](m M) M {
	r := make(M, len(m))
	for k, v := range m {
		r[k] = v
	}

	return r
}

// TODO: Consider IntoSliceP, etc, to keep naming uniform.
func IntoPairs[K Key, V Val, M ~map[K]V](m M) []Pair[K, V] {
	r := make([]Pair[K, V], 0, len(m))
	for k, v := range m {
		r = append(r, Pair[K, V]{K: k, V: v})
	}

	return r
}

func IntoKeys[K Key, V Val, M ~map[K]V](m M) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}

	return r
}

func IntoVals[K Key, V Val, M ~map[K]V](m M) []V {
	r := make([]V, 0, len(m))
	for _, v := range m {
		r = append(r, v)
	}

	return r
}

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

func IntoChanV[K Key, V Val, M ~map[K]V](m M) <-chan V {
	ch := make(chan V)
	go func() {
		defer close(ch)

		for _, v := range m {
			ch <- v
		}
	}()

	return ch
}

func IntoGenerator[K Key, V Val, M ~map[K]V](m M) func() (Pair[K, V], bool) {
	ks := IntoKeys(m)
	i := 0

	return func() (p Pair[K, V], ok bool) {
		if i >= len(ks) || len(m) == 0 {
			return
		}

		p.K = ks[i]
		p.V = m[p.K]

		ok = true
		i++
		return
	}
}

func IntoGeneratorK[K Key, V Val, M ~map[K]V](m M) func() (K, bool) {
	g := IntoGenerator(m)
	return func() (k K, ok bool) {
		p, ok := g()
		return p.K, ok
	}
}

func IntoGeneratorV[K Key, V Val, M ~map[K]V](m M) func() (V, bool) {
	g := IntoGenerator(m)
	return func() (v V, ok bool) {
		p, ok := g()
		return p.V, ok
	}
}
