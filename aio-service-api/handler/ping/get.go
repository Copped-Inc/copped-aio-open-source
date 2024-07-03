package ping

import (
	"github.com/Copped-Inc/aio-types/responses"
	"net/http"
)

func get(w http.ResponseWriter, r *http.Request) {

	responses.SendJson(struct {
		Status string `json:"status"`
	}{
		Status: http.StatusText(http.StatusOK),
	}, http.StatusOK, w, r)

}
