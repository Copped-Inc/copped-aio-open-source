package instance

import (
	"encoding/json"
	"github.com/Copped-Inc/aio-types/console"
	"net/http"
	"service-api/handler/discord/interactions/vars"
	"strings"

	"github.com/Copped-Inc/aio-types/branding"
	"github.com/infinitare/disgo"
)

func Update(interaction *disgo.Interaction, w http.ResponseWriter) {
	text := interaction.Data.Components[0].Components[0].Value
	var payload []string

	if len(strings.Split(text, ";"))+len(strings.Split(text, "\n")) == 2 {
		payload = []string{text}
	} else if len(strings.Split(text, ";")) > 1 {
		payload = strings.Split(text, ";")
	} else {
		payload = strings.Split(text, "\n")
	}

	data, err := json.Marshal(struct {
		Changelog []string `json:"changelog"`
	}{payload})
	if err != nil {
		console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
		return
	}

	_, err = vars.InternalRequest(data, "instance/update/client")
	if err != nil {
		console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
		return
	}

	response := disgo.InteractionResponse{Type: disgo.CHANNEL_MESSAGE_WITH_SOURCE, Data: disgo.InteractionCallbackData{Flags: disgo.EPHEMERAL, Embeds: []disgo.Embed{{Title: "Update notice posted successfully!", Color: branding.Green, Description: "Check out <#1006491954080653312> for the complete update notice including the following changelog:"}}}}
	for _, v := range payload {
		response.Data.Embeds[0].Description += "\n• " + v
	}

	if err := vars.Respond(response, w); err != nil {
		console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
	}
}
