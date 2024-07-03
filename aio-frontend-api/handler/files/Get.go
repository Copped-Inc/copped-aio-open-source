package files

import (
	"net/http"
	"strings"
)

func Get(w http.ResponseWriter, r *http.Request) {

	path := strings.Split(r.URL.Path, "/")
	file := path[len(path)-1]
	if file == "favicon.ico" {
		http.ServeFile(w, r, "html/img/"+file)
		return
	}

	fileType := path[len(path)-2]

	if fileType == "json" && strings.Contains(r.Header.Get("accept-encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")
		http.ServeFile(w, r, "html/json/"+file+".gz")
	} else {
		http.ServeFile(w, r, "html/"+fileType+"/"+file)
	}

}
