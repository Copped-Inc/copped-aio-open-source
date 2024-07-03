package console

import (
	"net/http"
)

func Request(r *http.Request) {

	if r.URL.Path == "/ping" {
		return
	}
	RequestLog(r, r.Method, r.URL.Path)

}
