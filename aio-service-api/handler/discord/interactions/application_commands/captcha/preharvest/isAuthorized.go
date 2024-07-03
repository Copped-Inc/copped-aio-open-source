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
)

func isAuthorized(interaction *disgo.Interaction, taskID string, w http.ResponseWriter, response disgo.InteractionResponse, action string) (is bool) {
	if interaction.Member.Permissions.Has(disgo.ADMINISTRATOR) {
		return true
	}

	req, err := http.NewRequest(http.MethodGet, helper.ActiveData+"/captcha/preharvest/"+taskID, nil)
	if err != nil {
		console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
		return
	}

	req.Header.Set("Password", secrets.API_Admin_PW)

	res, err := (&http.Client{}).Do(req)
	if err != nil {
		console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
		return
	} else if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNoContent && res.StatusCode != http.StatusNotModified {
		console.ErrorRequest(w, nil, errors.New("request to "+res.Request.URL.String()+" failed with response "+res.Status), http.StatusInternalServerError)
		return
	} else

	// prompt the user to retry action or go back to list of preharvest tasks if the task specified couldn't be found
	if res.StatusCode == http.StatusNotFound {
		NotFound(interaction, response, taskID, w)
		return
	}

	var task preharvest.Task

	if err := json.NewDecoder(res.Body).Decode(&task); err != nil {
		console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
		return
	}

	// only administrators can perform action on the preharvest tasks of others
	if task.User_ID != interaction.Member.User.ID {

		response.Data.Embeds[0].Description = "You aren't allowed to " + action + " the preharvest task **` " + taskID + " `**! You can only " + action + " your own preharvest tasks."

		InvalidPermissions(interaction, response, w)
		return
	}

	return true
}
