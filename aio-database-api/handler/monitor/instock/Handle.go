package instock

import (
	"database-api/user"
	"net/http"
)

func Handle(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		user.VerifyAdmin(post)(w, r)
	}

}
