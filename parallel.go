package parallel

import (
	"golang.org/x/net/context"
)

var Canceled = context.Canceled

func Parallel(done <-chan struct{}, cancel context.CancelFunc, fns ...func() error) error {
	errs := make(chan error)
	for _, fn := range fns {
		thisFn := fn
		go func() {
			errs <- thisFn()
		}()
	}

	for count := 0; count < len(fns); count++ {
		select {
		case e := <-errs:
			if e != nil {
				if cancel != nil {
					cancel()
				}
				return e
			}
		case <-done:
			return Canceled
		}
	}

	return nil
}
