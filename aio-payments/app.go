package main

import (
	"changeme/settings"
	"changeme/websocket"
	"context"
	"fmt"
	"net/http"
)

type App struct {
	ctx      context.Context
	settings *settings.Settings
	ws       *websocket.Client
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {

	println("Starting up")

	a.ctx = ctx
	s, err := settings.Get()
	if err != nil {
		return
	}

	a.settings = s

}

func (a *App) Loggedin() bool {

	println("Checking if logged in")

	if a.Check() {
		return true
	}

	if a.settings == nil {
		return false
	}

	err := a.settings.Update()
	if err != nil {
		panic(err)
	}

	a.connect()
	return a.Check()

}

func (a *App) Login(code string) bool {

	println("Logging in")

	s, err := settings.Login(code)
	if err != nil {
		return false
	}

	a.settings = s
	err = a.settings.Update()
	if err != nil {
		panic(err)
	}

	a.connect()
	return true

}

func (a *App) Check() bool {
	return a.ws != nil && a.ws.Conn != nil
}

func (a *App) connect() {

	println("Connecting to websocket")

	client := &websocket.Client{
		Header: http.Header{
			"User-Agent": []string{"Payments"},
			"Price":      []string{fmt.Sprintf("%f", a.settings.Price)},
			"Provider":   []string{a.settings.Provider},
			"Task-Max":   []string{fmt.Sprintf("%d", a.settings.TaskMax)},
			"Region":     []string{a.settings.Region},
			"Id":         []string{a.settings.Id},
			"Cookie":     []string{"authorization=" + a.settings.Authorization},
		},
	}

	a.ws = client
	client.Connect()

}
