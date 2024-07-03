package handler

import (
	"errors"
	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/cookies"
	"github.com/Copped-Inc/aio-types/responses"
	"net/http"
)

func get(w http.ResponseWriter, r *http.Request) {

	redirect := cookies.Get(r, "redirect")
	if redirect == nil {
		console.ErrorRequest(w, r, errors.New("redirect cookie is empty"), http.StatusBadRequest)
		return
	}

	keys, ok := r.URL.Query()["code"]

	if !ok || len(keys[0]) < 1 {
		console.ErrorRequest(w, r, errors.New("code is empty"), http.StatusBadRequest)
		return
	}

	responses.Redirect(w, r, redirect.Value+"?code="+keys[0])

}
