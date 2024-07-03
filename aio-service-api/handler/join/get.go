package join

import (
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {

	location := strings.ToLower(mux.Vars(r)["location"])
	switch location {
	case "copped-inc":
		http.Redirect(w, r, "https://discord.com/servers/copped-inc-811233408474677289", http.StatusFound)
	default:
		http.Redirect(w, r, "https://discord.com/servers/copped-inc-811233408474677289", http.StatusFound)
	}

}
