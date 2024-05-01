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
