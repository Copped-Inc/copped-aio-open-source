package insta

import (
	"database-api/handler/websocket"
	"database-api/user"
	"encoding/json"
	"net/http"

	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/responses"
)

func patch(w http.ResponseWriter, r *http.Request, database *user.Database) {

	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusBadRequest)
		return
	}

	err = database.UpdateInstance(req.Instance).Update()
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	websocket.Websocket{
		Action: websocket.UpdateInstances,
		Body:   database.Instances,
	}.Send(database.User.ID)

	websocket.Websocket{
		Action: websocket.UpdateSession,
		Body: struct {
			Session user.Session `json:"session"`
		}{
			Session: *database.Session,
		},
	}.Send(database.User.ID)

	responses.SendOk(w, r)

}

type request struct {
	ID       string        `json:"id"`
	Instance user.Instance `json:"instance"`
}
