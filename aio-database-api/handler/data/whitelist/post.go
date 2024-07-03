package whitelist

import (
	"database-api/handler/websocket"
	"database-api/user"
	"encoding/json"
	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/responses"
	"net/http"
)

func post(w http.ResponseWriter, r *http.Request, database *user.Database) {

	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusBadRequest)
		return
	}

	database, err = database.AddWhitelist(req.Product)
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	err = database.Update()
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	websocket.Websocket{
		Action: websocket.AddWhitelist,
		Body:   req,
	}.Send(database.User.ID)

	responses.SendOk(w, r)

}

type request struct {
	ID      string `json:"id"`
	Product string `json:"product"`
}
