package proxies

import (
	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/helper"
	"net/http"
)

func Handle(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		verifyUser(get)(w, r)
	}

}

func verifyUser(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := helper.GetClaim("id", r)
		if err != nil && !helper.IsMaster(r.Header.Get("Password")) {
			console.ErrorRequest(w, r, err, http.StatusUnauthorized)
		}

		f(w, r)
	}
}
