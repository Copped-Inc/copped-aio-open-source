package preharvest

import (
	"net/http"
	"service-api/handler/discord/interactions/application_commands/captcha/preharvest"

	"strings"

	consts "github.com/Copped-Inc/aio-types/captcha/preharvest"

	"github.com/infinitare/disgo"
)

func Buttons(interaction *disgo.Interaction, w http.ResponseWriter) {

	// redirect captcha preharvest button interactions to their corresponding slash command counterparts
	action := strings.Split(interaction.Data.Custom_ID, "-")[2]
	if strings.HasPrefix(action, "active") {
		if len(strings.Split(interaction.Data.Custom_ID, "-")) > 3 {
			navigate(interaction, w)
		} else {
			preharvest.Active(interaction, func() (id *disgo.Snowflake) {
				if len(strings.Split(action, ":")) != 1 {
					val := disgo.Snowflake(strings.Split(action, ":")[1])
					id = &val
				}
				return
			}(), "", w)
		}

	} else if strings.Split(action, ":")[0] == "remove" {
		preharvest.Remove(interaction, strings.Split(action, ":")[1], w)

	} else if strings.Split(action, ":")[0] == "schedule" {
		preharvest.Schedule(interaction, strings.Split(action, ":")[1], w)

	} else if strings.Split(action, ":")[0] == "patch" {
		preharvest.UpdateState(interaction, strings.Split(action, ":")[1], consts.Running, w)

	} else if strings.Split(action, ":")[0] == "stop" {
		preharvest.UpdateState(interaction, strings.Split(action, ":")[1], consts.Stopped, w)

	} else if strings.Split(action, ":")[0] == "restart" {
		preharvest.UpdateState(interaction, strings.Split(action, ":")[1], consts.Running, w)
	}
}
