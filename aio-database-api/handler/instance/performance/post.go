package performance

import (
	"github.com/Copped-Inc/aio-types/responses"
	"net/http"
)

func post(w http.ResponseWriter, r *http.Request) {

	responses.SendJson(response{Performance: "success"}, http.StatusOK, w, r)

}

type response struct {
	Performance string `json:"performance"`
}
