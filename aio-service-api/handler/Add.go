package handler

import (
	"github.com/Copped-Inc/aio-types/requests"
	"github.com/gorilla/mux"
	"net/http"
	"service-api/handler/discord/interactions"
	"service-api/handler/discord/linked_roles"
	"service-api/handler/join"
	"service-api/handler/ping"
	"service-api/handler/proxies"
	"service-api/handler/redirect"
)

func Add() *mux.Router {

	router := mux.NewRouter()
	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(requests.HandleCORS)

	router.HandleFunc("/ping", ping.Handle).Methods(http.MethodGet)
	router.HandleFunc("/discord/interactions", interactions.Handle).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/discord/linked-roles", linked_roles.Handle).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/proxies", proxies.Handle).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/redirect", redirect.Handle).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/join/{location}", join.Handle).Methods(http.MethodGet, http.MethodOptions)

	http.Handle("/", router)
	return router

}

/*func logHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		x, err := httputil.DumpRequest(r, true)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
		log.Printf("%q", x)
		rec := httptest.NewRecorder()
		fn(rec, r)
		y, err := httputil.DumpResponse(rec.Result(), true)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
		log.Printf("%q", y)

		// this copies the recorded response to the response writer
		for k, v := range rec.Result().Header {
			w.Header()[k] = v
		}
		w.WriteHeader(rec.Code)
		rec.Body.WriteTo(w)
	}
}*/
