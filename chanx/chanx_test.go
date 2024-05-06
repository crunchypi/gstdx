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

func TestFilterFnIdeal(t *testing.T) {
	init := New(1, 2, 3)
	have := sliceFromChan(FilterFn(init)(func(v int) bool { return v%2 != 0 }))
	want := []int{1, 3}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestFilterFnWithNilC(t *testing.T) {
	init := *new(chan int)
	have := sliceFromChan(FilterFn(init)(func(v int) bool { return true }))
	want := []int{}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestFilterFnWithNilF(t *testing.T) {
	init := New(1, 2, 3)
	have := sliceFromChan(FilterFn(init)(nil))
	want := []int{1, 2, 3}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestMapFnIdeal(t *testing.T) {
	init := New(1, 2, 3)
	have := sliceFromChan(MapFn[int, int](init)(func(v int) int { return v + 1 }))
	want := []int{2, 3, 4}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestMapFnWithNilC(t *testing.T) {
	init := *new(chan int)
	have := sliceFromChan(MapFn[int, int](init)(func(v int) int { return v + 1 }))
	want := []int{}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestMapFnWithNilF(t *testing.T) {
	init := New(1, 2, 3)
	have := sliceFromChan(MapFn[int, int](init)(nil))
	want := []int{}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestReduceFnIdeal(t *testing.T) {
	want := 6
	have := ReduceFn(New(1, 2, 3))(
		func(acc, curr int) int {
			return acc + curr
		},
	)

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestReduceFnWithNilC(t *testing.T) {
	want := 0
	have := ReduceFn(*new(chan int))(
		func(acc, curr int) int {
			return acc + curr
		},
	)

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestReduceFnWithNilF(t *testing.T) {
	want := 0
	have := ReduceFn(New(1, 2, 3))(nil)

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoSliceIdeal(t *testing.T) {
	want := []int{1, 2, 3}
	have := IntoSlice(New(want...), 8)

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoSliceWithNilC(t *testing.T) {
	want := []int{}
	have := IntoSlice(*new(chan int))

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoMapKFnIdeal(t *testing.T) {
	init := New(1, 2, 3)
	have := IntoMapKFn[int, int](init)(func(k int) int { return k })
	want := map[int]int{1: 1, 2: 2, 3: 3}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoMapKFnWithNilC(t *testing.T) {
	init := *new(chan int)
	have := IntoMapKFn[int, int](init)(func(k int) int { return k })
	want := map[int]int{}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoMapKFnWithNilF(t *testing.T) {
	init := New(1, 2, 3)
	have := IntoMapKFn[int, int](init)(nil)
	want := map[int]int{}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}
