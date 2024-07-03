package worker

import (
	"errors"
	"reflect"
	"time"

	"github.com/Copped-Inc/aio-types/console"
)

func selfRecover(exec func(), log *[]time.Time, limiter *Limiter) {
	defer func() {
		var err error
		if recovered := recover(); recovered == nil {
			err = errors.New("function execution of worker routine failed with recovery value being nil")
		} else {
			switch recovered.(type) {
			case error:
				err = recovered.(error)
			default:
				err = errors.New("function execution of worker routine failed with recovery value being of a type different than error, being " + reflect.TypeOf(recovered).Kind().String())
			}
		}

		console.ErrorLog(err)
		if limiter != nil {
			var newlog []time.Time

			for _, entry := range *log {
				if time.Since(entry) < limiter.Interval {
					newlog = append(newlog, entry)
				}
			}

			if newlog = append(newlog, time.Now()); len(newlog) > limiter.Limit {
				if limiter.Handler != nil {
					limiter.Handler(err)
				} else {
					console.ErrorLog(errors.New("go routine failure limit has been exceeded and no custom handler was provided"))
				}
				return
			}

			*log = newlog
			time.Sleep(limiter.Timeout)
		}

		selfRecover(exec, log, limiter)
	}()

	exec()
}
