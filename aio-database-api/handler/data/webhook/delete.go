package webhook

import (
	"database-api/handler/websocket"
	"database-api/user"
	"encoding/json"
	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/responses"
	"net/http"
)

func delete(w http.ResponseWriter, r *http.Request, database *user.Database) {

	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusBadRequest)
		return
	}

	err = database.DeleteWebhook(req.Webhook).Update()
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	websocket.Websocket{
		Action: websocket.DeleteWebhook,
		Body:   req,
	}.Send(database.User.ID)

	responses.SendOk(w, r)

}
