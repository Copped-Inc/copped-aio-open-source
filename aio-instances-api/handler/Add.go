package handler

import (
	"github.com/Copped-Inc/aio-types/requests"
	"github.com/gorilla/mux"
	"instances-api/handler/ping"
	"instances-api/handler/update"
	"net/http"
)

func Add() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(requests.HandleCORS)

	router.HandleFunc("/instance/{id}", Handle).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/update", update.Handle).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/ping", ping.Handle).Methods(http.MethodGet)

	return router

}
