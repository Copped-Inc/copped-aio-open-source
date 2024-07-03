package update

import (
	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/responses"
	"github.com/Copped-Inc/aio-types/statistic"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {

	f := mux.Vars(r)["file"]
	version, err := os.ReadFile("version-" + strings.Split(f, ".")[0])
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusBadRequest)
		return
	}

	if r.Header.Get("version") == string(version) {
		responses.SendOk(w, r)
		return
	}

	exe, err := os.ReadFile(f)
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusBadRequest)
		return
	}

	w.Header().Set("Version", string(version))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusCreated)
	statistic.Status(r, http.StatusCreated)
	_, _ = w.Write(exe)

}
