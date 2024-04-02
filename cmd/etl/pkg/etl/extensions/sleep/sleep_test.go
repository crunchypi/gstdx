package sleep

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

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

func TestNewSleepNormalisedReader(t *testing.T) {
	r := iox.NewV2VReader(0, 1, 2, 3, 4, 5, 6, 7, 8, 9)

	for i := 0; i < 1; i++ {
		r = newSleepFuzzReader(r)
	}
	for i := 0; i < 9; i++ {
		r = NewSleepVReader(r, time.Second)
	}

	stamp := time.Now()
	for v, ok, _ := r.Read(nil); ok; v, ok, _ = r.Read(nil) {
		now := time.Now()
		t.Log("###rf", v, ok, now.Sub(stamp))

		stamp = now
	}
}

func TestNewSleepVWriter(t *testing.T) {
	var w iox.Writer[int]

	stamp := time.Now()
	w = iox.WriterImpl[int]{
		Impl: func(ctx context.Context, v int) error {
			delta := time.Now().Sub(stamp)
			stamp = time.Now()

			t.Log("###r", v, delta)
			return nil
		},
	}

	w = NewSleepVWriter(w, time.Second)
	w = newSleepFuzzWriter(w)

	w.Write(nil, 1)
	w.Write(nil, 2)
	w.Write(nil, 3)
}
