package stats

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/crunchypi/gstdx/iox"
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

func TestNewReaderIdeal(t *testing.T) {
	vr := iox.NewV2VReader(1, 2)
	sr := NewReader(vr, "test", "k")

	var ctx = context.Background()
	var stat Stat[int]
	var err error

	ctx = context.WithValue(ctx, "k", "v")
	stat, err = sr.Read(ctx)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("tag", "test", stat.Tag, func(s string) { t.Fatal(s) })
	assertEq("val", 1, stat.Val, func(s string) { t.Fatal(s) })
	assertEq("err", *new(error), stat.Err, func(s string) { t.Fatal(s) })
	assertEq("ctx", map[string]any{"k": "v"}, stat.CtxVals, func(s string) { t.Fatal(s) })

	ctx = context.WithValue(ctx, "k", "w")
	stat, err = sr.Read(ctx)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("tag", "test", stat.Tag, func(s string) { t.Fatal(s) })
	assertEq("val", 2, stat.Val, func(s string) { t.Fatal(s) })
	assertEq("err", *new(error), stat.Err, func(s string) { t.Fatal(s) })
	assertEq("ctx", map[string]any{"k": "w"}, stat.CtxVals, func(s string) { t.Fatal(s) })

	ctx = context.WithValue(ctx, "k", "w")
	stat, err = sr.Read(ctx)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("tag", "", stat.Tag, func(s string) { t.Fatal(s) })
	assertEq("val", 0, stat.Val, func(s string) { t.Fatal(s) })
	assertEq("err", *new(error), stat.Err, func(s string) { t.Fatal(s) })
	assertEq("ctx", nil, stat.CtxVals, func(s string) { t.Fatal(s) })
}

func TestNewTeeReaderIdeal(t *testing.T) {
	rw := iox.NewV2VReadWriter[Stat[int]]()
	sr := NewTeeReaderFn(iox.NewV2VReader(1, 2), "test", "k")(rw)

	var stat Stat[int]
	var err error

	ctx := context.Background()
	ctx = context.WithValue(ctx, "k", "v")

	// Fill up.
	for _, err := sr.Read(ctx); !errors.Is(err, io.EOF); _, err = sr.Read(ctx) {
	}

	stat, err = rw.Read(ctx)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("tag", "test", stat.Tag, func(s string) { t.Fatal(s) })
	assertEq("val", 1, stat.Val, func(s string) { t.Fatal(s) })
	assertEq("err", *new(error), stat.Err, func(s string) { t.Fatal(s) })
	assertEq("ctx", map[string]any{"k": "v"}, stat.CtxVals, func(s string) { t.Fatal(s) })

	ctx = context.WithValue(ctx, "k", "v")
	stat, err = rw.Read(ctx)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("tag", "test", stat.Tag, func(s string) { t.Fatal(s) })
	assertEq("val", 2, stat.Val, func(s string) { t.Fatal(s) })
	assertEq("err", *new(error), stat.Err, func(s string) { t.Fatal(s) })
	assertEq("ctx", map[string]any{"k": "v"}, stat.CtxVals, func(s string) { t.Fatal(s) })

	ctx = context.WithValue(ctx, "k", "v")
	stat, err = rw.Read(ctx)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("tag", "", stat.Tag, func(s string) { t.Fatal(s) })
	assertEq("val", 0, stat.Val, func(s string) { t.Fatal(s) })
	assertEq("err", *new(error), stat.Err, func(s string) { t.Fatal(s) })
	assertEq("ctx", nil, stat.CtxVals, func(s string) { t.Fatal(s) })
}

func TestNewTeeWriterIdeal(t *testing.T) {
	rw := iox.NewV2VReadWriter[Stat[int]]()
	sr := NewTeeWriterFn(iox.NewV2VReadWriter[int](), "test", "k")(rw)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "k", "v")

	sr.Write(ctx, 1)
	sr.Write(ctx, 2)

	var stat Stat[int]
	var err error

	stat, err = rw.Read(ctx)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("tag", "test", stat.Tag, func(s string) { t.Fatal(s) })
	assertEq("val", 1, stat.Val, func(s string) { t.Fatal(s) })
	assertEq("err", *new(error), stat.Err, func(s string) { t.Fatal(s) })
	assertEq("ctx", map[string]any{"k": "v"}, stat.CtxVals, func(s string) { t.Fatal(s) })

	ctx = context.WithValue(ctx, "k", "v")
	stat, err = rw.Read(ctx)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("tag", "test", stat.Tag, func(s string) { t.Fatal(s) })
	assertEq("val", 2, stat.Val, func(s string) { t.Fatal(s) })
	assertEq("err", *new(error), stat.Err, func(s string) { t.Fatal(s) })
	assertEq("ctx", map[string]any{"k": "v"}, stat.CtxVals, func(s string) { t.Fatal(s) })

	ctx = context.WithValue(ctx, "k", "v")
	stat, err = rw.Read(ctx)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("tag", "", stat.Tag, func(s string) { t.Fatal(s) })
	assertEq("val", 0, stat.Val, func(s string) { t.Fatal(s) })
	assertEq("err", *new(error), stat.Err, func(s string) { t.Fatal(s) })
	assertEq("ctx", nil, stat.CtxVals, func(s string) { t.Fatal(s) })
}

/*
func TestX(t *testing.T) {
	ir := iox.NewV2VReader(1, 2, 3)
	sr := NewMappedReader(ir, "test", "k")
	lr := log.NewReader(log.NewArgs[Stat[int]]{Reader: sr})

	i := 0
	for v, err := lr.Read(nil); !errors.Is(err, io.EOF); v, err = lr.Read(nil) {
		t.Log(v)

		if i >= 3 {
			break
		}

		i++
	}
}
*/
