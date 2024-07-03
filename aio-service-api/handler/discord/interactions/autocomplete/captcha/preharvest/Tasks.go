package preharvest

import (
	"encoding/json"
	"errors"
	"net/http"
	internal "service-api/handler/discord/interactions/application_commands/captcha/preharvest"
	"service-api/handler/discord/interactions/vars"

	"github.com/Copped-Inc/aio-types/captcha/preharvest"
	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/helper"
	"github.com/Copped-Inc/aio-types/modules"
	"github.com/Copped-Inc/aio-types/secrets"
	trie "github.com/Vivino/go-autocomplete-trie"
	"github.com/infinitare/disgo"
)

var (
	stores = func() *trie.Trie {
		t := trie.New()
		t.Insert(func() (sites []string) {
			for _, site := range modules.Sites {
				if site.GetData().Runable && site.GetData().CaptchaRequired {
					sites = append(sites, site.GetData().Name, string(site), site.GetData().URL)
				}
			}
			return sites
		}()...)
		return t
	}()
)

func Respond(interaction *disgo.Interaction, w http.ResponseWriter) {

	var (
		response = disgo.InteractionResponse{Type: disgo.APPLICATION_COMMAND_AUTOCOMPLETE_RESULT, Data: disgo.InteractionCallbackData{}}
		inputs   = interaction.Data.Options[0].Options[0].Options
		user_id  string
	)

	// interaction field that holds the string to complete
	input := inputs[0].Value.(string)

	// creating a new preharvest tasks suggests sites, all other subcommands suggest preharvest task ids
	if interaction.Data.Options[0].Options[0].Options[0].Name == "new" {

		// check whether there's input to match against
		if len(input) != 0 {

			// stores holds the site const, as well as the monitor url domain (including prefix) and the displayed name for each preharvestable site
			// any matches to the user's input, up to 25 (maximum amount of autocomplete suggestions) are returned
			for _, match := range stores.Search(input, 25) {
				// matches then are reassigned their site const
				for _, site := range modules.Sites {
					if match == string(site) || match == site.GetData().Name || match == site.GetData().URL {
						match := false
						// to avoid having duplicate suggestions with the same site const value (e.g. Kith EU and eu.kith.com)
						// the current choices are ranged to check, whether the site is already part of the autocomplete suggestions
						for _, choice := range response.Data.Choices {
							if match = choice.Value.(modules.Site) == site; match {
								break
							}
						}

						if !match {
							response.Data.Choices = append(response.Data.Choices, disgo.ApplicationCommandOptionChoice{Name: site.GetData().Name, Value: site})
						}

						break
					}
				}
			}

		} else
		// otherwise return up to 25 random stores, same as task ID recommendation below

		{
			for i, site := range modules.Sites {
				if i >= 25 {
					break
				}

				response.Data.Choices = append(response.Data.Choices, disgo.ApplicationCommandOptionChoice{Name: site.GetData().Name, Value: site})
			}
		}

	} else {

		if len(inputs) > 1 {
			if inputs[0].Name == "id" {
				user_id = inputs[1].Value.(string)
			} else {
				user_id = inputs[0].Value.(string)
				input = inputs[1].Value.(string)
			}
		}

		// to suggest preharvest task ids, the tasks are fetched from the api
		req, err := http.NewRequest(http.MethodGet, helper.ActiveData+"/captcha/preharvest"+func() (path string) {
			// whether preharvest task id suggestion are based out of all tasks or only those of the user who triggered the command, is based on their permissions
			if !interaction.Member.Permissions.Has(disgo.ADMINISTRATOR) {
				return "/user/" + string(interaction.Member.User.ID)
			} else if user_id != "" {
				return "/user/" + user_id
			}
			return
		}(), nil)
		if err != nil {
			console.Log(err)
			return
		}

		req.Header.Set("Password", secrets.API_Admin_PW)

		res, err := (&http.Client{}).Do(req)
		if err != nil {
			console.Log(err)
			return
		} else if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNoContent {
			console.Log(errors.New("request to " + res.Request.URL.String() + " failed with response " + res.Status))
			return
		} else if res.StatusCode == http.StatusNoContent {
			response.Data.Choices = []disgo.ApplicationCommandOptionChoice{{Name: func() string {
				if user_id != "" {
					return "The specified user doesn't have any captcha preharvest tasks."
				}
				return "There are no captcha preharvest tasks. Create one for it to show up here."
			}(), Value: internal.NoContent}}
			vars.Respond(response, w)
			return
		}

		var tasks []preharvest.Task

		if err = json.NewDecoder(res.Body).Decode(&tasks); err != nil {
			console.Log(err)
			return
		}

		// check if there's actual input to match task IDs against
		if len(input) != 0 {

			// make a new trie populated with the preharvest tasks returned from fetching the database api
			suggestions := trie.New()

			for _, task := range tasks {
				suggestions.Insert(task.ID)
			}

			// search for matching suggestions for the input
			for _, suggestion := range suggestions.Search(input, 25) {

				for _, task := range tasks {

					// reassign task to suggested IDs
					if suggestion == task.ID {
						match := false
						// to avoid having duplicate suggestions, the current choices are ranged to check whether the id is already getting suggested
						for _, choice := range response.Data.Choices {
							if match = choice.Value.(string) == suggestion; match {
								break
							}
						}

						if !match {
							response.Data.Choices = append(response.Data.Choices, disgo.ApplicationCommandOptionChoice{Name: task.Site.GetData().Name + " — " + suggestion, Value: suggestion})
						}

						break
					}
				}
			}

		} else
		// else return some tasks to pick from randomly (could add recommendations based on popularity here at some point)

		{
			for i, task := range tasks {
				if i > 24 {
					break
				}

				response.Data.Choices = append(response.Data.Choices, disgo.ApplicationCommandOptionChoice{Name: task.Site.GetData().Name + " — " + task.ID, Value: task.ID})
			}
		}
	}

	vars.Respond(response, w)
}
