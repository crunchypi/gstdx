package testingx

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func AssertEq[T, U any](subject string, a T, b U, f func(string)) {
	if f == nil {
		return
	}

	ab, _ := json.Marshal(a)
	bb, _ := json.Marshal(b)

	as := string(ab)
	bs := string(bb)

	if as == bs {
		return
	}

	s := "unexpected '%v':\n\twant: '%v'\n\thave: '%v'\n"
	f(fmt.Sprintf(s, subject, as, bs))
}

func AssertNeq[T, U any](subject string, a T, b U, f func(string)) {
	if f == nil {
		return
	}

	ab, _ := json.Marshal(a)
	bb, _ := json.Marshal(b)

	as := string(ab)
	bs := string(bb)

	if as != bs {
		return
	}

	s := "unexpected '%v':\n\thave: '%v'\n"
	f(fmt.Sprintf(s, subject, as))
}

type ResponseWriterImpl struct {
	ImplHeader      func() http.Header
	ImplWrite       func([]byte) (int, error)
	ImplWriteHeader func(int)
}

func (impl ResponseWriterImpl) Header() http.Header {
	if impl.ImplHeader == nil {
		panic("testingx: used ResponseWriterImpl.Header without impl")
	}
	return impl.ImplHeader()
}

func (impl ResponseWriterImpl) Write(b []byte) (int, error) {
	if impl.ImplWrite == nil {
		panic("testingx: used ResponseWriterImpl.Write without impl")
	}

	return impl.ImplWrite(b)
}

func (impl ResponseWriterImpl) WriteHeader(statusCode int) {
	if impl.ImplWriteHeader == nil {
		panic("testingx: used ResponseWriterImpl.WriteHeader without impl")
	}

	impl.ImplWriteHeader(statusCode)
}

func NewResponseWriter[T any](
	want T,
	wantStatus int,
	onErr func(string),
) (
	r http.ResponseWriter,
) {
	return ResponseWriterImpl{
		ImplWriteHeader: func(status int) {
			const s = "status code (ResponseWriter.WriteHeader)"
			AssertEq(s, wantStatus, status, onErr)
		},
		ImplWrite: func(b []byte) (int, error) {
			var v T
			err := json.Unmarshal(b, &v)
			const s1 = "unmarshal err (ResponseWriter.Write)"
			AssertEq(s1, *new(error), err, onErr)

			const s2 = "response body (ResponseWriter.Write)"
			AssertEq(s2, want, v, onErr)

			return 0, nil
		},
	}
}
