package newsletter

import (
	"embed"
	"encoding/json"
	"github.com/Copped-Inc/aio-types/console"
	"github.com/infinitare/disgo"
	"net/http"
	"service-api/handler/discord/interactions/vars"
)

//go:embed callback.json
var callback embed.FS

func Send(interaction *disgo.Interaction, w http.ResponseWriter) {

	response := disgo.InteractionResponse{Type: disgo.MODAL}

	data, err := callback.ReadFile("callback.json")
	if err != nil {
		console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(data, &response.Data)
	if err != nil {
		console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
		return
	}

	for interaction_id, cached_interaction := range vars.Cache {
		if cached_interaction.Member.User.ID == interaction.Member.User.ID {
			if input := interaction.Data.Options; len(input) > 0 {
				if input[0].Value.(bool) {
					textfields := *response.Data.Components
					cached := cached_interaction.Data.Components

					textfields[0].Components[0].Value = cached[0].Components[0].Value

					if len(textfields) <= 2 {
						textfields[2].Components[0].Value = cached[1].Components[0].Value
					} else {
						textfields[1].Components[0].Value = cached[1].Components[0].Value
						textfields[2].Components[0].Value = cached[2].Components[0].Value
					}
				}
			}

			vars.Remove <- interaction_id
			break
		}
	}

	if err := vars.Respond(response, w); err != nil {
		console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
	}
}
