package message_components

import (
	"errors"
	"github.com/Copped-Inc/aio-types/console"
	"net/http"
	"strconv"
	"strings"

	"service-api/handler/discord/interactions/message_components/captcha/preharvest"
	"service-api/handler/discord/interactions/message_components/newsletter"

	"github.com/infinitare/disgo"
)

func Handle(interaction *disgo.Interaction, w http.ResponseWriter) {
	switch interaction.Data.Component_Type {
	case disgo.Button:
		switch strings.Split(interaction.Data.Custom_ID, "-")[0] {
		case "newsletter":
			newsletter.Buttons(interaction, w)
		case "captcha":
			preharvest.Buttons(interaction, w)

		default:
			console.ErrorRequest(w, nil, errors.New("unhandled / unknown button component name: "+interaction.Data.Custom_ID), http.StatusBadRequest)
		}
	case disgo.String_Select:
		switch strings.Split(interaction.Data.Custom_ID, "-")[0] {
		case "captcha":
			preharvest.SelectMenu(interaction, w)

		default:
			console.ErrorRequest(w, nil, errors.New("unhandled / unknown string select component name: "+interaction.Data.Custom_ID), http.StatusBadRequest)
		}
	case disgo.User_Select:
		switch strings.Split(interaction.Data.Custom_ID, "-")[0] {
		case "captcha":
			preharvest.SelectMenu(interaction, w)

		default:
			console.ErrorRequest(w, nil, errors.New("unhandled / unknown string select component name: "+interaction.Data.Custom_ID), http.StatusBadRequest)
		}

	default:
		console.ErrorRequest(w, nil, errors.New("unhandled / unknown component type: "+strconv.Itoa(int(interaction.Data.Component_Type))), http.StatusBadRequest)
	}
}
