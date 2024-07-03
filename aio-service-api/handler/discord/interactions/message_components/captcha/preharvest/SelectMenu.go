package preharvest

import (
	"net/http"
	"service-api/handler/discord/interactions/application_commands/captcha/preharvest"
	"strings"

	"github.com/infinitare/disgo"
)

func SelectMenu(interaction *disgo.Interaction, w http.ResponseWriter) {
	action := strings.Split(interaction.Data.Custom_ID, "-")[3]

	preharvest.Active(interaction, func() (id *disgo.Snowflake) {
		if action == "user" {
			val := disgo.Snowflake(interaction.Data.Values[0])
			id = &val
		}
		return
	}(), func() (task string) {
		if action == "task" {
			task = interaction.Data.Values[0]
		}
		return
	}(), w)
}
