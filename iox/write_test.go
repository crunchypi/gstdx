package iox

import (
	"bytes"
	"context"
	"encoding/gob"
	"io"
	"testing"
)

// TODO Modify for xxxFn
func TestNewV2BWriterIdeal(t *testing.T) {
	wo := bytes.NewBuffer(nil)
	wx := NewV2BWriterFn[int](wo)(nil)

	var val int
	var err error

	err = wx.Write(nil, 79)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })

	err = wx.Write(nil, 80)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })

	err = gob.NewDecoder(wo).Decode(&val)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 79, val, func(s string) { t.Fatal(s) })

	err = gob.NewDecoder(wo).Decode(&val)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 80, val, func(s string) { t.Fatal(s) })

	err = gob.NewDecoder(wo).Decode(&val)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", 80, val, func(s string) { t.Fatal(s) })
}

func TestNewB2VWriterIdeal(t *testing.T) {
	var val int
	var err error

	wx := WriterImpl[int]{func(_ context.Context, v int) error { val = v; return nil }}
	wo := NewB2VWriterFn(wx)(nil)

	err = gob.NewEncoder(wo).Encode(1)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 1, val, func(s string) { t.Fatal(s) })

	err = gob.NewEncoder(wo).Encode(8)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 8, val, func(s string) { t.Fatal(s) })
}
