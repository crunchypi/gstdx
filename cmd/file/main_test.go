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

func TestCopyFile(t *testing.T) {
	dec := func(r io.Reader) iox.Decoder { return json.NewDecoder(r) }

	//dec = func(r io.Reader) iox.Decoder {
	//	br := bufio.NewScanner(r)
	//	br.

	//
	//	return iox.DecoderImpl{}
	//}

	r, err := fileReader[map[string]any]("./test.txt", dec)
	if err != nil {
		t.Fatal(err)
	}

	defer r.Close()

	enc := func(w io.Writer) iox.Encoder { return json.NewEncoder(w) }
	w, err := fileWriter[map[string]any]("./test2.txt", enc)
	if err != nil {
		t.Fatal(err)
	}

	defer w.Close()

	for v, ok, err := r.Read(nil); ; v, ok, err = r.Read(nil) {
		t.Log(v)
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			break
		}

		if err := w.Write(nil, v); err != nil {
			t.Fatal(err)
		}
	}
}
