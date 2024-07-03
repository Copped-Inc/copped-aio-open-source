package newsletter

import (
	"embed"
	"encoding/json"
	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/discord"
	"net/http"
	"service-api/handler/discord/interactions/vars"
	"strings"

	"github.com/Copped-Inc/aio-types/branding"
	"github.com/infinitare/disgo"
)

//go:embed callback.json
var callback embed.FS

func Buttons(interaction *disgo.Interaction, w http.ResponseWriter) {
	original_interaction_id := disgo.Snowflake(strings.Split(interaction.Data.Custom_ID, "-")[2])

	if _, ok := vars.Cache[original_interaction_id]; !ok {
		if err := vars.Respond(disgo.InteractionResponse{Type: disgo.UPDATE_MESSAGE, Data: disgo.InteractionCallbackData{Flags: disgo.EPHEMERAL, Components: &[]disgo.Component{}, Embeds: []disgo.Embed{{Title: "Interaction failed", Color: branding.Red, Description: "An internal server error occured during this interaction, please try again.", Fields: []disgo.EmbedField{{Name: "\u200b", Value: "</newsletter:1041752868216111115>"}}, Footer: &disgo.EmbedFooter{Text: "You can give it another try by clicking the mention above.\nWe apologize for the inconvenience."}}}}}, w); err != nil {
			console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
		}
		return
	}

	switch strings.Split(strings.Split(interaction.Data.Custom_ID, "-")[1], ":")[0] {
	case "cancel":
		vars.Remove <- original_interaction_id

		if err := vars.Respond(disgo.InteractionResponse{Type: disgo.UPDATE_MESSAGE, Data: disgo.InteractionCallbackData{Flags: disgo.EPHEMERAL, Components: &[]disgo.Component{}, Embeds: []disgo.Embed{{Title: "Interaction cancelled", Color: branding.Yellow, Description: "The current newsletter creation interaction was cancelled. Previously submitted data were discarded.", Fields: []disgo.EmbedField{{Name: "\u200b", Value: "</newsletter:1041752868216111115>"}}, Footer: &disgo.EmbedFooter{Text: "In case the interaction was terminated on accident, you can give it another try by clicking the mention above. Better luck next time."}}}}}, w); err != nil {
			console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
		}

	case "button":
		if strings.Split(strings.Split(interaction.Data.Custom_ID, "-")[1], ":")[1] == "yes" {
			vars.ACK <- original_interaction_id
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

			response.Data.Custom_ID += string(original_interaction_id)

			if err := vars.Respond(response, w); err != nil {
				console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
				return
			}

			req, err := http.NewRequest(http.MethodDelete, "https://discord.com/api/v"+discord.API_Version+"/webhooks/"+string(vars.Cache[original_interaction_id].Application_ID)+"/"+string(vars.Cache[original_interaction_id].Token)+"/messages/@original", nil)
			if err != nil {
				console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
				return
			}

			req.Header.Set("Content-type", "application/json")
			req.Header.Add("Authorization", discord.Bearer)

			res, err := (&http.Client{}).Do(req)
			if err != nil || res.StatusCode != http.StatusNoContent {
				console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
			}

		} else {
			vars.CacheToMail(original_interaction_id, w)
		}
	}
}
