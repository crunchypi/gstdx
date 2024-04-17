package mapx

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

func TestNewIdeal(t *testing.T) {
	p1 := Pair[int, int]{1, 2}
	p2 := Pair[int, int]{2, 3}

	assertEq("r", map[int]int{1: 2, 2: 3}, New(p1, p2), func(s string) { t.Fatal(s) })
}

func TestFilterFnIdeal(t *testing.T) {
	init := map[int]int{1: 2, 2: 3}
	have := FilterFn(init)(func(p Pair[int, int]) bool { return p.K > 1 })
	want := map[int]int{2: 3}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestFilterFnWithNilM(t *testing.T) {
	init := *new(map[int]int)
	have := FilterFn(init)(func(p Pair[int, int]) bool { return p.K > 1 })
	want := map[int]int{}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestFilterFnWithNilF(t *testing.T) {
	init := map[int]int{1: 2, 2: 3}
	have := FilterFn(init)(nil)
	want := map[int]int{1: 2, 2: 3}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestFilterKFnIdeal(t *testing.T) {
	init := map[int]int{1: 2, 2: 3}
	have := FilterKFn(init)(func(k int) bool { return k > 1 })
	want := map[int]int{2: 3}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestFilterKFnWithNilM(t *testing.T) {
	init := *new(map[int]int)
	have := FilterKFn(init)(func(k int) bool { return k > 1 })
	want := map[int]int{}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestFilterKFnWithNilF(t *testing.T) {
	init := map[int]int{1: 2, 2: 3}
	have := FilterKFn(init)(nil)
	want := map[int]int{1: 2, 2: 3}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestFilterVFnIdeal(t *testing.T) {
	init := map[int]int{1: 2, 2: 3}
	have := FilterVFn(init)(func(v int) bool { return v > 2 })
	want := map[int]int{2: 3}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestFilterVFnWithNilM(t *testing.T) {
	init := *new(map[int]int)
	have := FilterVFn(init)(func(v int) bool { return v > 2 })
	want := map[int]int{}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestFilterVFnWithNilF(t *testing.T) {
	init := map[int]int{1: 2, 2: 3}
	have := FilterVFn(init)(nil)
	want := map[int]int{1: 2, 2: 3}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestMapFnIdeal(t *testing.T) {
	init := map[int]int{1: 2, 2: 3}
	conv := func(p Pair[int, int]) Pair[int, int] { p.K++; p.V++; return p }
	have := MapFn[int, int, int, int](init)(conv)
	want := map[int]int{2: 3, 3: 4}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestMapFnWithNilM(t *testing.T) {
	init := *new(map[int]int)
	conv := func(p Pair[int, int]) Pair[int, int] { p.K++; p.V++; return p }
	have := MapFn[int, int, int, int](init)(conv)
	want := map[int]int{}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestMapFnWithNilF(t *testing.T) {
	init := map[int]int{1: 2, 2: 3}
	have := MapFn[int, int, int, int](init)(nil)
	want := map[int]int{}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestMapKFnIdeal(t *testing.T) {
	init := map[int]int{1: 2, 2: 3}
	have := MapKFn[int, int](init)(func(k int) int { return k + 1 })
	want := map[int]int{2: 2, 3: 3}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestMapKFnWithNilM(t *testing.T) {
	init := *new(map[int]int)
	have := MapKFn[int, int](init)(func(k int) int { return k + 1 })
	want := map[int]int{}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestMapKFnWithNilF(t *testing.T) {
	init := map[int]int{1: 2, 2: 3}
	have := MapKFn[int, int](init)(nil)
	want := map[int]int{}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestMapVFnIdeal(t *testing.T) {
	init := map[int]int{1: 1, 2: 2}
	have := MapVFn[int, int, int](init)(func(v int) int { return v + 1 })
	want := map[int]int{1: 2, 2: 3}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestMapVFnWithNilM(t *testing.T) {
	init := *new(map[int]int)
	have := MapVFn[int, int, int](init)(func(v int) int { return v + 1 })
	want := map[int]int{}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestMapVFnWithNilF(t *testing.T) {
	init := map[int]int{1: 1, 2: 2}
	have := MapVFn[int, int, int](init)(nil)
	want := map[int]int{}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}
