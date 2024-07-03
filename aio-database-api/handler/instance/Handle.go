package instance

import (
	"database-api/user"
	"net/http"
)

func Handle(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		user.Get(get)(w, r)
	case http.MethodPost:
		post(w, r)
	}

}
