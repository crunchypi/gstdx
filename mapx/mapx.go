package mapx

type Key = comparable // Abbreviation, used a lot.
type Val = any        // Abbreviation, used a lot.

type Pair[K Key, V Val] struct {
	K K
	V V
}

// New returns a map[K]V from the given pairs.
func New[K Key, V Val](kv ...Pair[K, V]) map[K]V {
	r := make(map[K]V, len(kv))
	for _, pair := range kv {
		r[pair.K] = pair.V

	}

	return r
}
