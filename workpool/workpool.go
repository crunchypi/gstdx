package workpool

import (
	"sync"
)

type Result[T any] struct {
	Val T
	Err error
}

type Option[T any] struct {
	Val T
	Ok  bool
}

type NewArgs[I, O any] struct {
	N      int
	Work   <-chan I
	Filter func(I) bool
	Mapper func(I) O
}

func (args NewArgs[_, _]) Ok() (ok bool) {
	ok = true
	ok = ok && args.N > 0
	ok = ok && args.Work != nil
	ok = ok && args.Mapper != nil

	return ok
}

func New[I, O any](args NewArgs[I, O]) <-chan O {
	r := make(chan O)
	if !args.Ok() {
		close(r)
		return r
	}

	wg := sync.WaitGroup{}
	wg.Add(args.N)

	for i := 0; i < args.N; i++ {
		go func() {
			defer wg.Done()

			for vi := range args.Work {
				if args.Filter != nil && !args.Filter(vi) {
					continue
				}

				r <- args.Mapper(vi)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(r)
	}()

	return r
}
