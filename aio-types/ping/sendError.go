package ping

import (
	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/statistic"
	"github.com/Copped-Inc/aio-types/webhook"
	"time"
)

func sendError(domain string, error error) {
	if status[domain] == 0 {
		time.Sleep(time.Second)
		err := request(domain)
		if err == nil {
			return
		}

		console.Log("Ping failed", domain, error)
		status[domain] = time.Now().UnixMilli()
		err = statistic.SetOffline(domain)
		if err != nil {
			console.Log("Error", "SetOffline", error.Error())
		}

		err = webhook.New().AddEmbed(webhook.PingFailed, domain, error.Error()).Send("") // INSERT Webhook URL here
		if err != nil {
			console.Log("Error", "SendWebhook", error.Error())
		}
	}
}
