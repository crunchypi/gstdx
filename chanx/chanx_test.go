package chanx

import (
	"encoding/json"
	"fmt"
	"testing"
)

func assertEq[T any](subject string, a T, b T, f func(string)) {
	if f == nil {
		return
	}

	ab, _ := json.Marshal(a)
	bb, _ := json.Marshal(b)

	as := string(ab)
	bs := string(bb)

	if as == bs {
		return
	}

	s := "unexpected '%v':\n\twant: '%v'\n\thave: '%v'\n"
	f(fmt.Sprintf(s, subject, as, bs))
}

func sliceFromChan[T any](c <-chan T) []T {
	s := make([]T, 0, 8)
	for v := range c {
		s = append(s, v)
	}

	return s
}

func TestNewIdeal(t *testing.T) {
	want := []int{1, 2, 3}
	have := sliceFromChan(New(want...))

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestNewWithNilS(t *testing.T) {
	want := []int{}
	have := sliceFromChan(New[int]())

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}
