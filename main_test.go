package main

import (
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/crunchypi/gstdx/iox"
)

func fileWriter[T any](
	p string,
	enc ...func(io.Writer) iox.Encoder,
) (
	w iox.WriteCloser[T],
	err error,
) {
	f, err := os.Create(p)
	if err != nil {
		return iox.WriteCloserImpl[T]{}, err
	}

	w = iox.NewV2BWriteCloser[T](f, enc...)
	return w, err
}

func fileReader[T any](
	p string,
	dec ...func(io.Reader) iox.Decoder,
) (
	r iox.ReadCloser[T],
	err error,
) {
	f, err := os.Open(p)
	if err != nil {
		return iox.ReadCloserImpl[T]{}, err
	}

	r = iox.NewB2VReadCloser[T](f, dec...)
	return r, err
}

func TestWriteFile(t *testing.T) {
	type point struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	}

	enc := func(w io.Writer) iox.Encoder { return json.NewEncoder(w) }
	w, err := fileWriter[point]("./test.json", enc)
	if err != nil {
		t.Fatal(err)
	}

	w.Write(nil, point{1, 2})
	w.Write(nil, point{2, 3})
	err = w.Close()
	if err != nil {
		t.Fatal(err)
	}
	err = w.Write(nil, point{2, 3})
	if err != nil {
		t.Fatal(err)
	}
}

func TestReadFile(t *testing.T) {
	dec := func(w io.Reader) iox.Decoder { return json.NewDecoder(w) }
	r, err := fileReader[int]("./test.json", dec)
	if err != nil {
		t.Fatal(err)
	}

	r.Read(nil)
	r.Read(nil)
	r.Close()
}

func TestStuff(t *testing.T) {
	f, err := os.Open("")
	if err != nil {
		t.Fatal(err)
	}

	var v int
	err = json.NewDecoder(f).Decode(v)
	if err != nil {
		t.Fatal(err)
	}
}
