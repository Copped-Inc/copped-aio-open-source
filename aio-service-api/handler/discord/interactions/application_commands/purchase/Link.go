package purchase

import (
	"encoding/json"
	"io"
	"net/http"
	"service-api/handler/discord/interactions/vars"
	"strconv"

	"github.com/Copped-Inc/aio-types/branding"
	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/subscriptions"
	"github.com/infinitare/disgo"
)

func Link(interaction *disgo.Interaction, w http.ResponseWriter) {

	input := interaction.Data.Options[0].Options[0].Options

	plan := subscriptions.Plan(input[0].Value.(float64))
	stock := int(input[1].Value.(float64))
	limit := int(input[2].Value.(float64))

	data, err := json.Marshal(struct {
		Plan          subscriptions.Plan `json:"plan"`
		Stock         int                `json:"stock"`
		InstanceLimit int                `json:"instance_limit"`
	}{plan, stock, limit})
	if err != nil {
		console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
		return
	}

	res, err := vars.InternalRequest(data, "purchase")
	if err != nil {
		console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
		return
	}

	data, err = io.ReadAll(res.Body)
	if err != nil {
		console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
		return
	}

	var payload struct {
		Link string `json:"link"`
	}

	err = json.Unmarshal(data, &payload)
	if err != nil {
		console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
		return
	}

	response := disgo.InteractionResponse{Type: disgo.CHANNEL_MESSAGE_WITH_SOURCE, Data: disgo.InteractionCallbackData{Flags: disgo.EPHEMERAL, Embeds: []disgo.Embed{{Title: "Purchase link generated successfully!", Color: branding.Green, Description: "A new purchase link with the following conditions has been generated successfully:"}}}}
	response.Data.Embeds[0].Fields = []disgo.EmbedField{{Name: "Plan", Value: plan.GetData().Name, Inline: true}, {Name: "Stock", Value: strconv.Itoa(stock), Inline: true}, {Name: "Instance Limit", Value: strconv.Itoa(limit), Inline: true}, {Name: "\u200b", Value: "<:white_copy_link:1041714302370988132> [purchase link](" + payload.Link + ")", Inline: false}}

	if err := vars.Respond(response, w); err != nil {
		console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
	}
}
