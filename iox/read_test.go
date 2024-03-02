package iox

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"io"
	"testing"
)

func TestNewV2VReaderIdeal(t *testing.T) {
	r := NewV2VReader[int](79, 89)

	var val int
	var err error

	val, err = r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 79, val, func(s string) { t.Fatal(s) })

	val, err = r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 89, val, func(s string) { t.Fatal(s) })

	val, err = r.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", 0, val, func(s string) { t.Fatal(s) })
}

func TestNewB2VReaderIdeal(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)

	var val int
	var err error

	err = enc.Encode(79)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })

	err = enc.Encode(89)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })

	r := NewB2VReaderFn[int](buf)(nil)

	val, err = r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 79, val, func(s string) { t.Fatal(s) })

	val, err = r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 89, val, func(s string) { t.Fatal(s) })

	val, err = r.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", 0, val, func(s string) { t.Fatal(s) })
}

func TestNewD2VReaderIdeal(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)

	var val int
	var err error

	err = enc.Encode(79)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })

	err = enc.Encode(89)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })

	r := NewD2VReader[int](gob.NewDecoder(buf))

	val, err = r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 79, val, func(s string) { t.Fatal(s) })

	val, err = r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 89, val, func(s string) { t.Fatal(s) })

	val, err = r.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", 0, val, func(s string) { t.Fatal(s) })
}

func TestNewV2BReaderIdeal(t *testing.T) {
	rx := NewV2VReader(9, 8)
	ro := NewV2BReaderFn(rx)(nil)

	var val int
	var err error

	err = gob.NewDecoder(ro).Decode(&val)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 9, val, func(s string) { t.Fatal(s) })

	err = gob.NewDecoder(ro).Decode(&val)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 8, val, func(s string) { t.Fatal(s) })

	err = gob.NewDecoder(ro).Decode(&val)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", 8, val, func(s string) { t.Fatal(s) })
}

func TestNewV2BReaderWithJsonEncoderIdeal(t *testing.T) {
	type point struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	}

	p1 := point{X: 8, Y: 9}
	p2 := point{X: 5, Y: 3}

	rx := NewV2VReader(p1, p2)
	ro := NewV2BReaderFn(rx)(func(w io.Writer) Encoder { return json.NewEncoder(w) })

	var val point
	var err error

	err = json.NewDecoder(ro).Decode(&val)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", p1, val, func(s string) { t.Fatal(s) })

	err = json.NewDecoder(ro).Decode(&val)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", p2, val, func(s string) { t.Fatal(s) })

	err = json.NewDecoder(ro).Decode(&val)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", p2, val, func(s string) { t.Fatal(s) })
}

func TestNewBatchedVReaderIdeal(t *testing.T) {
	r := NewBatchedVReader(NewV2VReader(1, 2, 3), 2)
	for v, err := r.Read(nil); !errors.Is(err, io.EOF); v, err = r.Read(nil) {
		t.Log(v)
	}
}

func TestReadFilterFn(t *testing.T) {
	var r Reader[int]
	var val int
	var err error

	r = NewV2VReader(1, 2, 3)
	r = ReadFilterFn(r)(func(v int) bool { return v%2 != 0 })

	val, err = r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 1, val, func(s string) { t.Fatal(s) })

	val, err = r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 3, val, func(s string) { t.Fatal(s) })

	val, err = r.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", 0, val, func(s string) { t.Fatal(s) })
}

func TestReadMapFn(t *testing.T) {
	var r Reader[int]
	var val int
	var err error

	r = NewV2VReader(1, 2)
	r = ReadMapFn[int, int](r)(func(v int) int { return v + 1 })

	val, err = r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 2, val, func(s string) { t.Fatal(s) })

	val, err = r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 3, val, func(s string) { t.Fatal(s) })

	val, err = r.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", 0, val, func(s string) { t.Fatal(s) })
}

func TestReadReduceFn(t *testing.T) {
	r := NewV2VReader(1, 2, 3)
	var val int
	var err error

	val, err = ReadReduceFn(r)(func(acc, cur int) int { return acc + cur })
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 6, val, func(s string) { t.Fatal(s) })
}
