package iox

import (
	"io"
	"testing"
)

func TestNewV2VReadWriterIdeal(t *testing.T) {
	var val int
	var err error

	rw := NewV2VReadWriter[int]()

	val, err = rw.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", 0, val, func(s string) { t.Fatal(s) })

	err = rw.Write(nil, 1)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })

	val, err = rw.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 1, val, func(s string) { t.Fatal(s) })

	val, err = rw.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", 0, val, func(s string) { t.Fatal(s) })

	err = rw.Write(nil, 2)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })

	err = rw.Write(nil, 3)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })

	val, err = rw.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 2, val, func(s string) { t.Fatal(s) })

	val, err = rw.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 3, val, func(s string) { t.Fatal(s) })

	val, err = rw.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", 0, val, func(s string) { t.Fatal(s) })
}

func TestScratch(t *testing.T) {
	//rw := NewT2UReadWriter[int, string](nil, nil, nil)
	//rw.Write(nil, "a")

}
