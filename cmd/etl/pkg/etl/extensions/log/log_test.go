package log

import (
	"context"
	"testing"

	"github.com/crunchypi/gstdx/iox"
)

func TestNewReaderIdeal(t *testing.T) {
	vr := iox.NewV2VReader(1, 2, 3)
	sr := iox.ReadMapFn[int, map[string]any](vr)(
		func(v int) map[string]any { return map[string]any{"v": v} },
	)
	lr := NewReader(NewArgs[map[string]any]{Reader: sr, Tag: "test", CtxKeys: []string{"k"}})

	ctx := context.Background()
	ctx = context.WithValue(ctx, "k", "v")

	for _, ok, _ := lr.Read(ctx); ok; _, ok, _ = lr.Read(ctx) {

	}

}
