package handler

import (
	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/cookies"
	"github.com/Copped-Inc/aio-types/helper"
	"html/template"
	"net/http"
)

func Get(w http.ResponseWriter, r *http.Request) {

	dashboard := Template{
		IsDev: false,
	}
	templ := &template.Template{}

	if helper.System == "windows" {
		cookies.Add(w, "localhost", "true")
	}

	id, err := helper.GetClaim("id", r)
	if err != nil {
		http.Redirect(w, r, helper.ActiveData+"/login", http.StatusFound)
		return
	}

	if id == "573810185840361482" /* Valle */ || id == "432153783112433685" /* ileFix */ {
		dashboard.IsDev = true
	}

	templ, err = helper.ReadTemplate("html/dashboard.html")
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	err = templ.Execute(w, dashboard)
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
	}

}

type Template struct {
	IsDev bool
}
