package product

import (
	"database-api/global"
	"database-api/product"
	"encoding/json"
	"errors"
	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/helper"
	"github.com/Copped-Inc/aio-types/responses"
	"github.com/gorilla/mux"
	"golang.org/x/exp/slices"
	"net/http"
)

func post(w http.ResponseWriter, r *http.Request) {

	if !helper.IsMaster(r.Header.Get("Password")) {
		console.ErrorRequest(w, r, errors.New("invalid authorization password"), http.StatusUnauthorized)
		return
	}

	handle := mux.Vars(r)["handle"]

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		console.ErrorRequest(w, r, err, http.StatusBadRequest)
		return
	}

	p, err := product.Get(handle)
	if err != nil {
		p = product.New(handle, req.State)
		go global.AddProduct()
	}

	if req.State != product.None {
		p.UpdateState(req.State)
	}

	if req.Name != "" {
		p.UpdateName(req.Name)
	}

	if req.Image != "" {
		p.UpdateImage(req.Image)
	}

	if req.Handle != "" && !slices.Contains(p.Handles, req.Handle) {
		p.AddHandle(req.Handle)
	}

	if req.StockX != "" {
		p.UpdateStockX(req.StockX)
	}

	err = p.Save()
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	responses.SendJson(p, http.StatusOK, w, r)

}

type request struct {
	State  product.State `json:"state"`
	Image  string        `json:"image,omitempty"`
	Name   string        `json:"name,omitempty"`
	Handle string        `json:"handle,omitempty"`
	StockX string        `json:"stockx,omitempty"`
}
