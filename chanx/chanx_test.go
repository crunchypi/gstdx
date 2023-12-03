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

func intoSlice[T any](ch <-chan T) []T {
	if ch == nil {
		return []T{}
	}

	r := make([]T, 0, 8)
	for v := range ch {
		r = append(r, v)
	}

	return r
}

func TestNewIdeal(t *testing.T) {
	want := []int{1, 2, 3}
	have := intoSlice(New(want...))

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestNewWithNilS(t *testing.T) {
	want := []int{}
	have := intoSlice(New[int]())

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestFilterFnIdeal(t *testing.T) {
	ch1 := New(1, 1, 2, 3)
	ch2 := FilterFn(ch1)(func(v int) bool { return v%2 == 0 })

	want := []int{2}
	have := intoSlice(ch2)

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestFilterFnWithNilC(t *testing.T) {
	ch1 := *new(chan int)
	ch2 := FilterFn(ch1)(func(v int) bool { return v%2 == 0 })

	want := []int{}
	have := intoSlice(ch2)

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestFilterFnWithNilF(t *testing.T) {
	ch1 := New(1, 1, 2, 3)
	ch2 := FilterFn(ch1)(nil)

	want := []int{1, 1, 2, 3}
	have := intoSlice(ch2)

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestMapFnIdeal(t *testing.T) {
	ch1 := New(1, 2, 3)
	ch2 := MapFn[int, int](ch1)(func(v int) int { return v + 1 })

	want := []int{2, 3, 4}
	have := intoSlice(ch2)

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestMapFnWithNilC(t *testing.T) {
	ch1 := *new(chan int)
	ch2 := MapFn[int, int](ch1)(func(v int) int { return v + 1 })

	want := []int{}
	have := intoSlice(ch2)

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestMapFnWithNilF(t *testing.T) {
	ch1 := New(1, 2, 3)
	ch2 := MapFn[int, int](ch1)(nil)

	want := []int{}
	have := intoSlice(ch2)

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
	have := ReduceFn[int](nil)(
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
	have := IntoSlice(New(want...), 3)

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoSliceWithNilC(t *testing.T) {
	want := []int{}
	have := IntoSlice[int](nil)

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoMapKFnIdeal(t *testing.T) {
	want := map[int]int{1: 2, 2: 3, 3: 4}
	have := IntoMapKFn[int, int](New(1, 2, 3))(
		func(k int) int {
			return k + 1
		},
	)

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoMapKFnWithNilC(t *testing.T) {
	want := map[int]int{}
	have := IntoMapKFn[int, int](nil)(
		func(k int) int {
			return k + 1
		},
	)

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoMapKFnWithNilF(t *testing.T) {
	want := map[int]int{}
	have := IntoMapKFn[int, int](New(1, 2, 3))(nil)

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoMapVFnIdeal(t *testing.T) {
	want := map[int]int{1: 2, 2: 3, 3: 4}
	have := IntoMapVFn[int, int](New(2, 3, 4))(
		func(k int) int {
			return k - 1
		},
	)

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoMapVFnWithNilC(t *testing.T) {
	want := map[int]int{}
	have := IntoMapVFn[int, int](nil)(
		func(k int) int {
			return k - 1
		},
	)

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoMapVFnWithNilF(t *testing.T) {
	want := map[int]int{}
	have := IntoMapVFn[int, int](New(2, 3, 4))(nil)

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoGenIdeal(t *testing.T) {
	want := []int{1, 2, 3}
	have := make([]int, 0, 3)

	g := IntoGen(New(want...))
	for v, cont := g(); cont; v, cont = g() {
		have = append(have, v)
	}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoGenWithNilC(t *testing.T) {
	want := []int{}
	have := make([]int, 0, 3)

	g := IntoGen[int](nil)
	for v, cont := g(); cont; v, cont = g() {
		have = append(have, v)
	}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}
