package preharvest

import (
	"strconv"
	"time"

	"github.com/Copped-Inc/aio-types/captcha/preharvest"
)

func ParseSchedule(schedule string) (duration time.Duration) {
	input := preharvest.Schedule_Pattern.FindStringSubmatch(schedule)
	multiply := 1
	var value string

	if input[4] != "" {
		switch input[4] {
		case "daily":
			duration = time.Until(time.Now().AddDate(0, 0, 1))
		case "weekly":
			duration = time.Until(time.Now().AddDate(0, 0, 7))
		case "monthly":
			duration = time.Until(time.Now().AddDate(0, 1, 0))
		}

	} else {
		if input[3] != "" {
			value = input[3]
		} else {
			multiplicator, _ := strconv.ParseInt(input[1], 10, 64)
			multiply = int(multiplicator)
			value = input[2]
		}

		switch value {
		case "day":
			duration = time.Until(time.Now().AddDate(0, 0, 1))
		case "week":
			duration = time.Until(time.Now().AddDate(0, 0, 7))
		case "month":
			duration = time.Until(time.Now().AddDate(0, 1, 0))
		}
	}

	return duration * time.Duration(multiply)
}
