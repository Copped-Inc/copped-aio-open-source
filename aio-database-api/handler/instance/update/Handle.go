package update

import (
	"database-api/user"
	"net/http"
)

func Handle(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		user.Verify(get)(w, r)
	case http.MethodPost:
		user.VerifyAdmin(post)(w, r)
	}

}
