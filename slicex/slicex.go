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
