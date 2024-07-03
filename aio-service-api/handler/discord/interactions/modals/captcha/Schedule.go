package captcha

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/helper"
	"net/http"
	"service-api/handler/discord/interactions/vars"

	"github.com/Copped-Inc/aio-types/branding"
	"github.com/Copped-Inc/aio-types/captcha/preharvest"
	"github.com/Copped-Inc/aio-types/secrets"

	"github.com/infinitare/disgo"
)

func PreharvestSchedule(interaction *disgo.Interaction, w http.ResponseWriter) {

	response := disgo.InteractionResponse{Type: func() disgo.InteractionCallbackType {
		if interaction.Message != nil {
			return disgo.UPDATE_MESSAGE
		}
		return disgo.CHANNEL_MESSAGE_WITH_SOURCE
	}(), Data: disgo.InteractionCallbackData{Flags: disgo.EPHEMERAL}}

	var schedule, taskID string

	if len(interaction.Data.Components) > 1 {
		schedule = interaction.Data.Components[1].Components[0].Value
		taskID = interaction.Data.Components[1].Components[0].Custom_ID
	} else {
		schedule = interaction.Data.Components[0].Components[0].Value
		taskID = interaction.Data.Components[0].Components[0].Custom_ID
	}

	// check if schedule provided matches regex pattern
	// if not, send an error with the option to delete or edit the task optionally
	if !preharvest.Schedule_Pattern.MatchString(schedule) {
		embed := disgo.Embed{
			Color:       branding.Red,
			Title:       "Invalid schedule!",
			Description: "The schedule provided by you (` " + interaction.Data.Components[0].Value + " `) doesn't match the pattern expected. Please make sure to follow the instructions given in the pop-up next time.",
			Footer:      &disgo.EmbedFooter{Text: "Note that preharvest tasks without a routine will be deleted after running once. To prevent this, click below to retry scheduling this task. Alternatively, go back to the list of preharvest tasks or delete the one you just created preharvest task."},
		}

		response.Data.Embeds = []disgo.Embed{embed}
		response.Data.Components = &[]disgo.Component{{Type: disgo.Action_Row, Components: []disgo.Component{{Label: "edit task", Type: disgo.Button, Style: 3, Custom_ID: "captcha-preharvest-schedule:" + taskID, Emoji: &disgo.Emoji{ID: "1076270745535136024", Name: "black_clock_add"}}, {Type: disgo.Button, Style: 1, Custom_ID: "captcha-preharvest-active", Emoji: &disgo.Emoji{ID: "1077610684143108208", Name: "white_list"}}, {Type: disgo.Button, Style: 1, Emoji: &disgo.Emoji{ID: "1041707451667468419", Name: "white_discard"}, Label: "delete task", Custom_ID: "captcha-preharvest-remove:" + taskID}}}}

		if err := vars.Respond(response, w); err != nil {
			console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
		}

		return
	}

	// for valid schedules, update the preharvest task, determined by the custom id of the modal
	payload, err := json.Marshal(preharvest.Task_Edit{Schedule: schedule, State: preharvest.Running})
	if err != nil {
		console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest(http.MethodPatch, helper.ActiveData+"/captcha/preharvest/"+taskID, bytes.NewBuffer(payload))
	if err != nil {
		console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
		return
	}

	req.Header.Set("Password", secrets.API_Admin_PW)
	req.Header.Set("Content-Type", "application/json")

	res, err := (&http.Client{}).Do(req)
	if err != nil {
		console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
		return
	} else if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNotFound {
		console.ErrorRequest(w, nil, errors.New("request to "+res.Request.URL.String()+" failed with response "+res.Status), http.StatusInternalServerError)
		return
	}

	// parse the response payload
	var new_task preharvest.Task
	if err := json.NewDecoder(res.Body).Decode(&new_task); err != nil {
		console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
		return
	}

	embed := disgo.Embed{
		Color:       branding.Green,
		Title:       "Preharvest task created successfully!",
		Description: "A new scheduled preharvest task was created successfully! For more detailed information see the embed below.",
		Footer:      &disgo.EmbedFooter{Text: "Click the button at the very bottom to go to the list of active preharvest task."},
	}

	// respond with 2 embeds, one indicating the success of the action, the other one providing information about the preharvest task
	// includes a button to go to the list of active preharvest tasks
	response.Data.Components = &[]disgo.Component{{Type: disgo.Action_Row, Components: []disgo.Component{{Type: disgo.Button, Style: 1, Custom_ID: "captcha-preharvest-active", Emoji: &disgo.Emoji{ID: "1077610684143108208", Name: "white_list"}}}}}
	response.Data.Embeds = append(response.Data.Embeds, embed, vars.TaskToEmbed(new_task))

	if err := vars.Respond(response, w); err != nil {
		console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
	}
}
