package autocomplete

import (
	"errors"
	"github.com/Copped-Inc/aio-types/console"
	"github.com/infinitare/disgo"
	"net/http"
	"service-api/handler/discord/interactions/autocomplete/captcha/preharvest"
)

func Handle(interaction *disgo.Interaction, w http.ResponseWriter) {
	if interaction.Guild_ID == disgo.Snowflake("") { // Insert Server ID here
		switch interaction.Data.Name {
		case "captcha":
			preharvest.Respond(interaction, w)

		default:
			console.ErrorRequest(w, nil, errors.New("unhandled / unknown application command name: "+interaction.Data.Name), http.StatusBadRequest)
		}
	} else {
		console.ErrorRequest(w, nil, errors.New("request guild origin ("+string(interaction.Guild_ID)+") isn't Copped AIO discord"), http.StatusBadRequest)
	}

}
