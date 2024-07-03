package preharvest

import (
	_ "embed"
	"encoding/json"
	"github.com/Copped-Inc/aio-types/console"
	"github.com/infinitare/disgo"
	"net/http"
	"service-api/handler/discord/interactions/vars"
)

//go:embed callback.json
var callback []byte

func Schedule(interaction *disgo.Interaction, taskID string, w http.ResponseWriter) {
	response := disgo.InteractionResponse{Type: disgo.MODAL}

	if err := json.Unmarshal(callback, &response.Data); err != nil {
		console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
		return
	}

	(*response.Data.Components)[1].Components[0].Custom_ID = taskID

	if err := vars.Respond(response, w); err != nil {
		console.ErrorLog(err)
	}
}
