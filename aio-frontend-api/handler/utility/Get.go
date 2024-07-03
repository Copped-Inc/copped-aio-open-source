package utility

import (
	"html/template"
	"net/http"
	"net/url"

	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/helper"
	"github.com/gorilla/mux"
)

func Get(w http.ResponseWriter, r *http.Request) {

	dashboard := func() Template {
		vars := mux.Vars(r)
		params := r.URL.Query()

		return Template{
			ErrorTitle: func() string {
				if p := params.Get("error"); p != "" {
					t, _ := url.QueryUnescape(p)
					return t
				}
				switch vars["type"] {
				case "discord":
					return "Success"
				case "403":
					return "Unauthorised"
				case "500":
					return "Internal Server error"
				default:
					return "Not Found"
				}
			}(),
			ErrorMessage: func() string {
				if p := params.Get("message"); p != "" {
					t, _ := url.QueryUnescape(p)
					return t
				}
				switch vars["type"] {
				case "discord":
					return "You have successfully linked your Discord account to Discord."
				case "403":
					return "You do not have permission to visit this page."
				case "500":
					return "Something went wrong, please try again :C"
				default:
					return "This page does not exist. Return to the Dashboard."
				}
			}(),
			Button: func() bool {
				if p := params.Get("button"); p != "" {
					t, _ := url.QueryUnescape(p)
					return t == "true"
				}
				switch vars["type"] {
				case "403":
					return false
				default:
					return true
				}
			}(),
			ButtonTitle: func() string {
				if p := params.Get("title"); p != "" {
					t, _ := url.QueryUnescape(p)
					return t
				}
				switch vars["type"] {
				case "discord":
					return "Back to Discord"
				case "403":
					return ""
				case "500":
					return "try again"
				default:
					return "Dashboard"
				}
			}(),
			ButtonLocation: func() string {
				if p := params.Get("location"); p != "" {
					t, _ := url.QueryUnescape(p)
					return t
				}
				switch vars["type"] {
				case "discord":
					return "" // Insert Message Location here
				case "500":
					if r.Referer() != "" {
						return r.Referer()
					}
				}
				return "/"
			}(),
		}
	}()

	templ := &template.Template{}

	templ, err := helper.ReadTemplate("html/utility.html")
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
	ErrorTitle     string `query:"error"`
	ErrorMessage   string `query:"message"`
	Button         bool
	ButtonTitle    string `query:"button"`
	ButtonLocation string `query:"location"`
}
