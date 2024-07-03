package webhook

import (
	"database-api/user"
	"net/http"
)

func Handle(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		user.Get(post)(w, r)
	case http.MethodDelete:
		user.Get(delete)(w, r)
	}

}
