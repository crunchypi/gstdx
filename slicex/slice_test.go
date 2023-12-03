package slicex

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"testing"
	"time"
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

func TestMapFnIdeal(t *testing.T) {
	s1 := New(1, 2, 3)
	s2 := MapFn[int, int](s1)(
		func(v int) int {
			return v + 1
		},
	)

	assertEq("slice", []int{2, 3, 4}, s2, func(s string) { t.Fatal(s) })
}

func TestMapFnWithNilS(t *testing.T) {
	s1 := *new([]int)
	s2 := MapFn[int, int](s1)(func(int) int { return 0 })

	assertEq("slice", []int{}, s2, func(s string) { t.Fatal(s) })
}

func TestMapFnWithNilF(t *testing.T) {
	s1 := New(1, 2, 3)
	s2 := MapFn[int, int](s1)(nil)

	assertEq("slice", []int{}, s2, func(s string) { t.Fatal(s) })
}

func TestReduceFnIdeal(t *testing.T) {
	want := 6
	have := ReduceFn(New(1, 2, 3))(
		func(accumulate, current int) int {
			return accumulate + current
		},
	)

	assertEq("result", want, have, func(s string) { t.Fatal(s) })
}

func TestReduceFnWithNilS(t *testing.T) {
	want := 0
	have := ReduceFn[int, []int](nil)(
		func(int, int) int {
			return 0
		},
	)

	assertEq("slice", want, have, func(s string) { t.Fatal(s) })
}

func TestReduceFnWithNilF(t *testing.T) {
	want := 0
	have := ReduceFn[int](New(1, 2, 3))(nil)

	assertEq("slice", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoCloneIdeal(t *testing.T) {
	s1 := New(1, 2, 3)
	s2 := IntoClone(s1)

	s1[0] = 0

	assertEq("slice", New(1, 2, 3), s2, func(s string) { t.Fatal(s) })
}

func TestIntoCloneWithNilS(t *testing.T) {
	s1 := *new([]int)
	s2 := IntoClone(s1)

	assertEq("slice", []int{}, s2, func(s string) { t.Fatal(s) })
}

func TestIntoMapKFnIdeal(t *testing.T) {
	s1 := New(1, 2)
	sk := New[int]()
	sv := New[int]()

	m := IntoMapKFn[int, int](s1)(func(v int) int { return v + 1 })
	for k, v := range m {
		sk = append(sk, k)
		sv = append(sv, v)
	}

	sort.Ints(sk)
	sort.Ints(sv)

	assertEq("sk", s1, sk, func(s string) { t.Fatal(s) })
	assertEq("sv", New(2, 3), sv, func(s string) { t.Fatal(s) })
}

func TestIntoMapKFnWithNilS(t *testing.T) {
	m := IntoMapKFn[int, int, []int](nil)(func(int) int { return 0 })
	assertEq("len", 0, len(m), func(s string) { t.Fatal(s) })
}

func TestIntoMapKFnWithNilF(t *testing.T) {
	s1 := New(1, 2)
	sk := New[int]()
	sv := New[int]()

	m := IntoMapKFn[int, int](s1)(nil)
	for k, v := range m {
		sk = append(sk, k)
		sv = append(sv, v)
	}

	sort.Ints(sk)
	sort.Ints(sv)

	assertEq("sk", s1, sk, func(s string) { t.Fatal(s) })
	assertEq("sv", New(0, 0), sv, func(s string) { t.Fatal(s) })
}

func TestIntoMapVFnIdeal(t *testing.T) {
	s1 := New(1, 2)
	sk := New[int]()
	sv := New[int]()

	m := IntoMapVFn[int, int](s1)(func(v int) int { return v + 1 })
	for k, v := range m {
		sk = append(sk, k)
		sv = append(sv, v)
	}

	sort.Ints(sk)
	sort.Ints(sv)

	assertEq("sv", s1, sv, func(s string) { t.Fatal(s) })
	assertEq("sk", New(2, 3), sk, func(s string) { t.Fatal(s) })
}

func TestIntoMapVFnWithNilS(t *testing.T) {
	m := IntoMapVFn[int, int, []int](nil)(func(int) int { return 0 })
	assertEq("len", 0, len(m), func(s string) { t.Fatal(s) })
}

func TestIntoMapVFnWithNilF(t *testing.T) {
	s1 := New(1, 2)
	sk := New[int]()
	sv := New[int]()

	m := IntoMapVFn[int, int](s1)(nil)
	for k, v := range m {
		sk = append(sk, k)
		sv = append(sv, v)
	}

	sort.Ints(sk)
	sort.Ints(sv)

	assertEq("sv", New(2), sv, func(s string) { t.Fatal(s) })
	assertEq("sk", New(0), sk, func(s string) { t.Fatal(s) })
}

func TestIntoChanIdeal(t *testing.T) {
	want := New(1, 2, 3)
	have := New[int]()

	ctx, ctxCancel := context.WithCancel(context.Background())
	go func() {
		defer ctxCancel()
		for v := range IntoChan(want) {
			have = append(have, v)
		}
	}()

	select {
	case <-time.After(time.Second * 2):
		t.Fatal("test hung")
	case <-ctx.Done():
	}

	assertEq("slice", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoChanWithNilS(t *testing.T) {
	want := New[int]()
	have := []int{}

	ctx, ctxCancel := context.WithCancel(context.Background())
	go func() {
		defer ctxCancel()
		for v := range IntoChan[int, []int](nil) {
			have = append(have, v)
		}
	}()

	select {
	case <-time.After(time.Second * 2):
		t.Fatal("test hung")
	case <-ctx.Done():
	}

	assertEq("slice", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoGeneratorIdeal(t *testing.T) {
	want := New(1, 2, 3)
	have := New[int]()

	gen := IntoGenerator(want)
	for v, ok := gen(); ok; v, ok = gen() {
		have = append(have, v)
	}

	assertEq("slice", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoGeneratorWithNilS(t *testing.T) {
	_, ok := IntoGenerator[int, []int](nil)()
	assertEq("bool", false, ok, func(s string) { t.Fatal(s) })
}
