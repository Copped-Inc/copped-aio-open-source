package purchase

import (
	"database-api/user"
	"net/http"
)

func Handle(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		getHandle(w, r)
	case http.MethodPost:
		user.VerifyAdmin(post)(w, r)
	}

}
