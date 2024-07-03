package console

import (
	"errors"
	"github.com/Copped-Inc/aio-types/statistic"
	"net/http"
)

func ErrorRequest(w http.ResponseWriter, r *http.Request, err error, code int) {

	if err == nil {
		err = errors.New("")
	}

	text := http.StatusText(code)
	if text == "" {
		switch code {
		case 601:
			text = "OOS"
		}
	}

	if err == nil {
		err = errors.New("")
	}

	statistic.AddError(r, err)

	RequestLog(r, "Error", text, err.Error())
	statistic.Status(r, code)
	http.Error(w, text, code)

}

func ErrorLog(err error) {

	Log("Error", err.Error())

}
