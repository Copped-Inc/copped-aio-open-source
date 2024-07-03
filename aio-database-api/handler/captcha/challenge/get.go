package challenge

import (
	"github.com/Copped-Inc/aio-types/helper"
	"github.com/Copped-Inc/aio-types/responses"
	"net/http"
)

func get(w http.ResponseWriter, r *http.Request) {

	responses.SendJson(response{Sitekey: helper.KithEUSitekey}, http.StatusOK, w, r)

}

type response struct {
	Sitekey string `json:"sitekey"`
}
