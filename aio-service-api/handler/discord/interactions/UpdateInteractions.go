package interactions

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/discord"
	"github.com/Copped-Inc/aio-types/modules"
	"github.com/Copped-Inc/aio-types/subscriptions"
	"github.com/infinitare/disgo"
	"net/http"
)

//go:embed interactions.json
var application_command []byte

func Update() {
	var (
		data         map[string][]disgo.ApplicationCommand
		plans, sites []disgo.ApplicationCommandOptionChoice
	)

	if err := json.Unmarshal(application_command, &data); err != nil {
		console.ErrorLog(err)
		return
	}

	for _, plan := range subscriptions.Plans {
		plans = append(plans, disgo.ApplicationCommandOptionChoice{Name: plan.GetData().Name, Value: plan})
	}

	data["application_commands"][3].Options[0].Options[0].Options[0].Choices = plans

	preharvest_sites := data["application_commands"][5].Options[0].Options[0].Options[0]

	for _, site := range modules.Sites {
		if site.GetData().Runable && site.GetData().CaptchaRequired {
			sites = append(sites, disgo.ApplicationCommandOptionChoice{Name: site.GetData().Name, Value: site})
		}
	}

	if len(sites) <= 25 {
		autocomplete := false
		preharvest_sites.Autocomplete = &autocomplete
		preharvest_sites.Choices = append(preharvest_sites.Choices, sites...)
	}

	data["application_commands"][5].Options[0].Options[0].Options[0] = preharvest_sites

	new, err := json.Marshal(data["application_commands"])
	if err != nil {
		console.ErrorLog(err)
		return
	}

	req, err := http.NewRequest(http.MethodPut, "https://discord.com/api/v"+discord.API_Version+"/applications/"+discord.Application_ID+"/commands", bytes.NewBuffer(new))
	if err != nil {
		console.ErrorLog(err)
		return
	}

	req.Header.Set("Content-type", "application/json")
	req.Header.Add("Authorization", discord.Bearer)

	res, err := (&http.Client{}).Do(req)
	if err != nil {
		console.ErrorLog(err)
	} else if res.StatusCode != http.StatusOK {
		console.ErrorLog(errors.New("request to " + res.Request.URL.String() + " failed with response " + res.Status))
	}

	console.Log("interactions updated successfully")
}
