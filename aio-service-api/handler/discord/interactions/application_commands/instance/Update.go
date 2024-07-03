package instance

import (
	"embed"
	"encoding/json"
	"github.com/Copped-Inc/aio-types/console"
	"github.com/infinitare/disgo"
	"net/http"
	"service-api/handler/discord/interactions/vars"
)

//go:embed callback.json
var callback embed.FS

func Update(interaction *disgo.Interaction, w http.ResponseWriter) {

	response := disgo.InteractionResponse{Type: disgo.MODAL}

	data, err := callback.ReadFile("callback.json")
	if err != nil {
		console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(data, &response.Data)
	if err != nil {
		console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
		return
	}

	if err := vars.Respond(response, w); err != nil {
		console.ErrorRequest(w, nil, err, http.StatusInternalServerError)
	}
}
