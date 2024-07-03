package newsletter

import (
	"github.com/Copped-Inc/aio-types/console"
	"net/http"
	"service-api/handler/discord/interactions/vars"

	"github.com/Copped-Inc/aio-types/branding"
	"github.com/infinitare/disgo"
)

func Basic(interaction *disgo.Interaction, w http.ResponseWriter) {

	if vars.Add == nil {
		vars.CacheHandler(interaction)
	} else {
		vars.Add <- interaction
	}

	response := disgo.InteractionResponse{Type: disgo.CHANNEL_MESSAGE_WITH_SOURCE, Data: disgo.InteractionCallbackData{Flags: disgo.EPHEMERAL, Embeds: []disgo.Embed{{Title: "Newsletter", Color: branding.Green, Description: "Do you want the newsletter to include a button too?", Footer: &disgo.EmbedFooter{Text: "Note that this interaction will be automatically cancelled after a prolonged period of inactivity, resulting in the loss of the data already entered."}}}}}
	response.Data.Components = &[]disgo.Component{{Type: disgo.Action_Row, Components: []disgo.Component{{Type: disgo.Button, Style: 3, Emoji: &disgo.Emoji{Animated: false, ID: "1041703241064394776", Name: "white_checkbox"}, Custom_ID: "newsletter-button:yes-" + string(interaction.ID)}, {Type: disgo.Button, Style: 4, Emoji: &disgo.Emoji{Animated: false, ID: "1041703355023642694", Name: "white_close"}, Custom_ID: "newsletter-button:no-" + string(interaction.ID)}, {Type: disgo.Button, Style: 1, Emoji: &disgo.Emoji{Animated: false, ID: "1041707451667468419", Name: "white_discard"}, Label: "discard mail", Custom_ID: "newsletter-cancel-" + string(interaction.ID)}}}}

	if err := vars.Respond(response, w); err != nil {
		console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
	}
}
