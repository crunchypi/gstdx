package iox

import (
	"bytes"
	"encoding/json"
	"io"
)

type Buffer[T any] struct {
	writer Writer[T]
	buf    *bytes.Buffer
}

func NewBuffer[T any](w Writer[T]) Buffer[T] {
	r := Buffer[T]{}
	r.writer = w
	r.buf = bytes.NewBuffer([]byte{}) // <-- TODO
	return r
}

func (buf Buffer[T]) Write(v T) error {
	b, _ := json.Marshal(v)
	buf.buf.Read(b)
	return nil
}

func (buf Buffer[T]) Read() (r T, err error) {
	b, err := json.Marshal(r)
	if err != nil {
		return r, err
	}

	_, err = buf.buf.Read(b)
	if err != nil {
		return r, err
	}

	err = json.Unmarshal(b, &r)
	return r, err
}

type buffer[T any] struct {
	Buffer[T]
}

func (buf buffer[T]) Read(b []byte) (int, error) {
	return buf.buf.Read(b)
}

func (buf buffer[T]) Write(b []byte) (int, error) {
	return buf.buf.Write(b)
}

func IntoIOReader[T any](r Reader[T]) io.Reader {
	v, err := r.Read()
	if err != nil {
		return bytes.NewReader([]byte{})
	}

	b, _ := json.Marshal(v)
	return bytes.NewReader(b)
}

func IntoIOWriter[T any](r Writer[T]) io.Writer {
	return buffer[T]{Buffer[T]{writer: r, buf: bytes.NewBuffer([]byte{})}}
}
