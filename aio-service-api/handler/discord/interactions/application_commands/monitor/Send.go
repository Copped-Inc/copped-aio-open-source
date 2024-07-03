package monitor

import (
	"encoding/json"
	"github.com/Copped-Inc/aio-types/branding"
	"github.com/Copped-Inc/aio-types/console"
	"github.com/infinitare/disgo"
	"net/http"
	"regexp"
	"service-api/handler/discord/interactions/vars"
	"strings"
)

var (
	regx       = regexp.MustCompile(`^\[Link\]\(https:\/\/(?:(?:\w|-)+\.)*(?:\w|-)+\.(?:\w|-)+\/products(?:\/(?:\w|-)+)+\)$`)
	urlPattern = regexp.MustCompile(`https:\/\/(?:(?:\w|-)+\.)*(?:\w|-)+\.(?:\w|-)+\/products(?:\/(?:\w|-)+)+`)
)

func Send(interaction *disgo.Interaction, w http.ResponseWriter) {
	var (
		response         = disgo.InteractionResponse{Type: disgo.CHANNEL_MESSAGE_WITH_SOURCE, Data: disgo.InteractionCallbackData{Flags: disgo.EPHEMERAL, Embeds: []disgo.Embed{{Color: branding.Green}}}}
		handle, sku, url string
		state            = func() int {
			if interaction.Data.Type == disgo.ApplicationCommandType_MESSAGE {
				return 1
			} else if interaction.Data.Options[0].Name == "blacklist" {
				return 1
			}
			return 0
		}()
	)

	// validate message commands
	if interaction.Data.Type == disgo.ApplicationCommandType_MESSAGE {
		valid := false

		if interaction.Channel_ID != "1022857400170070116" {
			response.Data.Embeds = []disgo.Embed{{Title: "Invalid Channel", Color: branding.Red, Description: "Items can only be blacklisted inside <#1022857400170070116>!"}}
		} else {
			if embeds := interaction.Data.Resolved.Messages[interaction.Data.Target_ID].Embeds; len(embeds) != 0 && interaction.Data.Resolved.Messages[interaction.Data.Target_ID].Webhook_ID == "1022859099962101832" {
				if fields := embeds[0].Fields; len(fields) != 0 {
					if value := fields[0].Value; regx.MatchString(value) {
						url = urlPattern.FindString(value)
						handle = strings.Split(url, "/")[len(strings.Split(url, "/"))-1]
						response.Data.Embeds[0].Description = "The item **`` " + embeds[0].Title + " ``** was classified as unprofitable and therefore added to the blacklisted products."
						valid = true
					}
				}
			}
		}

		if !valid {
			response.Data.Embeds = []disgo.Embed{{Title: "Invalid Message", Color: branding.Red, Description: "Can only blacklist items when provided a webhook message of <@1022859099962101832> monitor including a link to the product's page."}}
			if err := vars.Respond(response, w); err != nil {
				console.ErrorLog(err)
			}
			return
		}

		response.Data.Components = &[]disgo.Component{{Type: disgo.Action_Row, Components: []disgo.Component{{Type: disgo.Button, Style: 5, URL: "discord://-/channels/" + string(interaction.Guild_ID) + "/" + string(interaction.Channel_ID) + "/" + string(interaction.Data.Target_ID), Label: "message origin"}}}}

	} else
	// otherwise populate values with the input made to the slash command
	{
		inputs := interaction.Data.Options[0].Options
		if len(inputs) == 0 {
			response.Data.Embeds = []disgo.Embed{{Title: "Missing Input", Color: branding.Red, Description: "Must either provide a sku or the item's handle."}}
			if err := vars.Respond(response, w); err != nil {
				console.ErrorLog(err)
			}
			return
		}

		for _, input := range inputs {
			if input.Name == "sku" {
				sku = input.Value.(string)
			} else {
				handle = input.Value.(string)
			}
		}

		response.Data.Embeds[0].Description = func() string {
			if state == 0 {
				return "The item was classified as profitable and therefore added to the whitelisted products."
			}
			return "The item was classified as unprofitable and therefore added to the blacklisted products."
		}()
		response.Data.Embeds[0].Footer = &disgo.EmbedFooter{Text: "SKU | " + sku}
	}

	data, err := json.Marshal(struct {
		Handle string `json:"handle,omitempty"`
		State  int    `json:"state"`
	}{handle, state})
	if err != nil {
		console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
		return
	}

	_, err = vars.InternalRequest(data, "monitor/product/"+func() string {
		if sku != "" {
			return sku
		}
		return handle
	}())
	if err != nil {
		console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
		return
	}

	response.Data.Embeds[0].Title = "Item successfully added to " + func() string {
		if state == 0 {
			return "white"
		}
		return "black"
	}() + "list!"

	if err := vars.Respond(response, w); err != nil {
		console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
	}
}
