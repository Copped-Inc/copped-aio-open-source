package newsletter

import (
	"github.com/Copped-Inc/aio-types/console"
	"net/http"
	"service-api/handler/discord/interactions/vars"
	"strings"

	"github.com/Copped-Inc/aio-types/branding"
	"github.com/infinitare/disgo"
)

func Button(interaction *disgo.Interaction, w http.ResponseWriter) {
	original_interaction_id := disgo.Snowflake(strings.Split(interaction.Data.Custom_ID, "-")[1])

	if _, ok := vars.Cache[original_interaction_id]; !ok {
		if err := vars.Respond(disgo.InteractionResponse{Type: disgo.UPDATE_MESSAGE, Data: disgo.InteractionCallbackData{Flags: disgo.EPHEMERAL, Components: &[]disgo.Component{}, Embeds: []disgo.Embed{{Title: "Interaction failed", Color: branding.Red, Description: "An internal server error occured during this interaction, please try again.", Fields: []disgo.EmbedField{{Name: "\u200b", Value: "</newsletter:1041752868216111115>"}}, Footer: &disgo.EmbedFooter{Text: "You can give it another try by clicking the mention above.\nWe apologize for the inconvenience."}}}}}, w); err != nil {
			console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
		}
		return
	}

	vars.Cache[original_interaction_id].Data.Components = append(vars.Cache[original_interaction_id].Data.Components, interaction.Data.Components...)
	vars.CacheToMail(original_interaction_id, w)
}
