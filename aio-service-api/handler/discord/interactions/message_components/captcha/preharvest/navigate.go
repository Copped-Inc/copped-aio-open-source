package preharvest

import (
	"net/http"
	"service-api/handler/discord/interactions/application_commands/captcha/preharvest"
	"service-api/handler/discord/interactions/vars"
	"strings"

	"github.com/Copped-Inc/aio-types/console"
	"github.com/infinitare/disgo"
)

func navigate(interaction *disgo.Interaction, w http.ResponseWriter) {
	response := disgo.InteractionResponse{Type: disgo.UPDATE_MESSAGE, Data: disgo.InteractionCallbackData{}}

	switch action := strings.Split(strings.Split(interaction.Data.Custom_ID, "-")[3], ":")[0]; action {
	case "before":
	case "after":
	case "add":
		for i, button := range interaction.Message.Components[1].Components {
			if strings.Split(button.Custom_ID, "-")[3] == action {
				button.Style = 4
				button.Custom_ID = "captcha-preharvest-active-remove-filter-user"
				interaction.Message.Components[1].Components[i] = button
				break
			}
		}

		interaction.Message.Components[0].Components[0].Disabled = true
		response.Data.Components = &[]disgo.Component{interaction.Message.Components[0], {Type: disgo.Action_Row, Components: []disgo.Component{{Type: disgo.User_Select, Custom_ID: "captcha-preharvest-active-user", Placeholder: "filter tasks by user"}}}, interaction.Message.Components[1]}

	case "remove":
		for i, button := range interaction.Message.Components[2].Components {
			if strings.Split(button.Custom_ID, "-")[3] == action {
				button.Style = 3
				button.Custom_ID = "captcha-preharvest-active-add-filter-user"
				interaction.Message.Components[2].Components[i] = button
				break
			}
		}

		interaction.Message.Components[0].Components[0].Disabled = false
		response.Data.Components = &[]disgo.Component{interaction.Message.Components[0], interaction.Message.Components[2]}

	case "refresh":
		preharvest.Active(interaction, func() (user *disgo.Snowflake) {
			if ids := strings.Split(strings.Split(interaction.Data.Custom_ID, "-")[3], ":"); len(ids) > 1 {
				id := disgo.Snowflake(ids[1])
				user = &id
			}
			return
		}(), "", w)

		return
	}

	if err := vars.Respond(response, w); err != nil {
		console.ErrorLog(err)
	}
}
