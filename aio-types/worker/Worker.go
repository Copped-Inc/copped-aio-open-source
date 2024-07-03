package worker

import (
	"errors"
	"time"
)

// Worker Creates a self-healing go routine that executes the provided function exec. Parameters:
//   - exec: function to be executed by the worker
//   - config: optional Limiter object to handle repeated error occurence
//
// An error may only be returned, when providing a Limiter, as this is the only parameter that's being validated. Special cases are:
//
//	Worker(func()) = nil // No error handling is required if no limiter is passed.
//	Worker(func(), Limiter, Limiter, ...) = error // One mustn't use more than exactly one listener, otherwise Worker() will always return an error.
func Worker(exec func(), config ...Limiter) error {
	var log []time.Time
	var limiter *Limiter

	if len(config) > 0 {
		if len(config) > 1 {
			return errors.New("one may not provide more than one limiter object in a call to the Worker function")
		}

		limiter = &config[0]

		if limiter.Limit < 1 {
			return errors.New("the limit set for the Limiter in call to Worker must be greater or equal to 1")
		}
	}

	go selfRecover(exec, &log, limiter)

	return nil
}
