package workpool

import (
	"encoding/json"
	"runtime"
	"testing"

	"github.com/crunchypi/gstdx/chanx"
)

func TestX(t *testing.T) {
	ch := New(
		NewArgs[any, Result[string]]{
			N:    runtime.NumCPU(),
			Work: chanx.New[any](1, "2", "a", "b", "3", 1.1),
			Mapper: func(v any) (r Result[string]) {
				b, err := json.Marshal(v)
				if err != nil {
					return Result[string]{Err: err}
				}

				r.Err = json.Unmarshal(b, &r.Val)
				return r
			},
		},
	)

	for i, v := range chanx.IntoSlice(ch) {
		t.Log(i, v)
	}
}
