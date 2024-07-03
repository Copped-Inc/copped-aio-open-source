package handler

import (
	"database-api/handler/captcha"
	"database-api/handler/captcha/challenge"
	"database-api/handler/captcha/preharvest"
	"database-api/handler/data"
	"database-api/handler/data/billing"
	"database-api/handler/data/checkout"
	"database-api/handler/data/insta"
	"database-api/handler/data/session"
	"database-api/handler/data/shipping"
	"database-api/handler/data/stores"
	"database-api/handler/data/webhook"
	"database-api/handler/data/whitelist"
	"database-api/handler/instance"
	"database-api/handler/instance/download"
	"database-api/handler/instance/payments"
	"database-api/handler/instance/performance"
	"database-api/handler/instance/update"
	"database-api/handler/linked_roles"
	"database-api/handler/login"
	"database-api/handler/monitor/instock"
	"database-api/handler/monitor/product"
	"database-api/handler/monitor/restart"
	"database-api/handler/newsletter"
	"database-api/handler/notifications"
	"database-api/handler/ping"
	"database-api/handler/purchase"
	"database-api/handler/subscriptions"
	"database-api/handler/user"
	"database-api/handler/websocket"
	"net/http"

	"github.com/Copped-Inc/aio-types/helper"
	"github.com/Copped-Inc/aio-types/requests"

	"github.com/gorilla/mux"
)

func Add() *mux.Router {

	pool := websocket.NewPool()
	go pool.Start()

	router := mux.NewRouter()
	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(requests.HandleCORS)
	router.Methods()

	router.HandleFunc("/", get).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/ping", ping.Handle).Methods(http.MethodGet)
	router.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) { websocket.Get(pool, w, r) }).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/discord/linked-roles", linked_roles.Handle).Methods(http.MethodGet, http.MethodOptions)

	router.HandleFunc("/login", login.Handle).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/newsletter", newsletter.Handle).Methods(http.MethodPost, http.MethodOptions)

	router.HandleFunc("/purchase", purchase.Handle).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/purchase/{code}", purchase.Handle).Methods(http.MethodGet, http.MethodOptions)

	router.HandleFunc("/stripe/customer-portal", subscriptions.Portal).Methods(http.MethodGet, http.MethodPost)
	subscriptions.Add(router.PathPrefix("/stripe/webhooks").Methods(http.MethodPost, http.MethodOptions).Subrouter())

	router.HandleFunc("/user", user.Handle).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/user/{id}", user.Handle).Methods(http.MethodGet, http.MethodPatch, http.MethodOptions)

	router.HandleFunc("/instance", instance.Handle).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	router.HandleFunc("/instance/performance", performance.Handle).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/instance/download/{file}", download.Handle).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/instance/update/{file}", update.Handle).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	router.HandleFunc("/instance/payments", payments.Handle).Methods(http.MethodPost, http.MethodOptions)

	router.HandleFunc("/data", data.Handle).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/data/instance", insta.Handle).Methods(http.MethodPost, http.MethodPatch, http.MethodDelete, http.MethodOptions)
	router.HandleFunc("/data/webhook", webhook.Handle).Methods(http.MethodPost, http.MethodDelete, http.MethodOptions)
	router.HandleFunc("/data/stores", stores.Handle).Methods(http.MethodPatch, http.MethodOptions)
	router.HandleFunc("/data/session", session.Handle).Methods(http.MethodPatch, http.MethodOptions)
	router.HandleFunc("/data/billing", billing.Handle).Methods(http.MethodPatch, http.MethodOptions)
	router.HandleFunc("/data/shipping", shipping.Handle).Methods(http.MethodPatch, http.MethodOptions)
	router.HandleFunc("/data/checkout", checkout.Handle).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/data/whitelist", whitelist.Handle).Methods(http.MethodPost, http.MethodDelete, http.MethodOptions)

	router.HandleFunc("/monitor/instock", instock.Handle).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/monitor/product/{handle}", product.Handle).Methods(http.MethodPost, http.MethodGet, http.MethodOptions)
	router.HandleFunc("/monitor/restart/{site}", restart.Handle).Methods(http.MethodGet, http.MethodOptions)

	router.HandleFunc("/notifications", notifications.Handle).Methods(http.MethodGet, http.MethodDelete, http.MethodPost, http.MethodOptions)
	router.HandleFunc("/notifications/{notification-id}", notifications.Handle).Methods(http.MethodGet, http.MethodDelete, http.MethodPatch, http.MethodOptions)
	router.HandleFunc("/notifications/users/@me", notifications.Handle).Methods(http.MethodGet, http.MethodDelete, http.MethodPut, http.MethodOptions)
	router.HandleFunc("/notifications/users/{user-id}", notifications.Handle).Methods(http.MethodGet, http.MethodDelete, http.MethodPost, http.MethodPut, http.MethodOptions)
	router.HandleFunc("/notifications/{notification-id}/users/@me", notifications.Handle).Methods(http.MethodGet, http.MethodDelete, http.MethodPatch, http.MethodOptions)
	router.HandleFunc("/notifications/{notification-id}/users/{user-id}", notifications.Handle).Methods(http.MethodGet, http.MethodDelete, http.MethodPatch, http.MethodOptions)

	router.HandleFunc("/captcha/challenge", challenge.Handle).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/captcha/preharvest", preharvest.Handle).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/captcha/preharvest/{task-id}", preharvest.Handle).Methods(http.MethodGet, http.MethodPatch, http.MethodDelete, http.MethodOptions)
	router.HandleFunc("/captcha/preharvest/user/{user-id}", preharvest.Handle).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	router.HandleFunc("/captcha/{site}", captcha.Handle).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, helper.Active+"/utility/404", http.StatusSeeOther)
	})

	http.Handle("/", router)

	return router

}
