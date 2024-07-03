package statistic

import (
	"github.com/Copped-Inc/aio-types/statistic/realtimedb"
	"strconv"
	"strings"
	"time"
)

func SetOnline(server string, duration int64) error {

	now := time.Now()
	year, month, day := now.Date()

	today := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	t := now.Sub(today)

	if t.Milliseconds() > duration {
		return save(year, month, day, server, duration)
	}

	err := save(year, month, day, server, t.Milliseconds())
	if err != nil {
		return err
	}

	duration -= t.Milliseconds()
	for duration > 0 {
		now = now.AddDate(0, 0, -1)
		year, month, day = now.Date()
		err = save(year, month, day, server, duration)
		if err != nil {
			return err
		}

		duration -= 86400000
	}

	return err

}

type uptime struct {
	Status           int64      `json:"status"`
	DowntimeDuration int64      `json:"downtime_duration"`
	Downtimes        []downTime `json:"downtime"`
}

type downTime struct {
	End      time.Time `json:"time"`
	Duration int64     `json:"duration"`
}

func save(year int, month time.Month, day int, server string, duration int64) error {

	globaltotalRef := realtimedb.GetDatabase().NewRef("uptime/" + strings.ReplaceAll(server, ".", "-") + "/" + strconv.Itoa(year) + "-" + month.String() + "-" + strconv.Itoa(day))

	var u uptime
	err := globaltotalRef.Get(realtimedb.GetContext(), &u)
	if err != nil {
		return err
	}

	u.Status = 0
	u.DowntimeDuration += duration
	u.Downtimes = append(u.Downtimes, downTime{
		End:      time.Now(),
		Duration: duration,
	})

	return globaltotalRef.Set(realtimedb.GetContext(), &u)

}
