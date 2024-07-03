package restart

import (
	"errors"
	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/responses"
	"github.com/Copped-Inc/aio-types/secrets"
	"github.com/gorilla/mux"
	"net/http"
)

func get(w http.ResponseWriter, r *http.Request) {

	site := mux.Vars(r)["site"]
	if site == "" {
		console.ErrorRequest(w, r, errors.New("missing site"), http.StatusBadRequest)
		return
	}

	req, err := http.NewRequest(http.MethodPost, "https://monitor.copped-inc.com/restart/"+site, nil)
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	req.Header.Add("Password", secrets.API_Admin_PW)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	if res.StatusCode != http.StatusOK {
		console.ErrorRequest(w, r, errors.New("request to "+res.Request.URL.String()+" failed with response "+res.Status), res.StatusCode)
		return
	}

	responses.SendOk(w, r)

}
