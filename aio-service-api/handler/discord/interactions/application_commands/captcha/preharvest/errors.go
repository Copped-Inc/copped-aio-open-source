package preharvest

import (
	"github.com/Copped-Inc/aio-types/branding"
	"github.com/Copped-Inc/aio-types/console"
	"github.com/infinitare/disgo"
	"net/http"
	"service-api/handler/discord/interactions/vars"
)

const NoContent = "no captcha preharvest tasks"

func NotFound(interaction *disgo.Interaction, response disgo.InteractionResponse, taskID string, w http.ResponseWriter) {
	embed := &response.Data.Embeds[0]
	embed.Color = branding.Red
	embed.Title = "Preharvest task not found!"
	embed.Description = "The preharvest task specified couldn't be found! **` " + taskID + " `** is an invalid preharvest task ID. Make sure to provide a valid value next time."

	if err := vars.Respond(response, w); err != nil {
		console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
	}
}

func InvalidPermissions(interaction *disgo.Interaction, response disgo.InteractionResponse, w http.ResponseWriter) {
	embed := &response.Data.Embeds[0]
	embed.Color = branding.Red
	embed.Title = "Insufficient permissions!"

	if err := vars.Respond(response, w); err != nil {
		console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
	}
}

func NoTasks(interaction *disgo.Interaction, userID *disgo.Snowflake, w http.ResponseWriter) {
	foreign := false
	response := disgo.InteractionResponse{Type: func() disgo.InteractionCallbackType {
		if interaction.Type != disgo.APPLICATION_COMMAND {
			return disgo.UPDATE_MESSAGE
		}
		return disgo.CHANNEL_MESSAGE_WITH_SOURCE
	}(), Data: disgo.InteractionCallbackData{Flags: disgo.EPHEMERAL, Embeds: []disgo.Embed{
		{
			Color: branding.Orange,
			Title: "No captcha preharvest tasks",
			Description: func() string {
				if userID == nil {
					return "There are no preharvest tasks yet. Create one to get started."
				} else if *userID == interaction.Member.User.ID {
					return "You have no preharvest tasks yet. Create one to get started."
				}
				foreign = true
				return "User <@" + string(*userID) + "> has no preharvest tasks yet."
			}(),
		}}}}
	if !foreign {
		response.Data.Embeds[0].Fields = []disgo.EmbedField{{Name: "\u200b", Value: "</captcha preharvest new:1077626991227961355>"}}
		response.Data.Embeds[0].Footer = &disgo.EmbedFooter{Text: "Click above to create a new preharvest task."}
	}

	if err := vars.Respond(response, w); err != nil {
		console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
	}
}
