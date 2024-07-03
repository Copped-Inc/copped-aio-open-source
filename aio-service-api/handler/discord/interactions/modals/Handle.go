package modals

import (
	"errors"
	"github.com/Copped-Inc/aio-types/console"
	"github.com/infinitare/disgo"
	"net/http"
	"service-api/handler/discord/interactions/modals/captcha"
	"service-api/handler/discord/interactions/modals/instance"
	"service-api/handler/discord/interactions/modals/newsletter"
	"strings"
)

func Handle(interaction *disgo.Interaction, w http.ResponseWriter) {
	if strings.HasPrefix(interaction.Data.Custom_ID, "newsletter-") {
		newsletter.Button(interaction, w)
	} else {
		switch interaction.Data.Custom_ID {
		case "captcha preharvest":
			captcha.PreharvestSchedule(interaction, w)
		case "instance update":
			instance.Update(interaction, w)
		case "newsletter":
			newsletter.Basic(interaction, w)
		default:
			console.ErrorRequest(w, nil, errors.New("unhandled / unknown modal: "+interaction.Data.Custom_ID), http.StatusBadRequest)
		}
	}
}
