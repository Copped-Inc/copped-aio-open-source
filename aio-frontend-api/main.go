package main

import (
	"frontend-api/handler"
	"net/http"

	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/helper"
)

var port = "90"

func main() {

	console.Log("Initialize", "Url", "aio.copped-inc.com")
	console.Log("Initialize", "Port", port)
	helper.Set("aio.copped-inc.com")

	router := handler.Add()

	go console.Loop()

	console.Log("Initialize", "Finished", "Listen and Serve")
	console.Log(http.ListenAndServe(":"+port, router))

}
