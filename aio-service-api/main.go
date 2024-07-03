package main

import (
	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/helper"
	"math/rand"
	"net/http"
	"service-api/handler"
	"service-api/handler/discord/interactions"
	"service-api/handler/discord/linked_roles"
	"service-api/ping"
	"time"
)

var port = "93"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {

	console.Log("Initialize", "Url", "service.copped-inc.com")
	console.Log("Initialize", "Port", port)
	helper.Set("service.copped-inc.com")

	go ping.Start()
	go interactions.Update()
	go linked_roles.RegisterMetadata()
	router := handler.Add()

	go console.Loop()
	console.Log("Initialize", "Finished", "Listen and Serve")
	console.Log((&http.Server{Handler: router, Addr: ":" + port, WriteTimeout: time.Minute, ReadTimeout: time.Minute}).ListenAndServe())

}
