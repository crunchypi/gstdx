package iox

import (
	"bytes"
	"context"
	"encoding/gob"
	"io"
	"testing"
)

func TestNewV2BWriterIdeal(t *testing.T) {
	wo := bytes.NewBuffer(nil)
	wx := NewV2BWriterFn[int](wo)(nil)

	var v int
	var err error

	err = wx.Write(nil, 79)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })

	err = wx.Write(nil, 80)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })

	err = gob.NewDecoder(wo).Decode(&v)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("r", 79, v, func(s string) { t.Fatal(s) })

	err = gob.NewDecoder(wo).Decode(&v)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("r", 80, v, func(s string) { t.Fatal(s) })

	err = gob.NewDecoder(wo).Decode(&v)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("r", 80, v, func(s string) { t.Fatal(s) })
}

func TestNewB2VWriterIdeal(t *testing.T) {
	var v int
	var err error

	wx := WriterImpl[int]{func(_ context.Context, _v int) error { v = _v; return nil }}
	wo := NewB2VWriterFn(wx)(nil)

	err = gob.NewEncoder(wo).Encode(1)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("r", 1, v, func(s string) { t.Fatal(s) })

	err = gob.NewEncoder(wo).Encode(8)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("r", 8, v, func(s string) { t.Fatal(s) })
}

func TestWriteFilterFn(t *testing.T) {
	var w Writer[int]
	var v int
	var err error

	w = WriterImpl[int]{
		Impl: func(_ context.Context, _v int) error {
			v = _v
			return nil
		},
	}

	w = WriteFilterFn(w)(func(v int) bool { return v%2 != 0 })

	err = w.Write(nil, 1)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 1, v, func(s string) { t.Fatal(s) })

	err = w.Write(nil, 2)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 1, v, func(s string) { t.Fatal(s) })

	err = w.Write(nil, 3)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 3, v, func(s string) { t.Fatal(s) })
}

func TestWMapFn(t *testing.T) {
	var w Writer[int]
	var v int
	var err error

	w = WriterImpl[int]{
		Impl: func(_ context.Context, _v int) error {
			v = _v
			return nil
		},
	}

	w = WriteMapFn[int, int](w)(func(v int) int { return v + 1 })

	err = w.Write(nil, 1)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 2, v, func(s string) { t.Fatal(s) })

	err = w.Write(nil, 2)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 3, v, func(s string) { t.Fatal(s) })
}
