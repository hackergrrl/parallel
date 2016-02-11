package main

import (
	"errors"
)

var Canceled = errors.New("canceled")

func Parallel(cancel <-chan struct{}, fns ...func() error) error {
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
				return e
			}
		case <-cancel:
			return Canceled
		}
	}

	return nil
}
