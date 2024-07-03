package handler

import (
	"github.com/gorilla/mux"
	"monitor-api/handler/ping"
	"monitor-api/handler/restart"
	"net/http"
)

func Add() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/ping", ping.Handle).Methods(http.MethodGet)
	router.HandleFunc("/restart/{site}", restart.Handle).Methods(http.MethodPost)

	return router

}
