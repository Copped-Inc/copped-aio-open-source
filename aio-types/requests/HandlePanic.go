package requests

import (
	"errors"
	"github.com/Copped-Inc/aio-types/console"
	"net/http"
)

func HandlePanic(r *http.Request, w http.ResponseWriter) {

	rec := recover()
	if rec != nil {
		var err error
		switch t := rec.(type) {
		case string:
			err = errors.New(t)
		case error:
			err = t
		default:
			err = errors.New("unknown error")
		}
		console.RequestLog(r, "Panic", err.Error())
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
	}

}
