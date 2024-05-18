package generator

// Gen represents a generator func.
type Gen[T any] func() (v T, cont bool)

// New returns a generator which will yield the given values.
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

// FilterFn returns a func which filters 'g' using a given filter func 'rcv',
// and returns the resulting generator.
// Example:
//
//	g := FilterFn(New(1, 2, 3))(
//		func(v int) bool {
//			return v != 2
//		},
//	)
//
//	t.Log(g())	// 1, true
//	t.Log(g())	// 3, true
//	t.Log(g())	// 0, false
func FilterFn[T any](g Gen[T]) func(rcv func(T) bool) Gen[T] {
	return func(f func(v T) bool) Gen[T] {
		if g == nil {
			return New[T]()
		}
		if f == nil {
			return g
		}

		return func() (v T, cont bool) {
			v, cont = g()
			for ; cont && !f(v); v, cont = g() {
			}

			return v, cont
		}
	}
}

// MapFn returns a func which maps 'g' using a given mapper func 'rcv', and
// returns the resulting generator.
// Example:
//
//	g := MapFn[int, int](New(1, 2))(
//		func(v int) int {
//			return v + 1
//		},
//	)
//
//	t.Log(g())	// 2, true
//	t.Log(g())	// 3, true
//	t.Log(g())	// 0, false
func MapFn[T, U any](g Gen[T]) func(rcv func(T) U) Gen[U] {
	return func(f func(T) U) Gen[U] {
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
}

// ReduceFn returns a func which reduces 'g' using a given reducer func, and
// returns the resulting value.
// Example:
//
//	r := ReduceFn(New(1, 2, 3))(
//		func(accumulator, current int) int {
//			return accumulator + current
//		},
//	)
//
//	t.Log(r)	// 6.
func ReduceFn[T any](g Gen[T]) func(rcv func(T, T) T) T {
	return func(f func(acc T, curr T) T) (r T) {
		if g == nil || f == nil {
			return r
		}

		for v, cont := g(); cont; v, cont = g() {
			r = f(r, v)
		}

		return r
	}
}

// IntoSlice reads all values of 'g' and returns them in a slice with a small
// initial capacity. The optional 'size' may be specified for a specific cap.
func IntoSlice[T any](g Gen[T], size ...int) []T {
	if g == nil {
		return []T{}
	}

	c := 8
	if len(size) > 0 {
		c = 0
		for _, v := range size {
			c += v
		}
	}

	r := make([]T, 0, c)
	for v, cont := g(); cont; v, cont = g() {
		r = append(r, v)
	}

	return r
}

// IntoMapKFn returns a func which creates a map[K]V using the given generator.
// It does so by reading all values of 'g' and storing them as keys. Vals are
// defined using the given func 'f'.
// Example:
//
//	m := IntoMapKFn[int, int](New(1, 2, 3))(
//		func(key int) (val int) {
//			val = key + 1
//			return val
//		},
//	)
//
//	t.Log(m)	// map[1:2 2:3 3:4]
func IntoMapKFn[K comparable, V any](g Gen[K]) func(func(K) V) map[K]V {
	return func(f func(K) V) map[K]V {
		if g == nil || f == nil {
			return map[K]V{}
		}

		r := make(map[K]V)
		for k, cont := g(); cont; k, cont = g() {
			r[k] = f(k)
		}

		return r
	}
}

// IntoMapVFn returns a func which creates a map[K]V using the given generator.
// It does so by reading all values of 'g' and storing them as values for keys
// that are defined using the given func 'f'.
// Example:
//
//	m := IntoMapVFn[int](New(1, 2, 3))(
//		func(val int) (key int) {
//			key = val - 1
//			return key
//		},
//	)
//
//	t.Log(m)	// map[0:1 1:2 2:3]
func IntoMapVFn[K comparable, V any](g Gen[V]) func(func(V) K) map[K]V {
	return func(f func(V) K) map[K]V {
		if g == nil || f == nil {
			return map[K]V{}
		}

		r := make(map[K]V)
		for v, cont := g(); cont; v, cont = g() {
			r[f(v)] = v
		}

		return r
	}
}

// IntoChan returns a chan which reads from the given generator. Values are fed
// into the chan from a new goroutine so any use of 'g' should stop after this call.
func IntoChan[T any](g Gen[T]) <-chan T {
	ch := make(chan T)
	if g == nil {
		close(ch)
		return ch
	}

	go func() {
		defer close(ch)
		for v, cont := g(); cont; v, cont = g() {
			ch <- v
		}
	}()

	return ch
}
