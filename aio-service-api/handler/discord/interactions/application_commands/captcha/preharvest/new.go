package preharvest

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/Copped-Inc/aio-types/branding"
	"github.com/Copped-Inc/aio-types/captcha/preharvest"
	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/discord"
	"github.com/Copped-Inc/aio-types/helper"
	"github.com/Copped-Inc/aio-types/modules"
	"github.com/Copped-Inc/aio-types/secrets"
	"github.com/infinitare/disgo"
	"net/http"
	"service-api/handler/discord/interactions/vars"
	"time"
)

func new(interaction *disgo.Interaction, w http.ResponseWriter) {
	defer func() {
		err := recover()
		if err != nil {
			console.ErrorLog(err.(error))
		}
	}()

	var (
		captchaResult chan error
		inputs        = interaction.Data.Options[0].Options[0].Options
		response      = disgo.InteractionResponse{Type: disgo.CHANNEL_MESSAGE_WITH_SOURCE, Data: disgo.InteractionCallbackData{Flags: disgo.EPHEMERAL}}
		embed         = disgo.Embed{Color: branding.Green}
	)

	// check if site value is valid, else return error message
	if !func() bool {
		for _, site := range modules.Sites {
			if inputs[0].Value.(string) == string(site) {
				return true
			}
		}
		return false
	}() {
		embed.Color = branding.Red
		embed.Title = "Invalid site"
		embed.Description = "**` " + inputs[0].Value.(string) + " `** is not a valid site! If you aren't sure as to what sites are supported for captcha preharvesting, consider picking on of the autocomplete suggestions provided instead next time."
		embed.Fields = []disgo.EmbedField{{Name: "\u200b", Value: "</captcha preharvest new:1077626991227961355>"}}
		embed.Footer = &disgo.EmbedFooter{Text: "Click above to give it another try."}

		response.Data.Embeds = append(response.Data.Embeds, embed)

		if err := vars.Respond(response, w); err != nil {
			console.ErrorLog(err)
		}

		return
	}

	// populate task preharvest create payload with the input(s) provided
	task := preharvest.Task_Create{Site: modules.Site(inputs[0].Value.(string))}

	if len(inputs) > 2 {
		// since delay accepts floating point numbers, calculate trigger date as current time + the delay specified in minutes converted to seconds
		task.Date = time.Now().Add(time.Second * time.Duration(inputs[1].Value.(float64)*60))
		task.Routine = inputs[2].Value.(bool)
	} else if len(inputs) > 1 {
		switch inputs[1].Name {
		case "routine":
			task.Routine = inputs[1].Value.(bool)
		case "delay":
			task.Date = time.Now().Add(time.Second * time.Duration(inputs[1].Value.(float64)*60))
		}
	}

	// in case no delay was specified, instantly start generating captchas
	if task.Date.IsZero() {
		task.Date = time.Now()
		captchaResult = make(chan error, 1)
		go func() {
			captchaResult <- generateCaptcha(task.Site)
		}()
	}

	// awaits result of generating captchas and responds to the interaction

	// if the task shouldn't be repeated and executed instantly, creating a preharvest task isn't necessary
	if captchaResult != nil && !task.Routine {

		// to avoid interaction timeout due to slow captcha endpoint response, ACK the interaction first

		//  should work but doesn't, for reasons unknown
		if err := vars.Respond(disgo.InteractionResponse{Type: disgo.DEFERRED_CHANNEL_MESSAGE_WITH_SOURCE, Data: disgo.InteractionCallbackData{Flags: disgo.EPHEMERAL}}, w); err != nil {
			console.ErrorLog(err)
			return
		}

		// in case any actions after deferring the interaction fail, log the reason and return an error message to indicate execution failure, otherwise conclude by indicating execution success
		go func() {
			if err := sendFollowUp(captchaResult, embed, interaction, task.Site.GetData().Name); err != nil {
				console.ErrorLog(err)
			}
		}()

		return
	}

	// create a new preharvest task
	data, err := json.Marshal(task)
	if err != nil {
		console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
		return
	}

	res, err := vars.InternalRequest(data, "captcha/preharvest/user/"+string(interaction.Member.User.ID))
	if err != nil {
		console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
		return
	} else if res.StatusCode != http.StatusCreated {
		console.ErrorRequest(w, nil, errors.New("preharvest task creation failed"), http.StatusInternalServerError)
		return
	}

	var new_task preharvest.Task
	if err := json.NewDecoder(res.Body).Decode(&new_task); err != nil {
		console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
		return
	}

	// if the preharvest task shouldn't be repeated, the initial response returns the created preharvest task as well as a button to go to this user's list of active preharvest tasks
	if !task.Routine {
		embed.Title = "Preharvest task created successfully!"
		embed.Description = "A new preharvest task was created successfully! For more detailed information see the embed below."
		embed.Footer = &disgo.EmbedFooter{Text: "Click the button at the very bottom to go to the list of active preharvest tasks."}
		response.Data.Embeds = append(response.Data.Embeds, embed, vars.TaskToEmbed(new_task))
		response.Data.Components = &[]disgo.Component{{Type: disgo.Action_Row, Components: []disgo.Component{{Type: disgo.Button, Style: 1, Custom_ID: "captcha-preharvest-active", Emoji: &disgo.Emoji{ID: "1077610684143108208", Name: "white_list"}}}}}

		if err := vars.Respond(response, w); err != nil {
			console.ErrorLog(err)
		}

		return
	}

	// otherwise, a modal will be sent as follow up to define a schedule for the preharvest task to follow
	Schedule(interaction, new_task.ID, w)
}

func generateCaptcha(site modules.Site) error {
	req, err := http.NewRequest(http.MethodGet, helper.ActiveData+"/captcha/"+string(site), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Password", secrets.API_Admin_PW)

	res, err := (&http.Client{Timeout: time.Minute * 14}).Do(req)
	if err != nil {
		return err
	} else if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNotFound {
		err = errors.New("request to " + res.Request.URL.String() + " failed with response " + res.Status)
	}

	return err
}

func sendFollowUp(capRes chan error, embed disgo.Embed, interaction *disgo.Interaction, sitename string) error {
	if err := <-capRes; err != nil {
		return err
	}

	embed.Title = "Captchas generated successfully!"
	embed.Description = "Captchas were successfully generated for " + sitename + "."

	data, err := json.Marshal(disgo.InteractionCallbackData{Embeds: []disgo.Embed{embed}, Flags: disgo.EPHEMERAL})
	if err != nil {
		return err
	}

	// in case any actions after deferring the interaction fail, log the reason and return an error message to indicate execution failure, otherwise conclude by indicating execution success
	req, err := http.NewRequest(http.MethodPost, "https://discord.com/api/v"+discord.API_Version+"/webhooks/"+string(interaction.Application_ID)+"/"+string(interaction.Token), bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	} else if res.StatusCode != http.StatusOK {
		err = errors.New("request to " + res.Request.URL.String() + " failed with response " + res.Status)
	}

	return err
}
