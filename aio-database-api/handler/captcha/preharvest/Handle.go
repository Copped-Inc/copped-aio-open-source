package preharvest

import (
	"errors"
	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/helper"
	"net/http"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	if !helper.IsMaster(r.Header.Get("Password")) {
		console.ErrorRequest(w, r, errors.New("invalid authorization password"), http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case http.MethodGet:
		get(w, r)
	case http.MethodPost:
		post(w, r)
	case http.MethodDelete:
		delete(w, r)
	case http.MethodPatch:
		patch(w, r)
	}
}
