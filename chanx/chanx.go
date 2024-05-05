package chanx

// New returns a read-only chan which receives the values given as args here.
// Values are pushed through the chan using a new goroutine.
func New[T any](vs ...T) <-chan T {
	r := make(chan T)

	if len(vs) == 0 {
		close(r)
		return r
	}

	go func() {
		defer close(r)

		for _, v := range vs {
			r <- v
		}
	}()

	return r
}

// FilterFn returns a func which filters 'ch' using a given filter func, and
// returns a chan that reads the remaining values from a new goroutine.
// Example:
//
//	ch := FilterFn(New(1,2,3))(
//		func(v int) bool {
//			return v > 1
//		},
//	)
//
//	// ranging over ch will yield [2,3]
func FilterFn[T any](ch <-chan T) func(func(T) bool) <-chan T {
	return func(f func(T) bool) <-chan T {
		if ch == nil {
			return New[T]()
		}

		if f == nil {
			return ch
		}

		r := make(chan T)
		go func() {
			defer close(r)

			for v := range ch {
				if f(v) {
					r <- v
				}
			}
		}()

		return r
	}
}

// MapFn returns a func which maps 'ch' using the given mapper func, and
// returns a chan that reads the mapped values from a new goroutine.
// Example:
//
//	ch := MapFn[int, int](New(1, 2, 3))(
//		func(v int) int {
//			return v + 1
//		},
//	)
//
//	// ranging over ch will yield [2, 3, 4]
func MapFn[T, U any](ch <-chan T) func(func(T) U) <-chan U {
	return func(f func(T) U) <-chan U {
		if ch == nil || f == nil {
			return New[U]()
		}

		r := make(chan U)
		go func() {
			defer close(r)

			for v := range ch {
				r <- f(v)
			}
		}()

		return r
	}
}

// ReduceFn returns a func which reduces 'ch' using a given reducer func.
// Example:
//
//	r := ReduceFn(New(1, 2, 3))(
//		func(acc, curr int) int {
//			return acc + curr
//		},
//	)
//
//	// r is 6.
func ReduceFn[T any](ch <-chan T) func(func(acc, curr T) T) (r T) {
	return func(f func(acc T, curr T) T) (r T) {
		if ch == nil || f == nil {
			return r
		}

		for v := range ch {
			r = f(r, v)
		}

		return r
	}
}

// IntoSlice reads all values of 'ch' and returns them in a slice with a small
// initial capacity. The optional 'size' may be specified for a specific cap.
func IntoSlice[T any](ch <-chan T, size ...int) []T {
	if ch == nil {
		return []T{}
	}

	l := 8
	if len(size) > 0 {
		l = 0
		for _, v := range size {
			l += v
		}
	}

	r := make([]T, 0, l)
	for v := range ch {
		r = append(r, v)
	}

	return r
}
