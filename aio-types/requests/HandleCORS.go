package requests

import (
	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/helper"
	"github.com/Copped-Inc/aio-types/statistic"
	"net/http"
)

func HandleCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		statistic.New(r)
		console.Request(r)

		w.Header().Set("Access-Control-Allow-Origin", helper.Active)
		w.Header().Set("Access-Control-Allow-Headers", "password, confirm, code, sitekey, browser")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions && r.Header.Get("sec-Fetch-Mode") != "no-cors" {
			statistic.Status(r, http.StatusOK)
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Header.Get("KithEU-Current-Sitekey") != "" {
			helper.KithEUSitekey = r.Header.Get("KithEU-Current-Sitekey")
		}

		defer HandlePanic(r, w)
		h.ServeHTTP(w, r)
	})
}
