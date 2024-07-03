package vars

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/helper"

	"github.com/Copped-Inc/aio-types/branding"
	"github.com/Copped-Inc/aio-types/captcha/preharvest"
	"github.com/Copped-Inc/aio-types/discord"
	"github.com/Copped-Inc/aio-types/secrets"
	"github.com/infinitare/disgo"
)

var (
	Add    chan *disgo.Interaction
	Remove = make(chan disgo.Snowflake)
	ACK    = make(chan disgo.Snowflake)
	Cache  = make(map[disgo.Snowflake]*disgo.Interaction)
)

func InternalRequest(data []byte, path string) (*http.Response, error) {

	req, err := http.NewRequest(http.MethodPost, helper.ActiveData+"/"+path, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Password", secrets.API_Admin_PW)
	req.Header.Set("Content-Type", "application/json")

	res, err := (&http.Client{}).Do(req)
	if err == nil {
		if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
			err = errors.New("request to " + res.Request.URL.String() + " failed with response " + res.Status)
		}
	}

	return res, err
}

func Respond(response disgo.InteractionResponse, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

func CacheHandler(interaction *disgo.Interaction) {

	Add = make(chan *disgo.Interaction, 3)

	go func() {
		var interacted = make(map[disgo.Snowflake]chan bool, 1)

		defer func() {
			close(Add)
			Add = nil
		}()

		for {
			select {
			case interaction := <-Add:
				Cache[interaction.ID] = interaction
				interacted[interaction.ID] = make(chan bool, 1)

				go func() { // handler for a single interaction
					for {
						select {
						case <-(time.NewTimer(14 * time.Minute)).C:
							Remove <- interaction.ID

							data, err := json.Marshal(disgo.InteractionCallbackData{Embeds: []disgo.Embed{{Title: "Interaction timed out!", Color: branding.Yellow, Description: "No interaction with the ongoing interaction occured for an extended amount of time. The interaction was automatically cancelled due to this reason.", Fields: []disgo.EmbedField{{Name: "\u200b", Value: "</newsletter:1041752868216111115>"}}, Footer: &disgo.EmbedFooter{Text: "In case you didn't intent to cancel the interaction, feel free to give it another try, by clicking the mention above. We apologize for the inconvenience."}}}, Components: &[]disgo.Component{}})
							if err != nil {
								console.Log(err)
								break
							}

							req, err := http.NewRequest(http.MethodPatch, "https://discord.com/api/v"+discord.API_Version+"/webhooks/"+string(interaction.Application_ID)+"/"+string(interaction.Token)+"/messages/@original", bytes.NewBuffer(data))
							if err != nil {
								console.Log(err)
								break
							}

							req.Header.Set("Content-type", "application/json")
							req.Header.Add("Authorization", discord.Bearer)

							_, err = (&http.Client{}).Do(req)
							if err != nil {
								console.Log(err)
							}

						case <-interacted[interaction.ID]:
							delete(interacted, interaction.ID)
							return
						}
					}
				}()

			case id := <-Remove:
				if interacted, ok := interacted[id]; ok {
					interacted <- true
					close(interacted)
				}
				delete(Cache, id)
				if len(Cache) == 0 {
					return
				}

			case id := <-ACK:
				interacted[id] <- true
				close(interacted[id])
				if len(Cache) == 0 {
					return
				}
			}
		}
	}()

	Add <- interaction

}

func CacheToMail(id disgo.Snowflake, w http.ResponseWriter) {
	type mail struct {
		Title       string `json:"title"`
		Subtitle    string `json:"subtitle"`
		Text        string `json:"text"`
		Button      bool   `json:"button"`
		ButtonUrl   string `json:"button_url"`
		ButtonText  string `json:"button_text"`
		BelowButton string `json:"below_button"`
	}

	interaction := Cache[id]
	Remove <- id
	inputs := interaction.Data.Components
	data := mail{Title: inputs[0].Components[0].Value}
	embed := disgo.Embed{Title: data.Title, Color: branding.Orange}
	button := disgo.Component{Type: disgo.Button, Style: 5}

	for _, component := range inputs {
		component = component.Components[0]
		val := component.Value

		switch component.Custom_ID {
		case "subtitle":
			data.Subtitle = val
			embed.Author = &disgo.EmbedAuthor{Name: val}
		case "text":
			data.Text = val
			embed.Description = val
		case "label":
			data.ButtonText = val
			data.Button = true
			button.Label = val
			embed.Description += "\n\u200b"
		case "button subtitle":
			data.BelowButton = val
			embed.Footer = &disgo.EmbedFooter{Text: "note, this will be written underneath the button:\n" + val}
		case "url":
			data.ButtonUrl = val
			button.URL = val
		}
	}

	payload, err := json.Marshal(data)
	if err != nil {
		console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
		return
	}

	_, err = InternalRequest(payload, "newsletter")
	if err != nil {
		console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
		return
	}

	response := disgo.InteractionResponse{Type: disgo.UPDATE_MESSAGE, Data: disgo.InteractionCallbackData{Flags: disgo.EPHEMERAL, Components: &[]disgo.Component{}, Embeds: []disgo.Embed{{Title: "Newsletter successfully sent!", Color: branding.Green, Description: "Check out the preview below:"}, embed}}}
	if data.Button {
		response.Type = disgo.CHANNEL_MESSAGE_WITH_SOURCE
		response.Data.Components = &[]disgo.Component{{Type: disgo.Action_Row, Components: []disgo.Component{button}}}
	}

	if err := Respond(response, w); err != nil {
		console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
	}
}

func TaskToEmbed(task preharvest.Task) disgo.Embed {
	return disgo.Embed{
		Author: &disgo.EmbedAuthor{URL: task.Site.GetData().URL, Name: task.Site.GetData().Name},
		Title: func() (schedule string) {
			if task.Routine {
				schedule = "scheduled "
			}
			return
		}() + "preharvest task",
		Description: "This preharvest task will run " + func() string {
			if task.Schedule != "" {
				return task.Schedule + " starting"
			}
			return "once"
		}() + " " + func() string {
			seconds := strconv.FormatInt(task.Date.Unix(), 10)
			if task.Date.Day() != time.Now().Day() {
				return "at <t:" + seconds + ">."
			}
			return "<t:" + seconds + ":R>."
		}(),
		Footer: &disgo.EmbedFooter{Text: "ID: " + task.ID},
		Color:  branding.Yellow,
	}
}

func ActiveTaskList(userSpecific bool, response *disgo.InteractionResponse, tasks []preharvest.Task) {
	embed := &response.Data.Embeds[0]

	// append tasks as embed fields
	for index, task := range tasks {
		if index > 24 {
			break
		}

		name := strconv.Itoa(index+1) + " — " + task.Site.GetData().Name

		embed.Fields = append(embed.Fields, disgo.EmbedField{Name: "<:yellow_hashtag:1077598907950972968> " + name,
			Value: "> " + func() (user string) {
				if !userSpecific {
					user = "from <@" + string(task.User_ID) + ">\n> "
				}
				return
			}() + func() string {
				if task.State == preharvest.Running {
					return "<:green_running:1077601547397124167> running, next execution " + func() string {
						seconds := strconv.FormatInt(task.Date.Unix(), 10)
						if task.Date.Day() != time.Now().Day() {
							return "at <t:" + seconds + ">"
						}
						return "<t:" + seconds + ":R>"
					}()
				}
				return "<:red_stopped:1077603145561165885> stopped"
			}() + "\n> " + func() (routine string) {
				if task.Routine {
					routine = "<:black_schedule:1077603691848290305> scheduled"
				}
				return
			}() + func() (uses string) {
				if task.Routine && task.State == preharvest.Running {
					uses = ", uses remaining " + strconv.Itoa(task.Uses)
				}
				return
			}()},
		)

		(*response.Data.Components)[0].Components[0].Options = append((*response.Data.Components)[0].Components[0].Options, disgo.SelectMenuOption{Label: name, Value: task.ID, Emoji: &disgo.Emoji{ID: "1077598907950972968", Name: "yellow_hashtag"}})
	}
}
