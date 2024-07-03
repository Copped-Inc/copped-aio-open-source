package product

import (
	"database-api/product"
	"errors"
	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/helper"
	"github.com/Copped-Inc/aio-types/responses"
	"github.com/gorilla/mux"
	"net/http"
)

func get(w http.ResponseWriter, r *http.Request) {

	if !helper.IsMaster(r.Header.Get("Password")) {
		console.ErrorRequest(w, r, errors.New("invalid authorization password"), http.StatusUnauthorized)
		return
	}

	handle := mux.Vars(r)["handle"]
	p, err := product.Get(handle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responses.SendJson(p, http.StatusOK, w, r)

}
