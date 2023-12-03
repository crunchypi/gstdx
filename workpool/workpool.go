package workpool

import "sync"

type NewArgs[I, O any] struct {
	N        int
	Work     <-chan I
	WorkEval func(I) bool
	WorkFn   func(I) O
}

func (args NewArgs[_, _]) Ok() (ok bool) {
	ok = true
	return ok
}

func New[I, O any](args NewArgs[I, O]) <-chan O {
	r := make(chan O)

	wg := sync.WaitGroup{}
	wg.Add(args.N)

	for i := 0; i < args.N; i++ {
		go func() {
			defer wg.Done()
			for vi := range args.Work {
				if args.WorkEval != nil && !args.WorkEval(vi) {
					continue
				}

				vo := args.WorkFn(vi)
				r <- vo
			}

		}()
	}

	go func() {
		wg.Wait()
		close(r)
	}()

	return r
}
