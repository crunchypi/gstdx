package page

import (
	"encoding/json"
	"fmt"
	"io"
	"testing"

	"github.com/crunchypi/gstdx/iox"
)

func assertEq[T any](subject string, a T, b T, f func(string)) {
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

func TestPageReaderIdeal(t *testing.T) {
	nr := iox.NewV2VReader(1, 2, 3)
	pr := NewPageReader(nr, 2)

	var val Page
	var err error

	val, err = pr.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", Page{Skip: 0, Limit: 1}, val, func(s string) { t.Fatal(s) })

	val, err = pr.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", Page{Skip: 0, Limit: 2}, val, func(s string) { t.Fatal(s) })

	val, err = pr.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", Page{Skip: 0, Limit: 2}, val, func(s string) { t.Fatal(s) })

	val, err = pr.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", Page{Skip: 2, Limit: 1}, val, func(s string) { t.Fatal(s) })

	val, err = pr.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", Page{Skip: 0, Limit: 0}, val, func(s string) { t.Fatal(s) })
}

/*
func testSwap(t *testing.T) {
	type T = map[string]any
	type U struct {
		t T
		w iox.Writer[T]
	}

	ew := iox.WriterImpl[U]{
		Impl: func(ctx context.Context, v U) error {

			return nil
		},
	}

	end := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		wx := iox.NewV2BWriterFn[T](w)(nil)
		rx := iox.NewB2VReaderFn[T](r.Body)(nil)

		for v, err := rx.Read(ctx); !errors.Is(err, io.EOF); v, err = rx.Read(ctx) {
			ew.Write(ctx, U{t: v, w: wx})

		}
	}

	end = end
}
*/
