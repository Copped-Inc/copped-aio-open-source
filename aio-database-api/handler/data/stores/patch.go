package stores

import (
	"database-api/handler/websocket"
	"database-api/user"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/modules"
	"github.com/Copped-Inc/aio-types/responses"
	"golang.org/x/exp/slices"
)

func patch(w http.ResponseWriter, r *http.Request, database *user.Database) {
	var req request

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		console.ErrorRequest(w, r, err, http.StatusBadRequest)
		return
	} else if !slices.Contains(modules.Sites, req.Store) {
		console.ErrorRequest(w, r, errors.New("invalid store"), http.StatusBadRequest)
		return
	}

	database.UpdateStore(req.Store, req.Value)
	if err := database.Update(); err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	websocket.Websocket{
		Action: websocket.UpdateStores,
		Body:   req,
	}.Send(database.User.ID)

	responses.SendOk(w, r)
}

type request struct {
	ID    string       `json:"id"`
	Store modules.Site `json:"store"`
	Value bool         `json:"value"`
}
