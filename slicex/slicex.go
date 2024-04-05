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
