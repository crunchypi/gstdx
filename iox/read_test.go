package iox

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"io"
	"testing"
)

func tfNewValueReaderFrom[T any](vs ...T) Reader[T] {
	i := 0
	return ReaderImpl[T]{
		Impl: func(ctx context.Context) (val T, err error) {
			if i >= len(vs) {
				return val, io.EOF
			}

			val = vs[i]
			i++
			return
		},
	}
}

func TestReaderImplReadIdeal(t *testing.T) {
	r := ReaderImpl[int]{}
	r.Impl = func(ctx context.Context) (int, error) { return 1, nil }

	val, err := r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 1, val, func(s string) { t.Fatal(s) })
}

func TestReaderImplReadWithoutImpl(t *testing.T) {
	r := ReaderImpl[int]{}

	val, err := r.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", 0, val, func(s string) { t.Fatal(s) })
}

func TestReadCloserImplReadIdeal(t *testing.T) {
	rc := ReadCloserImpl[int]{}
	rc.ImplR = func(ctx context.Context) (int, error) { return 1, nil }

	val, err := rc.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 1, val, func(s string) { t.Fatal(s) })
}

func TestReadCloserImplReadWithoutImpl(t *testing.T) {
	rc := ReadCloserImpl[int]{}

	val, err := rc.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", 0, val, func(s string) { t.Fatal(s) })
}

func TestReadCloserImplCloseIdeal(t *testing.T) {
	rc := ReadCloserImpl[int]{}
	rc.ImplC = func() error { return nil }

	err := rc.Close()
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
}

func TestReadCloserImplCloseWithoutImpl(t *testing.T) {
	rc := ReadCloserImpl[int]{}

	err := rc.Close()
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
}

func TestNewValueReaderFnIdeal(t *testing.T) {
	b := bytes.NewBuffer(nil)
	json.NewEncoder(b).Encode("test1")
	json.NewEncoder(b).Encode("test2")

	f := func(r io.Reader) Decoder { return json.NewDecoder(r) }
	r := NewValueReaderFn[string](b)(f)

	err := *new(error)
	val := ""

	val, err = r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", "test1", val, func(s string) { t.Fatal(s) })

	val, err = r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", "test2", val, func(s string) { t.Fatal(s) })

	val, err = r.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", "", val, func(s string) { t.Fatal(s) })
}

func TestNewValueReaderFnWithNilReader(t *testing.T) {
	r := NewValueReaderFn[string](nil)(nil)

	err := *new(error)
	val := ""

	val, err = r.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", "", val, func(s string) { t.Fatal(s) })
}

func TestNewValueReaderFnWithNilDecoder(t *testing.T) {
	b := bytes.NewBuffer(nil)
	gob.NewEncoder(b).Encode("test1")
	gob.NewEncoder(b).Encode("test2")

	r := NewValueReaderFn[string](b)(nil)

	err := *new(error)
	val := ""

	val, err = r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", "test1", val, func(s string) { t.Fatal(s) })

	val, err = r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", "test2", val, func(s string) { t.Fatal(s) })

	val, err = r.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", "", val, func(s string) { t.Fatal(s) })
}

func TestNewValueReadCloserFnIdeal(t *testing.T) {
	closed := false

	brc := readWriteCloserImpl{ImplC: func() error { closed = true; return nil }}
	vrc := NewValueReadCloserFn[int](brc)(nil)

	vrc.Close()
	assertEq("closed", true, closed, func(s string) { t.Fatal(s) })
}

func TestNewValueReadCloserFnWithNilReader(t *testing.T) {
	vrc := NewValueReadCloserFn[int](nil)(nil)
	vrc.Close()
}

func TestNewByteReaderFnIdeal(t *testing.T) {
	vr := tfNewValueReaderFrom("test1", "test2")
	br := NewByteReaderFn(vr)(func(w io.Writer) Encoder { return json.NewEncoder(w) })

	dec := json.NewDecoder(br)
	err := *new(error)
	val := ""

	err = dec.Decode(&val)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", "test1", val, func(s string) { t.Fatal(s) })

	err = dec.Decode(&val)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", "test2", val, func(s string) { t.Fatal(s) })

	err = dec.Decode(&val)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", "test2", val, func(s string) { t.Fatal(s) })
}

func TestNewByteReaderFnWithNilReader(t *testing.T) {
	br := NewByteReaderFn[int](nil)(func(w io.Writer) Encoder { return json.NewEncoder(w) })

	dec := json.NewDecoder(br)
	err := *new(error)
	val := ""

	err = dec.Decode(&val)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", "", val, func(s string) { t.Fatal(s) })
}

func TestNewByteReaderFnWithNilEncoder(t *testing.T) {
	vr := tfNewValueReaderFrom("test1", "test2")
	br := NewByteReaderFn(vr)(nil)

	dec := gob.NewDecoder(br)
	err := *new(error)
	val := ""

	err = dec.Decode(&val)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", "test1", val, func(s string) { t.Fatal(s) })

	err = dec.Decode(&val)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", "test2", val, func(s string) { t.Fatal(s) })

	err = dec.Decode(&val)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", "test2", val, func(s string) { t.Fatal(s) })
}

func TestNewByteReaderFnWithEncodeError(t *testing.T) {
	vr := tfNewValueReaderFrom(make(chan int))
	br := NewByteReaderFn(vr)(func(w io.Writer) Encoder { return json.NewEncoder(w) })

	dec := json.NewDecoder(br)
	err := *new(error)
	val := ""

	err = dec.Decode(&val)

	want := "json: unsupported type: chan int"
	have := err.Error()
	assertEq("err", want, have, func(s string) { t.Fatal(s) })
}
