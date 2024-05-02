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
