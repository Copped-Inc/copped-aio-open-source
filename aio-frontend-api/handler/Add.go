package handler

import (
	"frontend-api/handler/captcha"
	"frontend-api/handler/files"
	"frontend-api/handler/ping"
	"frontend-api/handler/utility"
	"github.com/Copped-Inc/aio-types/requests"
	"github.com/gorilla/mux"
	"net/http"
)

func Add() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(requests.HandleCORS)

	router.HandleFunc("/", Get).Methods(http.MethodGet)
	router.HandleFunc("/captcha", captcha.Get).Methods(http.MethodGet)
	router.HandleFunc("/ping", ping.Handle).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/utility/{type}", utility.Get).Methods(http.MethodGet, http.MethodOptions)

	router.HandleFunc("/favicon.ico", files.Get).Methods(http.MethodGet)
	router.HandleFunc("/{type}/{file}", files.Get).Methods(http.MethodGet)
	router.HandleFunc("/{path}/{type}/{file}", files.Get).Methods(http.MethodGet)

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/utility/404", http.StatusFound)
	})
	return router

}
