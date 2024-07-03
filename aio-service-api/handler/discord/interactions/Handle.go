package interactions

import (
	"errors"
	"github.com/Copped-Inc/aio-types/console"
	"net/http"
	"service-api/handler/discord/interactions/application_commands"
	"service-api/handler/discord/interactions/autocomplete"
	"service-api/handler/discord/interactions/message_components"
	"service-api/handler/discord/interactions/modals"

	"github.com/Copped-Inc/aio-types/discord"
	"github.com/infinitare/disgo"
)

var client = disgo.NewInteractionClient(discord.PublicKey)

func Handle(w http.ResponseWriter, r *http.Request) {
	interaction, err := client.Verify(w, r)
	if interaction == nil {
		if err != nil {
			console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		}
		return
	}

	switch interaction.Type {
	case disgo.APPLICATION_COMMAND:
		application_commands.Handle(interaction, w)
	case disgo.MESSAGE_COMPONENT:
		message_components.Handle(interaction, w)
	case disgo.MODAL_SUBMIT:
		modals.Handle(interaction, w)
	case disgo.APPLICATION_COMMAND_AUTOCOMPLETE:
		autocomplete.Handle(interaction, w)
	default:
		console.ErrorRequest(w, r, errors.New("unhandled / unknown reaction type"), http.StatusBadRequest)
	}
}
