package ping

import (
	"github.com/Copped-Inc/aio-types/responses"
	"net/http"
)

func get(w http.ResponseWriter, r *http.Request) {

	responses.SendOk(w, r)

}
