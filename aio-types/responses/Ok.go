package responses

import (
	"github.com/Copped-Inc/aio-types/statistic"
	"net/http"
)

func SendOk(w http.ResponseWriter, r *http.Request) {

	http.Error(w, http.StatusText(http.StatusOK), http.StatusOK)
	statistic.Response(r, http.StatusText(http.StatusOK))
	statistic.Status(r, http.StatusOK)

}
