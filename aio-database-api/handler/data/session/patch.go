package session

import (
	"database-api/handler/websocket"
	"database-api/user"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/responses"
)

func patch(w http.ResponseWriter, r *http.Request, database *user.Database) {

	var req request

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		console.ErrorRequest(w, r, err, http.StatusBadRequest)
		return
	}

	if database.Session != nil {
		if database.Session.Status == "Stopped by Admin" {
			console.ErrorRequest(w, r, errors.New("your session was disabled by an admin"), http.StatusUnauthorized)
			return
		}
	}

	if err := database.UpdateSession(req.Session).Update(); err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	websocket.Websocket{
		Action: websocket.UpdateSession,
		Body:   req,
	}.Send(database.User.ID)

	responses.SendOk(w, r)

}

type request struct {
	ID      string       `json:"id"`
	Session user.Session `json:"session"`
}
