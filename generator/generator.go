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
