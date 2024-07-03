package payments

import (
	"database-api/handler/websocket"
	"database-api/user"
	"encoding/json"
	"net/http"

	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/modules"
	"github.com/Copped-Inc/aio-types/responses"
)

func post(w http.ResponseWriter, r *http.Request, database *user.Database) {

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		console.ErrorRequest(w, r, err, http.StatusBadRequest)
		return
	}

	websocket.UserPayments(req, database.User.ID)
	responses.SendOk(w, r)

}

type request struct {
	ID    string       `json:"id"`
	Store modules.Site `json:"store"`
	Data  string       `json:"data"`
	State state        `json:"state"`
}

type state int

const (
	Created state = iota
	Accepted
	Declined
	Finalized
)
