package login

import (
	"net/http"
)

func Handle(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		get(w, r)
	}

}
