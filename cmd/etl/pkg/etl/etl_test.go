package etl

import (
	"context"
	"testing"
	"time"

	"github.com/crunchypi/gstdx/iox"
)

func TestPageReader(t *testing.T) {
	r := NewPageReader(iox.NewV2VReader(4, 3), 2)
	//r := newPageReader(3, 2)
	t.Log(r.Read(nil))
	t.Log(r.Read(nil))
	t.Log(r.Read(nil))
	t.Log(r.Read(nil))
	t.Log(r.Read(nil))
}

func TestNewSleepNormalisedReader(t *testing.T) {
	r := iox.NewV2VReader(1, 2, 3, 4, 5, 6, 7, 8, 9)

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

func TestNewStatReader(t *testing.T) {
	r := iox.NewV2VReader(1, 2, 3, 4, 5, 6, 7, 8, 9)
	r = NewStatsReader(r, "dank")(
		iox.WriterImpl[ReadStat[int]]{
			Impl: func(ctx context.Context, v ReadStat[int]) error {
				t.Logf("%+v\n", v)
				return nil
			},
		},
	)

	for v, ok, _ := r.Read(nil); ok; v, ok, _ = r.Read(nil) {
		v = v
	}
}

func TestNewBatchedVReader(t *testing.T) {
	vr := iox.NewV2VReader(1, 2, 3)
	sr := NewBatchedVReader(vr, 2)

	for v, ok, err := sr.Read(nil); ok; v, ok, err = sr.Read(nil) {
		t.Log(v, ok, err)
	}
}

func TestNewBatchedSReader(t *testing.T) {
	vr := iox.NewV2VReader(1, 2, 3)
	sr := NewBatchedSReader(vr, 2)

	for v, ok, err := sr.Read(nil); ok; v, ok, err = sr.Read(nil) {
		t.Log(v, ok, err)
	}
}
