package worker

import (
	"time"
)

// Optional argument to be passed to the Worker() function:
//
//	Worker(func(), Limiter)
//
// Used to stop the execution of the function provided to the worker, if a certain amount of errors occurs within a certain timeframe.
type Limiter struct {
	Interval time.Duration // The time within the amount of errors may not be exceed the limit.
	Timeout  time.Duration // The time to wait after an error before restarting the worker.

	Limit int // The maximum amount of errors to occur within the specified interval.

	// Optional handler function to be executed once the error limit for the specified interval is exceeded.
	// The error that's being passed over is the most recent one that occured during the execution of the function by the worker.
	// Similarly to the handler function being optional, handling the error being passed over is mandatory too, as it's only passed over to satisfy any custom handler usecase.
	Handler func(error)
}
