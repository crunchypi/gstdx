package mapx

import "context"

type Key comparable
type Val any

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

func Filter[K Key, V Val, M ~map[K]V](
	m M,
	f func(Pair[K, V]) bool,
) (
	r M,
) {
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

func FilterK[K Key, V Val, M ~map[K]V](m M, f func(K) bool) M {
	r := make(map[K]V, len(m))
	for k, v := range m {
		if f(k) {
			r[k] = v
		}
	}

	return r
}

func FilterV[K Key, V Val, M ~map[K]V](
	m M,
	f func(V) bool,
) (
	r M,
) {
	r = make(map[K]V, len(m))
	for k, v := range m {
		if f(v) {
			r[k] = v
		}
	}

	return r
}

func Map[K1, K2 Key, V1, V2 Val](
	m map[K1]V1,
	f func(K1, V1) (K2, V2),
) (
	r map[K2]V2,
) {
	if len(m) == 0 || f == nil {
		return map[K2]V2{}
	}

	r = make(map[K2]V2, len(m))
	for k1, v1 := range m {
		k2, v2 := f(k1, v1)
		r[k2] = v2
	}

	return r
}

func MapK[KI, KO Key, V Val](m map[KI]V, f func(KI) KO) map[KO]V {
	return Map(m, func(k KI, v V) (KO, V) { return f(k), v })
}

func MapV[K Key, VI, VO Val](m map[K]VI, f func(VI) VO) map[K]VO {
	return Map(m, func(k K, v VI) (K, VO) { return k, f(v) })
}

func ReduceK[K Key, V Val, M ~map[K]V](m M, f func(K, K) K) (r K) {
	for k := range m {
		r = f(r, k)
	}

	return r
}

func ReduceV[K Key, V Val, M ~map[K]V](m M, f func(V, V) V) (r V) {
	for _, v := range m {
		r = f(r, v)
	}

	return r
}

func IntoClone[K Key, V Val, M ~map[K]V](m M) M {
	r := make(M, len(m))
	for k, v := range m {
		r[k] = v
	}

	return r
}

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

func IntoChanP[K Key, V Val, M ~map[K]V](
	ctx context.Context,
	m M,
) (
	r <-chan Pair[K, V],
) {
	ch := make(chan Pair[K, V])
	go func() {
		defer close(ch)

		for k, v := range m {
			select {
			case <-ctx.Done():
				return
			case ch <- Pair[K, V]{K: k, V: v}:
			}
		}
	}()

	return ch
}

func IntoChanK[K Key, V Val, M ~map[K]V](
	ctx context.Context,
	m M,
) (
	r <-chan K,
) {
	ch := make(chan K)
	go func() {
		defer close(ch)

		for k := range m {
			select {
			case <-ctx.Done():
				return
			case ch <- k:
			}
		}
	}()

	return ch
}

func IntoChanV[K Key, V Val, M ~map[K]V](
	ctx context.Context,
	m M,
) (
	r <-chan V,
) {
	ch := make(chan V)
	go func() {
		defer close(ch)

		for _, v := range m {
			select {
			case <-ctx.Done():
				return
			case ch <- v:
			}
		}
	}()

	return ch
}
