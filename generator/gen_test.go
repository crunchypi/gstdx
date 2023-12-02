package generator

import (
	"context"
	"encoding/json"
	"fmt"
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

func intoSlice[T any](g Gen[T]) []T {
	r := make([]T, 0, 10)
	for v, cont := g(); cont; v, cont = g() {
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
	g1 := New(1, 2, 3)
	g2 := FilterFn(g1)(func(v int) bool { return v%2 == 0 })

	want := []int{1, 3}
	have := intoSlice(g2)
	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestFilterFnWithNilG(t *testing.T) {
	g := FilterFn[int](nil)(func(v int) bool { return false })

	want := []int{}
	have := intoSlice(g)

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}
func TestFilterFnWithNilF(t *testing.T) {
	g := FilterFn[int](New(1, 2, 3))(nil)

	want := []int{1, 2, 3}
	have := intoSlice(g)

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestMapFn(t *testing.T) {
	g1 := New(1, 2, 3)
	g2 := MapFn[int, int](g1)(func(v int) int { return v + 1 })

	want := []int{2, 3, 4}
	have := intoSlice(g2)
	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestMapFnWithNilG(t *testing.T) {
	g1 := *new(Gen[int])
	g2 := MapFn[int, int](g1)(func(v int) int { return v + 1 })

	want := []int{}
	have := intoSlice(g2)
	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestMapFnWithNilF(t *testing.T) {
	g1 := New(1, 2, 3)
	g2 := MapFn[int, int](g1)(nil)

	want := []int{}
	have := intoSlice(g2)
	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestReduceFnIdeal(t *testing.T) {
	want := 6
	have := ReduceFn[int](New(1, 2, 3))(
		func(acc, cur int) int {
			return acc + cur
		},
	)

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestReduceFnWithNilG(t *testing.T) {
	want := 0
	have := ReduceFn[int](nil)(
		func(acc, cur int) int {
			return acc + cur
		},
	)

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestReduceFnWithNilF(t *testing.T) {
	want := 0
	have := ReduceFn[int](New(1, 2, 3))(nil)

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoSIdeal(t *testing.T) {
	want := []int{1, 2, 3}
	have := IntoS(New(want...), 3)

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoSWithNilG(t *testing.T) {
	want := []int{}
	have := IntoS[int](nil)

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoMapKFnIdeal(t *testing.T) {
	want := map[int]int{1: 2, 2: 3, 3: 4}
	have := IntoMapKFn[int, int](New(1, 2, 3))(
		func(k int) (v int) {
			return k + 1
		},
	)

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoMapKFnWithNilG(t *testing.T) {
	want := map[int]int{}
	have := IntoMapKFn[int, int](nil)(
		func(k int) (v int) {
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
		func(v int) (k int) {
			return v - 1
		},
	)

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoMapVFnWithNilG(t *testing.T) {
	want := map[int]int{}
	have := IntoMapVFn[int, int](nil)(
		func(v int) (k int) {
			return v - 1
		},
	)

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoMapVFnWithNilF(t *testing.T) {
	want := map[int]int{}
	have := IntoMapVFn[int, int](New(2, 3, 4))(nil)

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoChanIdeal(t *testing.T) {
	want := []int{1, 2, 3}
	have := make([]int, 0, 3)

	ctx, ctxCancel := context.WithCancel(context.Background())
	go func() {
		defer ctxCancel()
		for v := range IntoChan(ctx, New(want...)) {
			have = append(have, v)
		}
	}()

	select {
	case <-time.After(time.Second * 2):
		t.Fatal("test hung")
	case <-ctx.Done():
	}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoChanWithCancel(t *testing.T) {
	want := []int{1}
	have := make([]int, 0, 1)

	ctx, ctxCancel := context.WithCancel(context.Background())
	go func() {
		for v := range IntoChan(ctx, New(1, 2, 3)) {
			have = append(have, v)
			ctxCancel()
		}
	}()

	select {
	case <-time.After(time.Second * 2):
		t.Fatal("test hung")
	case <-ctx.Done():
	}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoChanWithNilCtx(t *testing.T) {
	want := []int{1, 2, 3}
	have := make([]int, 0, 3)

	ctx, ctxCancel := context.WithCancel(context.Background())
	go func() {
		defer ctxCancel()
		for v := range IntoChan(nil, New(want...)) {
			have = append(have, v)
		}
	}()

	select {
	case <-time.After(time.Second * 2):
		t.Fatal("test hung")
	case <-ctx.Done():
	}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}

func TestIntoChanWithNilG(t *testing.T) {
	want := []int{}
	have := make([]int, 0, 3)

	ctx, ctxCancel := context.WithCancel(context.Background())
	go func() {
		defer ctxCancel()
		for v := range IntoChan[int](ctx, nil) {
			have = append(have, v)
		}
	}()

	select {
	case <-time.After(time.Second * 2):
		t.Fatal("test hung")
	case <-ctx.Done():
	}

	assertEq("r", want, have, func(s string) { t.Fatal(s) })
}
