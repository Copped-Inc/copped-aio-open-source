package log

import (
	"encoding/json"
	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/statistic/realtimedb"
	"github.com/Copped-Inc/aio-types/webhook"
	"strconv"
	"time"
)

var lastWebhook = time.Now()

func Handle(user string, instance string, l []Log) {

	for _, log := range l {
		globaltotalRef := realtimedb.GetDatabase().NewRef("userstats/user/" + user + "/logs/" + instance + "/" + strconv.Itoa(int(log.Date.UnixMicro())))

		log.User = user
		log.Instance = instance

		if log.State == Error {
			b, err := json.MarshalIndent(&log, "", "  ")
			if err != nil {
				console.ErrorLog(err)
				b = []byte(err.Error())
			}

			if time.Since(lastWebhook) > time.Second*30 {
				err = webhook.New().AddEmbed(webhook.ErrorLog, log.Instance, string(b)).Send("") // INSERT Webhook URL here
				if err != nil {
					console.ErrorLog(err)
				}
				lastWebhook = time.Now()
			}
		}

		err := globaltotalRef.Set(realtimedb.GetContext(), &log)

		if err != nil {
			console.ErrorLog(err)
		}
	}

}
