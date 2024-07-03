package preharvest

import (
	"encoding/json"
	"errors"
	"github.com/Copped-Inc/aio-types/captcha/preharvest"
	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/helper"
	"github.com/Copped-Inc/aio-types/secrets"
	"github.com/infinitare/disgo"
	"net/http"
	"service-api/handler/discord/interactions/vars"
	"strconv"
	"time"
)

func Active(interaction *disgo.Interaction, userID *disgo.Snowflake, taskID string, w http.ResponseWriter) {
	var (
		task     preharvest.Task
		tasks    []preharvest.Task
		reqUrl   = helper.ActiveData + "/captcha/preharvest"
		response = disgo.InteractionResponse{Type: func() disgo.InteractionCallbackType {
			if interaction.Type != disgo.APPLICATION_COMMAND {
				return disgo.UPDATE_MESSAGE
			}
			return disgo.CHANNEL_MESSAGE_WITH_SOURCE
		}(), Data: disgo.InteractionCallbackData{Flags: disgo.EPHEMERAL, Embeds: []disgo.Embed{{Title: "captcha preharvest task list"}}}}
		embed = &response.Data.Embeds[0]
	)

	// placeholder ID returned by autcomplete when there are no captcha preharvest tasks
	if taskID == NoContent {
		NoTasks(interaction, nil, w)
		return
	}

	response.Data.Components = &[]disgo.Component{{Type: disgo.Action_Row, Components: []disgo.Component{{Type: disgo.String_Select, Custom_ID: "captcha-preharvest-active-task", Placeholder: "select a specific task for more details"}}}, {Type: disgo.Action_Row}}

	// only administrators are allowed to view preharvest tasks other than their own or even all
	if !interaction.Member.Permissions.Has(disgo.ADMINISTRATOR) {
		userID = &interaction.Member.User.ID
	}

	// if a specific task ID was provided, return this task, otherwise look up the query
	if taskID != "" {
		reqUrl += "/" + taskID
	} else if userID != nil {
		reqUrl += "/user/" + string(*userID)
	}

	// request preharvest tasks
	req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
		return
	}

	// unless a specific task should be returned, limit returned results to 26
	if taskID == "" {
		query := req.URL.Query()
		query.Set("limit", "26")
		req.URL.RawQuery = query.Encode()
	}

	req.Header.Set("Password", secrets.API_Admin_PW)

	res, err := (&http.Client{}).Do(req)
	if err != nil {
		console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
		return
	} else if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNotFound && res.StatusCode != http.StatusNoContent {
		console.ErrorRequest(w, nil, errors.New("request to "+res.Request.URL.String()+" failed with response "+res.Status), http.StatusInternalServerError)
		return
	} else

	// if there are no preharvest tasks for the specified query, offer to create a new preharvest task
	if res.StatusCode == http.StatusNoContent {
		NoTasks(interaction, userID, w)
		return

		// if only a specific task was selected but not found, offer to return to list of active tasks
	} else if res.StatusCode == http.StatusNotFound {
		NotFound(interaction, response, taskID, w)
		return
	}

	if err := json.NewDecoder(res.Body).Decode(func() any {
		if taskID != "" {
			return &task
		}
		return &tasks
	}()); err != nil {
		console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
		return
	}

	// response from here on varies on whether only one task was requested or multiple
	if len(tasks) != 0 {
		/*username, icon := func() (name string, icon string) {
			if
			return "", ""
		}()
		embed.Author = &disgo.EmbedAuthor{Name: username, URL: icon}*/

		// append tasks as embed fields
		vars.ActiveTaskList(userID != nil, &response, tasks)

		if len(tasks) > 25 {
			(*response.Data.Components)[1].Components = []disgo.Component{{Type: disgo.Button, Style: 1, Emoji: &disgo.Emoji{Name: "white_bwd", ID: "1089644608234995912"}, Disabled: true, Custom_ID: "captcha-preharvest-active-before"}, {Type: disgo.Button, Style: 1, Label: "1", Disabled: true, Custom_ID: "page"}, {Type: disgo.Button, Style: 1, Emoji: &disgo.Emoji{Name: "white_fwd", ID: "1089644669476012033"}, Custom_ID: "captcha-preharvest-active-after:" + strconv.Itoa(tasks[24].Date.Second())}}
		}

		if userID == nil {
			(*response.Data.Components)[1].Components = []disgo.Component{{Type: disgo.Button, Style: 3, Emoji: &disgo.Emoji{Name: "black_filter_user", ID: "1089186470021050369"}, Custom_ID: "captcha-preharvest-active-add-filter-user"}}
		}

		(*response.Data.Components)[1].Components = append((*response.Data.Components)[1].Components, disgo.Component{Type: disgo.Button, Style: 2, Emoji: &disgo.Emoji{Name: "black_refresh", ID: "1089189421175287969"}, Custom_ID: "captcha-preharvest-active-refresh" + func() (user string) {
			if userID != nil {
				user = ":" + string(*userID)
			}
			return
		}()})

	} else {

		// only administrators can retrieve others preharvest tasks
		if !interaction.Member.Permissions.Has(disgo.ADMINISTRATOR) && task.User_ID != interaction.Member.User.ID {
			embed.Description = "You aren't allowed to view the preharvest task **` " + taskID + " `**! You can only view your own preharvest tasks."

			InvalidPermissions(interaction, response, w)
			return
		}

		response.Data.Components = &[]disgo.Component{{Type: disgo.Action_Row, Components: []disgo.Component{{Type: disgo.Button, Style: 1, Custom_ID: "captcha-preharvest-active" + func() (user string) {
			if userID != nil {
				user = ":" + string(*userID)
			}
			return
		}(), Emoji: &disgo.Emoji{ID: "1077610684143108208", Name: "white_list"}}}}}

		embed := vars.TaskToEmbed(task)
		embed.Description = func() string {
			if task.State == preharvest.Running {
				return "next execution " + func() string {
					seconds := strconv.FormatInt(task.Date.Unix(), 10)
					if task.Date.Day() != time.Now().Day() {
						return "at <t:" + seconds + ">"
					}
					return "<t:" + seconds + ":R>"
				}()
			}
			return ""
		}()

		if task.Routine && task.Schedule != "" {
			embed.Fields = append(embed.Fields, disgo.EmbedField{Name: "schedule", Value: "` " + task.Schedule + " `", Inline: true})
		} else if task.Routine {
			embed.Fields = append(embed.Fields, disgo.EmbedField{Name: "schedule", Value: "No schedule has been defined yet. Note that preharvest tasks without a routine will be deleted after running once.", Inline: false})
			(*response.Data.Components)[0].Components = append((*response.Data.Components)[0].Components, disgo.Component{Label: "edit task", Type: disgo.Button, Style: 3, Custom_ID: "captcha-preharvest-schedule:" + taskID, Emoji: &disgo.Emoji{ID: "1076270745535136024", Name: "black_clock_add"}})
		}

		if task.State == preharvest.Running {
			embed.Fields = append(embed.Fields, disgo.EmbedField{Name: "state", Value: "<:green_running:1077601547397124167> ` running `", Inline: true})
			if task.Routine {
				(*response.Data.Components)[0].Components = append((*response.Data.Components)[0].Components, disgo.Component{Type: disgo.Button, Style: 4, Emoji: &disgo.Emoji{Name: "white_stop", ID: "1089644550286487704"}, Custom_ID: "captcha-preharvest-stop:" + taskID})
			}

		} else {
			embed.Fields = append(embed.Fields, disgo.EmbedField{Name: "state", Value: "<:red_stopped:1077603145561165885> ` stopped `", Inline: true})
			(*response.Data.Components)[0].Components = append((*response.Data.Components)[0].Components, disgo.Component{Type: disgo.Button, Style: 3, Emoji: &disgo.Emoji{Name: "white_fwd", ID: "1089644669476012033"}, Disabled: task.Schedule == "", Custom_ID: "captcha-preharvest-restart:" + taskID})
		}

		(*response.Data.Components)[0].Components = append((*response.Data.Components)[0].Components, disgo.Component{Type: disgo.Button, Style: 4, Emoji: &disgo.Emoji{ID: "1041707451667468419", Name: "white_discard"}, Label: "delete task", Custom_ID: "captcha-preharvest-remove:" + taskID})

		if task.State == preharvest.Running && task.Routine {
			embed.Fields = append(embed.Fields, disgo.EmbedField{Name: "uses remaining " + strconv.Itoa(task.Uses), Value: func() string {
				if task.Uses < 7 {
					(*response.Data.Components)[0].Components = append((*response.Data.Components)[0].Components, disgo.Component{Type: disgo.Button, Style: 2, Emoji: &disgo.Emoji{Name: "black_refresh", ID: "1089189421175287969"}, Custom_ID: "captcha-preharvest-patch:" + taskID})
					return "To refresh this task's remaining uses, click the button at the very bottom right."
				}
				return "\u200b"
			}()})
		}

		embed.Color = 0
		response.Data.Embeds = []disgo.Embed{embed}
	}

	if err := vars.Respond(response, w); err != nil {
		console.Log(err)
	}
}
