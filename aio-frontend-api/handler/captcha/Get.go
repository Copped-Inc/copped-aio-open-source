package captcha

import (
	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/cookies"
	"github.com/Copped-Inc/aio-types/helper"
	"net/http"
)

func Get(w http.ResponseWriter, r *http.Request) {

	captcha := Template{}

	if helper.System != "linux" {
		cookies.Add(w, "localhost", "true")
	}

	templ, err := helper.ReadTemplate("html/captcha.html")
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	err = templ.Execute(w, captcha)
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
	}

}

type Template struct{}
