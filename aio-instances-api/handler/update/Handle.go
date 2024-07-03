package update

import (
	"errors"
	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/helper"
	"net/http"
)

func Handle(w http.ResponseWriter, r *http.Request) {

	if !helper.IsMaster(r.Header.Get("Password")) {
		console.ErrorRequest(w, r, errors.New(http.StatusText(http.StatusUnauthorized)), http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case http.MethodPost:
		post(w, r)
	}

}
