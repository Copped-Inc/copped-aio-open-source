package preharvest

import (
	"errors"
	"github.com/Copped-Inc/aio-types/branding"
	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/helper"
	"github.com/Copped-Inc/aio-types/secrets"
	"github.com/infinitare/disgo"
	"net/http"
	"service-api/handler/discord/interactions/vars"
)

func Remove(interaction *disgo.Interaction, taskID string, w http.ResponseWriter) {
	var (
		response = disgo.InteractionResponse{Type: disgo.CHANNEL_MESSAGE_WITH_SOURCE, Data: disgo.InteractionCallbackData{Flags: disgo.EPHEMERAL, Embeds: []disgo.Embed{{Color: branding.Red, Footer: &disgo.EmbedFooter{Text: "Click below to go back to the list of preharvest tasks. Alternatively try deleting a preharvest again, by clicking the mention above."}, Fields: []disgo.EmbedField{{Name: "\u200b", Value: "</captcha preharvest remove:1077626991227961355>"}}}}, Components: &[]disgo.Component{{Type: disgo.Action_Row, Components: []disgo.Component{{Type: disgo.Button, Style: 1, Custom_ID: "captcha-preharvest-active", Emoji: &disgo.Emoji{ID: "1077610684143108208", Name: "white_list"}}}}}}}
		embed    = &response.Data.Embeds[0]
	)

	// placeholder ID returned by autcomplete when there are no captcha preharvest tasks
	if taskID == NoContent {
		NoTasks(interaction, nil, w)
		return
	}

	// application commands require an additional step to double check whether the user is allowed to perform this ation on the task specified
	if interaction.Type == disgo.APPLICATION_COMMAND && !isAuthorized(interaction, taskID, w, response, "delete") {
		return
	}

	// delete preharvest task
	req, err := http.NewRequest(http.MethodDelete, helper.ActiveData+"/captcha/preharvest/"+taskID, nil)
	if err != nil {
		console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
		return
	}

	req.Header.Set("Password", secrets.API_Admin_PW)

	res, err := (&http.Client{}).Do(req)
	if err != nil {
		console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
		return
	} else if res.StatusCode != http.StatusNoContent && res.StatusCode == http.StatusNotFound {
		console.ErrorRequest(w, nil, errors.New("request to "+res.Request.URL.String()+" failed with response "+res.Status), http.StatusInternalServerError)
		return
	}

	// prompt the user to retry action or go back to list of preharvest tasks if the task specified couldn't be found
	if res.StatusCode == http.StatusNotFound {
		NotFound(interaction, response, taskID, w)
		return
	}

	embed.Color = branding.Green
	embed.Title = "Deletion successful!"
	embed.Description = "The preharvest task **` " + taskID + " `** was successfully deleted."
	embed.Fields = []disgo.EmbedField{}
	embed.Footer.Text = "Click below to go back to the list of preharvest tasks."

	if interaction.Type == disgo.MESSAGE_COMPONENT {
		response.Type = disgo.UPDATE_MESSAGE
	}

	if err := vars.Respond(response, w); err != nil {
		console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
	}
}
