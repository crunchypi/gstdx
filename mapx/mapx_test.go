package mapx

import (
	"encoding/json"
	"fmt"
	"slices"
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

func TestFilterKFnIdeal(t *testing.T) {
	init := map[int]int{1: 2, 2: 3}
	have := FilterKFn(init)(func(k int) bool { return k > 1 })
	want := map[int]int{2: 3}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestFilterVFnIdeal(t *testing.T) {
	init := map[int]int{1: 2, 2: 3}
	have := FilterVFn(init)(func(v int) bool { return v > 2 })
	want := map[int]int{2: 3}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestMapFnIdeal(t *testing.T) {
	init := map[int]int{1: 2, 2: 3}
	conv := func(p Pair[int, int]) Pair[int, int] { p.K++; p.V++; return p }
	have := MapFn[int, int, int, int](init)(conv)
	want := map[int]int{2: 3, 3: 4}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestMapKFnIdeal(t *testing.T) {
	init := map[int]int{1: 2, 2: 3}
	have := MapKFn[int, string](init)(func(k int) string { return fmt.Sprint(k) })
	want := map[string]int{"1": 2, "2": 3}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestMapVFnIdeal(t *testing.T) {
	init := map[int]int{1: 2, 2: 3}
	have := MapVFn[int, int, string](init)(func(k int) string { return fmt.Sprint(k) })
	want := map[int]string{1: "2", 2: "3"}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestReduceFnIdeal(t *testing.T) {
	type P = Pair[int, int]

	init := map[int]int{1: 2, 2: 3}
	have := ReduceFn(init)(func(a, c P) P { a.K += c.K; a.V += c.V; return a })
	want := P{3, 5}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestReduceKFnIdeal(t *testing.T) {
	init := map[int]int{1: 2, 2: 3}
	have := ReduceKFn(init)(func(c, a int) int { return c + a })
	want := 3

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestReduceVFnIdeal(t *testing.T) {
	init := map[int]int{1: 2, 2: 3}
	have := ReduceVFn(init)(func(c, a int) int { return c + a })
	want := 5

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoCloneFnIdeal(t *testing.T) {
	init := map[int]int{1: 2, 2: 3}
	have := IntoClone(init)
	want := map[int]int{1: 2, 2: 3}

	init[1] = 3
	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoPairsIdeal(t *testing.T) {
	init := map[int]int{1: 2, 2: 3}
	have := IntoPairs(init)
	want := []Pair[int, int]{{K: 1, V: 2}, {K: 2, V: 3}}

	slices.SortFunc(have, func(a Pair[int, int], b Pair[int, int]) int { return a.K - b.K })
	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoKeysIdeal(t *testing.T) {
	init := map[int]int{1: 2, 2: 3}
	have := IntoKeys(init)
	want := []int{1, 2}

	slices.Sort(have)
	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoValsIdeal(t *testing.T) {
	init := map[int]int{1: 2, 2: 3}
	have := IntoVals(init)
	want := []int{2, 3}

	slices.Sort(have)
	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoChanIdeal(t *testing.T) {
	want := map[int]int{1: 2, 2: 3}
	have := make(map[int]int, len(want))

	for pair := range IntoChan(want) {
		have[pair.K] = pair.V
	}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoChanKIdeal(t *testing.T) {
	init := map[int]int{1: 2, 2: 3}
	want := []int{1, 2}
	have := make([]int, 0, len(want))

	for k := range IntoChanK(init) {
		have = append(have, k)
	}

	slices.Sort(have)
	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoChanVIdeal(t *testing.T) {
	init := map[int]int{1: 2, 2: 3}
	want := []int{2, 3}
	have := make([]int, 0, len(want))

	for v := range IntoChanV(init) {
		have = append(have, v)
	}

	slices.Sort(have)
	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoGeneratorIdeal(t *testing.T) {
	want := map[int]int{1: 2, 2: 3}
	have := make(map[int]int, len(want))

	g := IntoGenerator(want)
	for p, ok := g(); ok; p, ok = g() {
		have[p.K] = p.V
	}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoGeneratorKIdeal(t *testing.T) {
	init := map[int]int{1: 2, 2: 3}
	want := []int{1, 2}
	have := make([]int, 0, len(want))

	g := IntoGeneratorK(init)
	for k, ok := g(); ok; k, ok = g() {
		have = append(have, k)
	}

	slices.Sort(have)
	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoGeneratorVIdeal(t *testing.T) {
	init := map[int]int{1: 2, 2: 3}
	want := []int{2, 3}
	have := make([]int, 0, len(want))

	g := IntoGeneratorV(init)
	for v, ok := g(); ok; v, ok = g() {
		have = append(have, v)
	}

	slices.Sort(have)
	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}
