package session

import (
	"database-api/user"
	"net/http"
)

func Handle(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPatch:
		user.Get(patch)(w, r)
	}

}
