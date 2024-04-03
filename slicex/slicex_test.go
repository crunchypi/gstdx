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

func TestFilterFnIdeal(t *testing.T) {
	s1 := New(1, 2, 3, 4)
	s2 := FilterFn(s1)(
		func(v int) bool {
			return v%2 == 0
		},
	)

	assertEq("slice", []int{2, 4}, s2, func(s string) { t.Fatal(s) })
}

func TestFilterFnWithNilS(t *testing.T) {
	s1 := *new([]int)
	s2 := FilterFn[int](s1)(func(int) bool { return false })

	assertEq("slice", []int{}, s2, func(s string) { t.Fatal(s) })
}

func TestFilterFnWithNilF(t *testing.T) {
	s1 := New(1, 2, 3)
	s2 := FilterFn[int](s1)(nil)

	assertEq("slice", s1, s2, func(s string) { t.Fatal(s) })
}
