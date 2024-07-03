package linked_roles

import (
	"net/http"
	"net/url"
	"time"

	"github.com/Copped-Inc/aio-types/discord"
	"github.com/google/uuid"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	state := uuid.New().String()

	http.SetCookie(w, &http.Cookie{
		Name:     "client_state",
		Value:    state,
		Secure:   true,
		HttpOnly: true,
		Path:     "/discord/linked-roles",
		Domain:   ".copped-inc.com",
		Expires:  time.Now().Add(5 * time.Minute),
	})

	http.Redirect(w, r, "https://discord.com/api/oauth2/authorize?client_id="+discord.Application_ID+"&redirect_uri="+url.QueryEscape("https://database.copped-inc.com/discord/linked-roles")+"&response_type=code&scope=identify%20email%20role_connections.write&state="+state, http.StatusFound)
}
