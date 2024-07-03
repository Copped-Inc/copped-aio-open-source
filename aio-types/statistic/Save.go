package statistic

import (
	"fmt"
	"github.com/Copped-Inc/aio-types/helper"
	"github.com/Copped-Inc/aio-types/statistic/realtimedb"
	"net/http"
	"strconv"
	"time"
)

func (stat statistic) SaveRequest() error {

	if stat.Path == "/ping" {
		return nil
	}

	stat.End = time.Now()
	stat.Duration = stat.End.Sub(stat.Start)

	year, month, day := time.Now().Date()
	globaltotalRef := realtimedb.GetDatabase().NewRef("serverstats/" + realtimedb.GetServer() + "/" + strconv.Itoa(year) + "-" + month.String() + "-" + strconv.Itoa(day) + "/" + realtimedb.Requests + "/" + stat.Id)
	return globaltotalRef.Set(realtimedb.GetContext(), &stat)

}

func SaveLog(r *http.Request, id string, text ...any) {

	if realtimedb.GetContext() == nil || !helper.GeneralLog {
		return
	}

	l := log{
		Id: id,
		State: func() state {
			if text[0] == "Error" || text[0] == "Panic" {
				return stateErr
			} else {
				return stateOk
			}
		}(),
		Ref: func() string {
			if r == nil {
				return ""
			}
			return r.Context().Value("statistic").(statistic).Id
		}(),
		Time:    time.Now(),
		Content: text,
	}

	year, month, day := time.Now().Date()
	globaltotalRef := realtimedb.GetDatabase().NewRef("serverstats/" + realtimedb.GetServer() + "/" + strconv.Itoa(year) + "-" + month.String() + "-" + strconv.Itoa(day) + "/" + realtimedb.Logs + "/" + l.Id)
	err := globaltotalRef.Set(realtimedb.GetContext(), &l)
	if err != nil {
		fmt.Println("Error", err.Error())
	}

}
