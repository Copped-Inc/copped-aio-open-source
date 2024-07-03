package shipping

import (
	"database-api/handler/websocket"
	"database-api/user"
	"encoding/json"
	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/responses"
	"net/http"
)

func patch(w http.ResponseWriter, r *http.Request, data *user.Data, database *user.Database) {

	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusBadRequest)
		return
	}

	err = data.UpdateShipping(req.Shipping).UpdateData(r, database)
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	websocket.Websocket{
		Action: websocket.UpdateShipping,
		Body:   req,
	}.Send(database.User.ID)

	responses.SendOk(w, r)

}

type request struct {
	ID       string          `json:"id"`
	Shipping []user.Shipping `json:"shipping"`
}
