package console

import (
	"fmt"
	"github.com/Copped-Inc/aio-types/helper"
	"github.com/Copped-Inc/aio-types/webhook"
	"time"
)

func Loop() {
	if helper.System == "windows" {
		return
	}

	for {
		time.Sleep(1 * time.Minute)
		if len(LoopText) == 0 {
			continue
		}

		wh := webhook.New()
		text := "\n"
		for _, t := range LoopText {
			parsed := fmt.Sprintf("%v", t)
			if len(text)+len(parsed) >= 2048 {
				wh.AddEmbed(webhook.Log, text)
				text = "\n"
			} else {
				text = text + parsed[1:len(parsed)-1] + "\n"
			}
		}

		wh.AddEmbed(webhook.Log, text)
		_ = wh.Send(helper.Webhook)
		LoopText = nil
	}
}

var LoopText []any
