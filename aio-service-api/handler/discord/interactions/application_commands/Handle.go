package application_commands

import (
	"errors"
	"github.com/Copped-Inc/aio-types/console"
	"github.com/infinitare/disgo"
	"net/http"
	"service-api/handler/discord/interactions/application_commands/captcha/preharvest"
	"service-api/handler/discord/interactions/application_commands/instance"
	"service-api/handler/discord/interactions/application_commands/monitor"
	"service-api/handler/discord/interactions/application_commands/newsletter"
	"service-api/handler/discord/interactions/application_commands/purchase"
	"strconv"
)

func Handle(interaction *disgo.Interaction, w http.ResponseWriter) {
	if interaction.Guild_ID == disgo.Snowflake("") { // Insert Server ID here
		switch interaction.Data.Type {
		case disgo.ApplicationCommandType_CHAT_INPUT:
			switch interaction.Data.Name {
			case "monitor":
				monitor.Send(interaction, w)
			case "instance":
				switch interaction.Data.Options[0].Name {
				case "update":
					instance.Update(interaction, w)
				}
			case "purchase":
				purchase.Link(interaction, w)
			case "newsletter":
				newsletter.Send(interaction, w)
			case "captcha":
				preharvest.Handle(interaction, w)

			default:
				console.ErrorRequest(w, nil, errors.New("unhandled / unknown application command name: "+interaction.Data.Name), http.StatusBadRequest)
			}
		case disgo.ApplicationCommandType_MESSAGE:
			switch interaction.Data.Name {
			case "monitor blacklist":
				monitor.Send(interaction, w)
			default:
				console.ErrorRequest(w, nil, errors.New("unhandled / unknown application command name: "+interaction.Data.Name), http.StatusBadRequest)
			}
		default:
			console.ErrorRequest(w, nil, errors.New("unhandled / unknown application command type: "+strconv.Itoa(int(interaction.Data.Type))), http.StatusBadRequest)
		}
	} else {
		console.ErrorRequest(w, nil, errors.New("request guild origin ("+string(interaction.Guild_ID)+") isn't Copped AIO discord"), http.StatusBadRequest)
	}

}
