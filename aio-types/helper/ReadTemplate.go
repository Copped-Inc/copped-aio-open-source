package helper

import (
	"errors"
	"html/template"
	"net/http"
)

func ReadTemplate(path ...string) (*template.Template, error) {

	templ, err := template.ParseFiles(path...)
	if err == nil && templ == nil {
		err = errors.New(http.StatusText(http.StatusInternalServerError))
	}

	return templ, err

}
