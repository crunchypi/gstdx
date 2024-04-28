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

func sliceFromChan[T any](c <-chan T) []T {
	s := make([]T, 0, 8)
	for v := range c {
		s = append(s, v)
	}

	return s
}

func sliceFromGenerator[T any](g func() (T, bool)) []T {
	s := make([]T, 0, 8)
	for v, ok := g(); ok; v, ok = g() {
		s = append(s, v)
	}

	return s
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

func TestReduceFnIdeal(t *testing.T) {
	type P = Pair[int, int]

	init := map[int]int{1: 2, 2: 3}
	have := ReduceFn(init)(func(a, c P) P { a.K += c.K; a.V += c.V; return a })
	want := P{3, 5}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestReduceFnWithNilM(t *testing.T) {
	type P = Pair[int, int]

	init := *new(map[int]int)
	have := ReduceFn(init)(func(a, c P) P { a.K += c.K; a.V += c.V; return a })
	want := P{}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestReduceFnWithNilF(t *testing.T) {
	type P = Pair[int, int]

	init := map[int]int{1: 2, 2: 3}
	have := ReduceFn(init)(nil)
	want := P{}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestReduceKFnIdeal(t *testing.T) {
	init := map[int]int{1: 2, 2: 3}
	have := ReduceKFn(init)(func(c, a int) int { return c + a })
	want := 3

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestReduceKFnWithNilM(t *testing.T) {
	init := *new(map[int]int)
	have := ReduceKFn(init)(func(c, a int) int { return c + a })
	want := 0

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestReduceKFnWithNilF(t *testing.T) {
	init := map[int]int{1: 2, 2: 3}
	have := ReduceKFn(init)(nil)
	want := 0

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestReduceVFnIdeal(t *testing.T) {
	init := map[int]int{1: 2, 2: 3}
	have := ReduceVFn(init)(func(c, a int) int { return c + a })
	want := 5

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestReduceVFnWithNilM(t *testing.T) {
	init := *new(map[int]int)
	have := ReduceVFn(init)(func(c, a int) int { return c + a })
	want := 0

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestReduceVFnWithNilF(t *testing.T) {
	init := map[int]int{1: 2, 2: 3}
	have := ReduceVFn(init)(nil)
	want := 0

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoCloneIdeal(t *testing.T) {
	init := map[int]int{1: 1, 2: 2}
	have := IntoClone(init)
	want := map[int]int{1: 1, 2: 2}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoCloneWithNilM(t *testing.T) {
	init := *new(map[int]int)
	have := IntoClone(init)
	want := map[int]int{}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoSliceIdeal(t *testing.T) {
	have := IntoSlice(map[int]int{1: 1, 2: 2})
	want := []Pair[int, int]{{K: 1, V: 1}, {K: 2, V: 2}}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoSliceWithNilM(t *testing.T) {
	have := IntoSlice(*new(map[int]int))
	want := []Pair[int, int]{}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoSliceKIdeal(t *testing.T) {
	have := IntoSliceK(map[int]int{1: 1, 2: 2})
	want := []int{1, 2}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoSliceKWithNilM(t *testing.T) {
	have := IntoSliceK(*new(map[int]int))
	want := []int{}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoSliceVIdeal(t *testing.T) {
	have := IntoSliceV(map[int]int{1: 1, 2: 2})
	want := []int{1, 2}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoSliceVWithNilM(t *testing.T) {
	have := IntoSliceV(*new(map[int]int))
	want := []int{}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoChanIdeal(t *testing.T) {
	init := map[int]int{1: 2, 2: 3}
	have := sliceFromChan(IntoChan(init))
	want := []Pair[int, int]{{K: 1, V: 2}, {K: 2, V: 3}}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoChanWithNilM(t *testing.T) {
	init := *new(map[int]int)
	have := sliceFromChan(IntoChan(init))
	want := []Pair[int, int]{}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoChanKIdeal(t *testing.T) {
	init := map[int]int{1: 2, 2: 3}
	have := sliceFromChan(IntoChanK(init))
	want := []int{1, 2}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoChanKWithNilM(t *testing.T) {
	init := *new(map[int]int)
	have := sliceFromChan(IntoChanK(init))
	want := []int{}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoChanVIdeal(t *testing.T) {
	init := map[int]int{1: 2, 2: 3}
	have := sliceFromChan(IntoChanV(init))
	want := []int{2, 3}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoChanVWithNilM(t *testing.T) {
	init := *new(map[int]int)
	have := sliceFromChan(IntoChanV(init))
	want := []int{}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoGeneratorIdeal(t *testing.T) {
	init := map[int]int{1: 2, 2: 3}
	have := sliceFromGenerator(IntoGenerator(init))
	want := []Pair[int, int]{{K: 1, V: 2}, {K: 2, V: 3}}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoGeneratorWithNilM(t *testing.T) {
	init := *new(map[int]int)
	have := sliceFromGenerator(IntoGenerator(init))
	want := []Pair[int, int]{}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}
