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

func filter[T any, S ~[]T](s S, f func(T) bool) S {
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
		return filter(s, f)
	}
}
