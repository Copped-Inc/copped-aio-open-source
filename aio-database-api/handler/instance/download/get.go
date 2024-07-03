package download

import (
	"database-api/user"
	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/statistic"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

func get(w http.ResponseWriter, r *http.Request, _ *user.Database) {

	file := mux.Vars(r)["file"]
	exe, err := os.ReadFile(file + ".exe")
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	statistic.Status(r, http.StatusOK)
	_, _ = w.Write(exe)

}
