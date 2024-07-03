package redirect

import (
	"net/http"
)

func get(w http.ResponseWriter, r *http.Request) {

	link := r.URL.Query().Get("link")
	http.Redirect(w, r, link, http.StatusFound)

}
