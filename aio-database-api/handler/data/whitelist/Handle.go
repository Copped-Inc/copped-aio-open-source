package whitelist

import (
	"database-api/user"
	"net/http"
)

func Handle(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodDelete:
		user.Get(del)(w, r)
	case http.MethodPost:
		user.Get(post)(w, r)
	}

}
