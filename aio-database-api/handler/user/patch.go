package user

import (
	"database-api/handler/websocket"
	"database-api/user"
	"encoding/json"
	"errors"
	userTypes "github.com/Copped-Inc/aio-types/user"
	"net/http"

	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/modules"
	"github.com/Copped-Inc/aio-types/responses"
	"github.com/gorilla/mux"
)

func patch(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		console.ErrorRequest(w, r, errors.New("id is empty"), http.StatusBadRequest)
		return
	}

	var req user.Database
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusBadRequest)
		return
	}

	u, err := user.FromId(id)
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	if u.Session != nil && u.Session.Status != req.Session.Status {
		u.Session.Status = req.Session.Status

		websocket.Websocket{
			Action: websocket.UpdateSession,
			Body: struct {
				ID      string       `json:"id"`
				Session user.Session `json:"session"`
			}{
				ID:      "overriden",
				Session: *u.Session,
			},
		}.Send(u.User.ID)
	}

	if u.User.Subscription.Plan != req.User.Subscription.Plan {
		u.User.Subscription.Plan = req.User.Subscription.Plan
		if u.User.Subscription.Plan == 0 {
			u.User.Subscription.State = userTypes.Disabled
		}
		// TODO add handling
	}

	for _, store := range modules.Sites {
		if u.Settings != nil && req.Settings.Stores == u.Settings.Stores {
			continue
		} else if newSetting := req.Settings.Stores.IsEnabled(store); newSetting || u.Settings != nil && newSetting != u.Settings.Stores.IsEnabled(store) {
			console.Log(store, req.Settings, req.Settings.Stores.IsEnabled(store))
			if u.Settings == nil {
				u.Settings = &user.Settings{
					Stores: "",
				}
			}

			if !newSetting {
				u.Settings.Stores.Remove(store)
			} else {
				u.Settings.Stores.Add(store)
			}

			websocket.Websocket{
				Action: websocket.UpdateStores,
				Body: struct {
					ID    string       `json:"id"`
					Store modules.Site `json:"store"`
					Value bool         `json:"value"`
				}{
					ID:    "overriden",
					Store: store,
					Value: newSetting,
				},
			}.Send(u.User.ID)
		}
	}

	console.Log(u.Settings.Stores)
	u.User.InstanceLimit = req.User.InstanceLimit
	u.User.Picture = req.User.Picture
	if u.Update() != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	responses.SendOk(w, r)
}
