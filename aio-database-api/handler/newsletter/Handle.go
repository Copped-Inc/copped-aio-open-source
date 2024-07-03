package newsletter

import (
	"net/http"
)

func Handle(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		post(w, r)
	}

}
