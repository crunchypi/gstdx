package slicex

// New returns a slice containing a shallow copy of the given values.
func New[T any](vs ...T) (r []T) {
	r = make([]T, len(vs))

	if len(vs) == 0 {
		return
	}

	copy(r, vs)
	return
}

// FilterFn returns a func which filters 's' using a given filter func.
// Example:
//
//	s1 := []int{1,2,3,4}
//	s2 := FilterFn(s1)(
//		func(v int) bool {
//			return v%2 == 0
//		},
//	)
//
//	// s2 is []int{2,4}
func FilterFn[T any, S ~[]T](s S) func(filter func(T) bool) S {
	return func(f func(T) bool) S {
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
}

// MapFn returns a func which maps 's' using a given mapper func.
// Example:
//
//	s1 := New(1, 2, 3)
//	s2 := MapFn[int, int](s1)(
//		func(v int) int {
//			return v + 1
//		},
//	)
//
//	// s2 is []int{2, 3, 4}
func MapFn[T, U any, S ~[]T](s S) func(mapper func(T) U) []U {
	return func(f func(T) U) []U {
		if len(s) == 0 || f == nil {
			return []U{}
		}

		r := make([]U, 0, len(s))
		for _, v := range s {
			r = append(r, f(v))
		}

		return r
	}
}

// ReduceFn returns a func which reduces 's' into T using a given reducer.
// Example:
//
//	want := 6
//	have := ReduceFn(New(1, 2, 3))(
//		func(accumulate, current int) int {
//			return accumulate + current
//		},
//	)
//
//	// want == have is true.
func ReduceFn[T any, S ~[]T](s S) func(reducer func(acc, curr T) T) T {
	return func(f func(acc T, curr T) T) (r T) {
		if len(s) == 0 || f == nil {
			return r
		}

		for _, v := range s {
			r = f(r, v)
		}

		return r
	}
}

// IntoClone returns a slice containing a shallow copy of the given slice.
func IntoClone[T any, S ~[]T](s S) S {
	r := make([]T, 0, len(s))
	for _, v := range s {
		r = append(r, v)
	}

	return r
}

// IntoMapKFn returns a func which makes a map[K]V using the given slice "s".
// Keys in the map are elements in "s", and values are generated using the
// mapper func 'f', which maps the keys to vals.
// Example:
//
//	m := IntoMapKFn[int, int]([]int{1, 2, 3})(
//		func(k int) (v int) {
//			v = k + 1
//			return
//		},
//	)
//
//	// m is map[1:2 2:3 3:4]
func IntoMapKFn[K comparable, V any, S ~[]K](s S) func(f func(K) V) map[K]V {
	return func(f func(K) V) map[K]V {
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
}

// IntoMapVFn returns a func which makes a map[K]V using the given slice "s".
// Vals in the map are elements in "s", and their associated keys are generated
// using the mapper func "f".
// Example:
//
//	m := IntoMapVFn[int, int]([]int{1, 2, 3})(
//		func(v int) (k int) {
//			k = v + 1
//			return
//		},
//	)
//
//	// m is map[2:1 3:2 4:3]
func IntoMapVFn[K comparable, V any, S ~[]V](s S) func(f func(V) K) map[K]V {
	return func(f func(V) K) map[K]V {
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
}

// IntoChan returns a chan which is fed the contents of "s" from a new goroutine.
// There is no internal copying of "s", so its use should stop after calling this
// func. On empty or nil "s", the chan is returned pre-closed.
// Example:
//
//	s := make([]int, 0, 3)
//	for v := range IntoChan([]int{1, 2, 3}) {
//		s = append(s, v)
//	}
//
//	// s is []int{1, 2, 3}
func IntoChan[T any, S ~[]T](s S) <-chan T {
	if len(s) == 0 {
		r := make(chan T)
		close(r)
		return r
	}

	r := make(chan T)
	go func() {
		defer close(r)

		for _, v := range s {
			r <- v
		}
	}()

	return r
}

// IntoGenerator returns a func which, when called, yields values from 's'.
// If the bool is false, then the iteration is completed and the value is a
// default T. There is no copying of 's', so there should be no further use
// of 's' when it is given to this func.
// Example:
//
//	s := make([]int, 0, 3)
//	g := IntoGenerator([]int{1, 2, 3})
//
//	for v, ok := g(); ok; v, ok = g() {
//		s = append(s, v)
//	}
//
//	// s is []int{1, 2, 3}
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
