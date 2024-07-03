package notifications

import (
	"net/http"
)

func Handle(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		post(w, r)
	case http.MethodGet:
		get(w, r)
	case http.MethodPatch:
		patch(w, r)
	case http.MethodDelete:
		delete(w, r)
	case http.MethodPut:
		put(w, r)
	}

}

type Request struct {
	Title string `json:"title,omitempty"`
	Text  string `json:"text,omitempty"`
	Read  *bool  `json:"read,omitempty"`
}
