package ping

import (
	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/statistic"
	"github.com/Copped-Inc/aio-types/webhook"
	"time"
)

var status = make(map[string]int64)

func Check(domain string) {

	err := request(domain)
	if err != nil {
		sendError(domain, err)
		return
	}

	if status[domain] != 0 {
		console.Log("Domain", domain, "is back online")
		err = webhook.New().AddEmbed(webhook.PingSuccess, domain).Send("") // Insert Webhook URL here
		duration := time.Now().UnixMilli() - status[domain]
		status[domain] = 0

		err = statistic.SetOnline(domain, duration)
		if err != nil {
			console.Log("Error", "SetOnline", err.Error())
		}
	}

}
