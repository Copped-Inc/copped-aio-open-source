package main

import (
	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/helper"
	"instances-api/handler"
	"instances-api/handler/update"
	"net/http"
)

var port = "94"

func main() {

	console.Log("Initialize", "Url", "instances.copped-inc.com")
	console.Log("Initialize", "Port", port)
	helper.Set("instances.copped-inc.com")

	router := handler.Add()

	go console.Loop()

	if helper.System == "linux" {
		err := update.Update()
		if err != nil {
			console.Log("Initialize", "Error", err.Error())
			return
		}
	}

	console.Log("Initialize", "Finished", "Listen and Serve")
	console.Log(http.ListenAndServe(":"+port, router))

}
