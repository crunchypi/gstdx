package generator

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

func sliceFromGenerator[T any](g func() (T, bool)) []T {
	s := make([]T, 0, 8)
	for v, ok := g(); ok; v, ok = g() {
		s = append(s, v)
	}

	return s
}

func TestNewWithVals(t *testing.T) {
	have := sliceFromGenerator(New(1, 2, 3))
	want := []int{1, 2, 3}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestNewWithNone(t *testing.T) {
	have := sliceFromGenerator(New[int]())
	want := []int{}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestFilterFnIdeal(t *testing.T) {
	init := New(1, 2, 3)
	have := sliceFromGenerator(FilterFn(init)(func(v int) bool { return v != 2 }))
	want := []int{1, 3}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestFilterFnWithNilG(t *testing.T) {
	init := *new(Gen[int])
	have := sliceFromGenerator(FilterFn(init)(func(v int) bool { return v != 2 }))
	want := []int{}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestFilterFnWithNilF(t *testing.T) {
	init := New(1, 2, 3)
	have := sliceFromGenerator(FilterFn(init)(nil))
	want := []int{1, 2, 3}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestMapFnIdeal(t *testing.T) {
	init := New(1, 2, 3)
	have := sliceFromGenerator(MapFn[int, int](init)(func(v int) int { return v + 1 }))
	want := []int{2, 3, 4}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestMapFnWithNilG(t *testing.T) {
	init := *new(Gen[int])
	have := sliceFromGenerator(MapFn[int, int](init)(func(v int) int { return v + 1 }))
	want := []int{}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestMapFnWithNilF(t *testing.T) {
	init := New(1, 2, 3)
	have := sliceFromGenerator(MapFn[int, int](init)(nil))
	want := []int{}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}
