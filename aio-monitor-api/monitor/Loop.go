package monitor

import (
	"time"
)

func Loop(insert func() error, delay time.Duration) func() {
	return func() {
		for {
			err := insert()
			if err != nil {
				panic(err)
			}
			time.Sleep(delay)
		}
	}
}
