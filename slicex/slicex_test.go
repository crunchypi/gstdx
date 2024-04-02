package slicex

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

func TestNewWithVals(t *testing.T) {
	want := []int{1, 2, 3}
	have := New(1, 2, 3)

	assertEq("slice", want, have, func(s string) { t.Fatal(s) })
}

func TestNewWithNone(t *testing.T) {
	want := []int{}
	have := New[int]()

	assertEq("slice", want, have, func(s string) { t.Fatal(s) })
}
