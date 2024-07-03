package responses

import (
	"encoding/json"
	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/statistic"
	"net/http"
)

func SendJson(data interface{}, statusCode int, w http.ResponseWriter, r *http.Request) {

	if statusCode == http.StatusNoContent || statusCode == http.StatusNotModified {
		w.WriteHeader(statusCode)
		statistic.Status(r, statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	statistic.Response(r, data)
	statistic.Status(r, statusCode)

}
