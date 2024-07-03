package main

import (
	"database-api/database"
	"database-api/handler"
	"database-api/handler/instance/update"
	"database-api/preharvest"
	"math/rand"
	"net/http"
	"time"

	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/helper"
	"github.com/Copped-Inc/aio-types/secrets"
	"github.com/stripe/stripe-go/v74"
)

var port = "91"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {

	stripe.Key = secrets.Stripe_Secret

	console.Log("Initialize", "Url", "database.copped-inc.com")
	console.Log("Initialize", "Port", port)
	console.Log("Initialize", "DB", "Connecting")

	if err := database.Connect(); err != nil {
		console.Log("Initialize", "DB", "Error", err.Error())
		return
	}

	console.Log("Initialize", "DB", "Connected")

	helper.Set("database.copped-inc.com")
	go console.Loop()

	if helper.System == "linux" {
		err := update.DownloadAll()
		if err != nil {
			console.Log("Initialize", "Update", "Error", err.Error())
			return
		}
	}

	go func() {
		if err := preharvest.Initialize(); err != nil {
			console.Log("Initialize", "Captcha preharvest", "loading captcha preharvest tasks failed with:", err.Error())
			return
		}

		console.Log("Initialize", "Captcha preharvest", "Successfully initiated captcha preharvest tasks.")
	}()

	console.Log("Initialize", "Finished", "Listen and Serve")
	console.Log((&http.Server{Handler: handler.Add(), Addr: ":" + port, WriteTimeout: time.Minute, ReadTimeout: time.Minute}).ListenAndServe())
}
