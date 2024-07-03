package user

import (
	"database-api/user"
	"net/http"
)

func Handle(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		user.VerifyAdmin(get)(w, r)
	case http.MethodPatch:
		user.VerifyAdmin(patch)(w, r)
	}

}
