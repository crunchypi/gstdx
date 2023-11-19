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
	s := FilterFn[int](nil)(func(int) bool { return false })
	assertEq("slice", []int{}, s, func(s string) { t.Fatal(s) })
}

func TestFilterFnWithNilF(t *testing.T) {
	s := FilterFn[int](New(1, 2, 3))(nil)
	assertEq("slice", []int{}, s, func(s string) { t.Fatal(s) })
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
	s := MapFn[int, int](nil)(func(int) int { return 0 })
	assertEq("slice", []int{}, s, func(s string) { t.Fatal(s) })
}

func TestMapFnWithNilF(t *testing.T) {
	s := MapFn[int, int](New(1, 2, 3))(nil)
	assertEq("slice", []int{}, s, func(s string) { t.Fatal(s) })
}

func TestReduceFnIdeal(t *testing.T) {
	s := New(1, 2, 3)
	r := ReduceFn(s)(
		func(accumulate, current int) int {
			return accumulate + current
		},
	)

	assertEq("result", 6, r, func(s string) { t.Fatal(s) })
}

func TestReduceFnWithNilS(t *testing.T) {
	r := ReduceFn[int](nil)(func(int, int) int { return 0 })
	assertEq("slice", 0, r, func(s string) { t.Fatal(s) })
}

func TestReduceFnWithNilF(t *testing.T) {
	r := ReduceFn[int](New(1, 2, 3))(nil)
	assertEq("slice", 0, r, func(s string) { t.Fatal(s) })
}

func TestIntoCloneIdeal(t *testing.T) {
	s1 := New(1, 2, 3)
	s2 := IntoClone(s1)

	s1[0] = 0

	assertEq("slice", New(1, 2, 3), s2, func(s string) { t.Fatal(s) })
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
	m := IntoMapKFn[int, int](nil)(func(int) int { return 0 })
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
	m := IntoMapVFn[int, int](nil)(func(int) int { return 0 })
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
	s1 := New(1, 2, 3)
	s2 := New[int]()

	ctx, ctxCancel := context.WithCancel(context.Background())
	go func() {
		defer ctxCancel()
		for v := range IntoChan(context.Background(), s1) {
			s2 = append(s2, v)
		}
	}()

	select {
	case <-time.After(time.Second * 2):
		t.Fatal("test hung")
	case <-ctx.Done():
	}

	assertEq("slice", s1, s2, func(s string) { t.Fatal(s) })
}

func TestIntoChanWithCancel(t *testing.T) {
	s1 := New(1, 2, 3)
	s2 := New[int]()

	ctx := context.Background()

	ctx1, ctxCancel1 := context.WithCancel(ctx) // For test hung.
	ctx2, ctxCancel2 := context.WithCancel(ctx) // For abort.
	go func() {
		defer ctxCancel1()
		for v := range IntoChan(ctx2, s1) {
			s2 = append(s2, v)
			ctxCancel2()
		}
	}()

	select {
	case <-time.After(time.Second * 2):
		t.Fatal("test hung")
	case <-ctx1.Done():
	}

	assertEq("slice", New(1), s2, func(s string) { t.Fatal(s) })
}

func TestIntoChanWithNilC(t *testing.T) {
	s1 := New(1, 2, 3)
	s2 := New[int]()

	ctx, ctxCancel := context.WithCancel(context.Background())
	go func() {
		defer ctxCancel()
		for v := range IntoChan(nil, s1) {
			s2 = append(s2, v)
		}
	}()

	select {
	case <-time.After(time.Second * 2):
		t.Fatal("test hung")
	case <-ctx.Done():
	}

	assertEq("slice", s1, s2, func(s string) { t.Fatal(s) })
}

func TestIntoChanWithNilS(t *testing.T) {
	s := New[int]()

	ctx, ctxCancel := context.WithCancel(context.Background())
	go func() {
		defer ctxCancel()
		for v := range IntoChan[int](context.Background(), nil) {
			s = append(s, v)
		}
	}()

	select {
	case <-time.After(time.Second * 2):
		t.Fatal("test hung")
	case <-ctx.Done():
	}

	assertEq("slice", New[int](), s, func(s string) { t.Fatal(s) })
}

func TestIntoGeneratorIdeal(t *testing.T) {
	s1 := New(1, 2, 3)
	s2 := New[int]()

	gen := IntoGenerator(s1)
	for v, ok := gen(); ok; v, ok = gen() {
		s2 = append(s2, v)
	}

	assertEq("slice", s1, s2, func(s string) { t.Fatal(s) })
}

func TestIntoGeneratorWithNilS(t *testing.T) {
	_, ok := IntoGenerator[int](nil)()
	assertEq("bool", false, ok, func(s string) { t.Fatal(s) })
}
