package iox

import (
	"testing"
)

//type Filter struct {
//	Skip   int
//	Limit  int
//	Filter map[string]any
//}

//func newDBCaller[T any](url, usr, pwd string) Swapper[Filter, Reader[T]] {
//	return SwapperImpl[Filter, Reader[T]]{
//		Impl: func(ctx context.Context, filter Filter) (r Reader[T], err error) {
//			return
//		},
//	}
//}

func TestSwap(t *testing.T) {
	//s0 := SwapperImpl[int, int]{
	//	Impl: func(ctx context.Context, v int) (r int, err error) {
	//		r = v
	//		return
	//	},
	//}

	//// What i need to send to, what i have.
	//s1 := SwapI2MapFn[int, string](s0)(
	//	func(v string) int {
	//		x, _ := strconv.Atoi(v)
	//		return x
	//	},
	//)

	//s1 = s1

	////s0.Swap(nil, 1)
	//t.Log(s1.Swap(nil, "1"))

}
