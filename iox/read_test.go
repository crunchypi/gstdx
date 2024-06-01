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

func TestNewByteReadCloserFnIdeal(t *testing.T) {
	closed := false
	vrc := ReadCloserImpl[int]{ImplC: func() error { closed = true; return nil }}
	brc := NewByteReadCloserFn(vrc)(nil)

	brc.Close()
	assertEq("closed", true, closed, func(s string) { t.Fatal(s) })
}

func TestNewByteReadCloserFnWithNilReader(t *testing.T) {
	brc := NewByteReadCloserFn[int](nil)(nil)
	brc.Close()
}

func TestNewBatchedValueReaderIdeal(t *testing.T) {
	vs := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	vr := tfNewValueReaderFrom(vs...)
	sr := NewBatchedValueReader(vr, 0)

	s := []int{}
	err := *new(error)

	s, err = sr.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", vs[0:8], s, func(s string) { t.Fatal(s) })

	s, err = sr.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", vs[8:], s, func(s string) { t.Fatal(s) })
}

func TestNewBatchedValueReaderWithNilReader(t *testing.T) {
	sr := NewBatchedValueReader[int](nil, 0)

	s, err := sr.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", *new([]int), s, func(s string) { t.Fatal(s) })
}

func TestNewUnbatcherIdeal(t *testing.T) {
	sr := NewBatchedValueReader(tfNewValueReaderFrom(1, 3, 2), 2)
	vr := NewUnbatchedValueReader(sr)

	err := *new(error)
	val := 0

	val, err = vr.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 1, val, func(s string) { t.Fatal(s) })

	val, err = vr.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 3, val, func(s string) { t.Fatal(s) })

	val, err = vr.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 2, val, func(s string) { t.Fatal(s) })

	val, err = vr.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", 0, val, func(s string) { t.Fatal(s) })
}

func TestNewUnbatcherWithNilReader(t *testing.T) {
	vr := NewUnbatchedValueReader[int](nil)

	val, err := vr.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", 0, val, func(s string) { t.Fatal(s) })
}

func TestNewUnbatcherWithEmptyBatchAndNilErr(t *testing.T) {
	sr := ReaderImpl[[]int]{}
	sr.Impl = func(ctx context.Context) (s []int, err error) { return }
	vr := NewUnbatchedValueReader(sr)

	val, err := vr.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", 0, val, func(s string) { t.Fatal(s) })
}

func TestNewUnbatcherWithEmptyBatchAndErr(t *testing.T) {
	sr := ReaderImpl[[]int]{}
	sr.Impl = func(ctx context.Context) (s []int, err error) { err = io.EOF; return }
	vr := NewUnbatchedValueReader(sr)

	val, err := vr.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", 0, val, func(s string) { t.Fatal(s) })
}

func TestNewValueReaderWithFilterFnIdeal(t *testing.T) {
	r := tfNewValueReaderFrom(1, 2, 3)
	r = NewValueReaderWithFilterFn(r)(func(v int) bool { return v%2 == 0 })

	err := *new(error)
	val := 0

	val, err = r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 2, val, func(s string) { t.Fatal(s) })

	val, err = r.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", 0, val, func(s string) { t.Fatal(s) })
}

func TestNewValueReaderWithFilterFnWithNilReader(t *testing.T) {
	r := NewValueReaderWithFilterFn[int](nil)(func(v int) bool { return true })

	val, err := r.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", 0, val, func(s string) { t.Fatal(s) })
}

func TestNewValueReaderWithFilterFnWithNilFunc(t *testing.T) {
	r := tfNewValueReaderFrom(2)
	r = NewValueReaderWithFilterFn(r)(nil)

	err := *new(error)
	val := 0

	val, err = r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 2, val, func(s string) { t.Fatal(s) })

	val, err = r.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", 0, val, func(s string) { t.Fatal(s) })
}

func TestNewValueReaderWithMapperFnIdeal(t *testing.T) {
	r := tfNewValueReaderFrom(1, 2)
	r = NewValueReaderWithMapperFn[int, int](r)(
		func(v int) int {
			return v * -1
		},
	)

	err := *new(error)
	val := 0

	val, err = r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", -1, val, func(s string) { t.Fatal(s) })

	val, err = r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", -2, val, func(s string) { t.Fatal(s) })

	val, err = r.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", 0, val, func(s string) { t.Fatal(s) })
}

func TestNewValueReaderWithMapperFnWithNilReader(t *testing.T) {
	r := NewValueReaderWithMapperFn[int, int](nil)(
		func(v int) int {
			return v * -1
		},
	)

	err := *new(error)
	val := 0

	val, err = r.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", 0, val, func(s string) { t.Fatal(s) })
}

func TestNewValueReaderWithMapperFnWithNilFunc(t *testing.T) {
	r := tfNewValueReaderFrom(1, 2)
	r = NewValueReaderWithMapperFn[int, int](r)(nil)

	err := *new(error)
	val := 0

	val, err = r.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", 0, val, func(s string) { t.Fatal(s) })
}
