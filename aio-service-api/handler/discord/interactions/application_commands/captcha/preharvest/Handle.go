package preharvest

import (
	"net/http"

	"github.com/Copped-Inc/aio-types/captcha/preharvest"
	"github.com/infinitare/disgo"
)

func Handle(interaction *disgo.Interaction, w http.ResponseWriter) {
	switch interaction.Data.Options[0].Options[0].Name {
	case "new":
		new(interaction, w)
	case "active":
		Active(interaction, func() (id *disgo.Snowflake) {
			if inputs := interaction.Data.Options[0].Options[0].Options; len(inputs) != 0 {
				for _, input := range inputs {
					if input.Name == "user" {
						snowflake := disgo.Snowflake(input.Value.(string))
						id = &snowflake
						break
					}
				}
			}
			return
		}(), func() (taskID string) {
			if inputs := interaction.Data.Options[0].Options[0].Options; len(inputs) != 0 {
				for _, input := range inputs {
					if input.Name == "id" {
						taskID = input.Value.(string)
						break
					}
				}
			}
			return
		}(), w)
	case "stop":
		UpdateState(interaction, interaction.Data.Options[0].Options[0].Options[0].Value.(string), preharvest.Stopped, w)
	case "restart":
		UpdateState(interaction, interaction.Data.Options[0].Options[0].Options[0].Value.(string), preharvest.Running, w)
	case "remove":
		Remove(interaction, interaction.Data.Options[0].Options[0].Options[0].Value.(string), w)
	}
}
