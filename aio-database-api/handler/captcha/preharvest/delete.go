package preharvest

import (
	"database-api/database"
	"database-api/preharvest"
	"net/http"

	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/responses"

	"github.com/gorilla/mux"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func delete(w http.ResponseWriter, r *http.Request) {
	if _, err := database.GetDatabase().Doc("preharvest/" + mux.Vars(r)["task-id"]).Delete(database.GetContext()); err != nil {
		if status.Code(err) == codes.NotFound {
			console.ErrorRequest(w, r, err, http.StatusNotFound)
		} else {
			console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		}
		return
	}

	responses.SendJson(nil, http.StatusNoContent, w, r)
	preharvest.Remove <- mux.Vars(r)["task-id"]
}
