package iox

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"io"
	"testing"
)

func TestNewV2VReaderIdeal(t *testing.T) {
	r := NewV2VReader[int](79, 89)

	have, ok, err := r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("ok", true, ok, func(s string) { t.Fatal(s) })
	assertEq("r", 79, have, func(s string) { t.Fatal(s) })

	have, ok, err = r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("ok", true, ok, func(s string) { t.Fatal(s) })
	assertEq("r", 89, have, func(s string) { t.Fatal(s) })

	have, ok, err = r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("ok", false, ok, func(s string) { t.Fatal(s) })
	assertEq("r", 0, have, func(s string) { t.Fatal(s) })
}

func TestNewB2VReaderIdeal(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)

	err := enc.Encode(79)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	err = enc.Encode(89)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })

	r := NewB2VReader[int](buf)

	have, ok, err := r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("ok", true, ok, func(s string) { t.Fatal(s) })
	assertEq("r", 79, have, func(s string) { t.Fatal(s) })

	have, ok, err = r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("ok", true, ok, func(s string) { t.Fatal(s) })
	assertEq("r", 89, have, func(s string) { t.Fatal(s) })

	have, ok, err = r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("ok", false, ok, func(s string) { t.Fatal(s) })
	assertEq("r", 0, have, func(s string) { t.Fatal(s) })
}

func TestNewD2VReaderIdeal(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)

	err := enc.Encode(79)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	err = enc.Encode(89)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })

	r := NewD2VReader[int](gob.NewDecoder(buf))

	// TODO standardize "val" and not "r"
	have, ok, err := r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("ok", true, ok, func(s string) { t.Fatal(s) })
	assertEq("val", 79, have, func(s string) { t.Fatal(s) })

	have, ok, err = r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("ok", true, ok, func(s string) { t.Fatal(s) })
	assertEq("val", 89, have, func(s string) { t.Fatal(s) })

	have, ok, err = r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("ok", false, ok, func(s string) { t.Fatal(s) })
	assertEq("val", 0, have, func(s string) { t.Fatal(s) })
}

func TestNewV2BReaderIdeal(t *testing.T) {
	rx := NewV2VReader(9, 8)
	ro := NewV2BReader(rx)

	var val int
	var err error

	err = gob.NewDecoder(ro).Decode(&val)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("r", 9, val, func(s string) { t.Fatal(s) })
	val = 0

	err = gob.NewDecoder(ro).Decode(&val)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("r", 8, val, func(s string) { t.Fatal(s) })
	val = 0

	err = gob.NewDecoder(ro).Decode(&val)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("r", 0, val, func(s string) { t.Fatal(s) })
	val = 0
}

func TestNewV2BReaderWithJsonEncoderIdeal(t *testing.T) {
	type point struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	}

	p1 := point{X: 8, Y: 9}
	p2 := point{X: 5, Y: 3}

	rx := NewV2VReader(p1, p2)
	ef := func(w io.Writer) Encoder { return json.NewEncoder(w) }
	ro := NewV2BReader(rx, ef)

	var val point
	var err error

	err = json.NewDecoder(ro).Decode(&val)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("r", p1, val, func(s string) { t.Fatal(s) })
	val = point{}

	err = json.NewDecoder(ro).Decode(&val)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("r", p2, val, func(s string) { t.Fatal(s) })
	val = point{}

	err = json.NewDecoder(ro).Decode(&val)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("r", point{}, val, func(s string) { t.Fatal(s) })
	val = point{}
}

func TestReadFilterFn(t *testing.T) {
	var r Reader[int]
	var v int
	var ok bool
	var err error

	r = NewV2VReader(1, 2, 3)
	r = ReadFilterFn(r)(func(v int) bool { return v%2 != 0 })

	v, ok, err = r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("ok", true, ok, func(s string) { t.Fatal(s) })
	assertEq("val", 1, v, func(s string) { t.Fatal(s) })

	v, ok, err = r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("ok", true, ok, func(s string) { t.Fatal(s) })
	assertEq("val", 3, v, func(s string) { t.Fatal(s) })

	v, ok, err = r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("ok", false, ok, func(s string) { t.Fatal(s) })
	assertEq("val", 0, v, func(s string) { t.Fatal(s) })
}

func TestReadMapFn(t *testing.T) {
	var r Reader[int]
	var v int
	var ok bool
	var err error

	r = NewV2VReader(1, 2)
	r = ReadMapFn[int, int](r)(func(v int) int { return v + 1 })

	v, ok, err = r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("ok", true, ok, func(s string) { t.Fatal(s) })
	assertEq("val", 2, v, func(s string) { t.Fatal(s) })

	v, ok, err = r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("ok", true, ok, func(s string) { t.Fatal(s) })
	assertEq("val", 3, v, func(s string) { t.Fatal(s) })

	v, ok, err = r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("ok", false, ok, func(s string) { t.Fatal(s) })
	assertEq("val", 0, v, func(s string) { t.Fatal(s) })
}

func TestReadReduceFn(t *testing.T) {
	r := NewV2VReader(1, 2, 3)
	v, err := ReadReduceFn(r)(func(acc, cur int) int { return acc + cur })

	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 6, v, func(s string) { t.Fatal(s) })
}
