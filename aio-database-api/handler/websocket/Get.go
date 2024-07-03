package websocket

import (
	"database-api/user"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/helper"
	"github.com/Copped-Inc/aio-types/statistic"
	"github.com/Copped-Inc/aio-types/webhook"
)

func Get(pool *Pool, w http.ResponseWriter, r *http.Request) {

	if r.UserAgent() == "Monitor" && helper.IsMaster(r.Header.Get("Password")) {
		ws, err := Upgrade(w, r)
		if err != nil {
			console.ErrorRequest(w, r, err, http.StatusInternalServerError)
			return
		}

		err = ws.WriteJSON(Body{Op: Ping, Data: []byte("ping")})
		if err != nil {
			console.ErrorRequest(w, r, err, http.StatusInternalServerError)
			return
		}

		client := &Client{
			Conn: ws,
			Pool: pool,
			User: user.User{
				ID: "monitor",
			},
		}

		pool.Register <- client
		statistic.Status(r, http.StatusSwitchingProtocols)
		go client.Read()
		return
	}

	d, err := user.FromRequest(r)
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusUnauthorized)
		return
	}

	if c, ok := pool.Clients[d.User.ID]; ok && (r.UserAgent() == "Instance" || r.UserAgent() == "Payments") {
		for _, i := range c.Conns {
			if i != nil && i.ID == r.Header.Get("ID") {
				console.ErrorRequest(w, r, errors.New("instance is already online"), http.StatusNotAcceptable)
				return
			}
		}
	}

	ws, err := Upgrade(w, r)
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	err = ws.WriteJSON(Body{Op: Ping, Data: []byte("ping")})
	if err != nil {
		console.ErrorRequest(w, r, err, http.StatusInternalServerError)
		return
	}

	client := &Client{
		Login: time.Now(),
		Conn:  ws,
		Pool:  pool,
		User:  d.User,
	}

	if r.UserAgent() == "Instance" || r.UserAgent() == "Payments" {
		price, err := strconv.ParseFloat(r.Header.Get("Price"), 64)
		if err != nil {
			console.ErrorRequest(w, r, err, http.StatusBadRequest)
			return
		}

		i := user.Instance{
			Price:    price,
			Provider: r.Header.Get("Provider"),
			ID:       r.Header.Get("ID"),
			Status:   "Running",
			TaskMax:  r.Header.Get("Task-Max"),
			Region:   r.Header.Get("Region"),
		}

		i, err = func() (user.Instance, error) {
			if r.Header.Get("reconnect") == "" && r.UserAgent() == "Instance" {
				err = d.NeedPassword(true).Update()
				if err != nil {
					return i, err
				}

				go webhook.New().AddEmbed(webhook.DataRequest, i.ID).SendMultiple(d.Settings.Webhooks)
				return i, d.UpdateInstance(i).Update()
			}
			return i, d.UpdateInstance(i).Update()
		}()
		client.Instance = &i

		if err != nil {
			console.ErrorRequest(w, r, err, http.StatusInternalServerError)
			return
		}

		Websocket{
			Action: UpdateInstances,
			Body:   d.Instances,
		}.Send(d.User.ID)

		Websocket{
			Action: UpdateSession,
			Body: struct {
				Session user.Session `json:"session"`
			}{
				Session: *d.Session,
			},
		}.Send(d.User.ID)

		if r.Header.Get("reconnect") == "" {
			b := Body{
				id: d.User.ID,
				Op: NewQueueRequest,
			}
			Broadcast <- b

		}
	}

	statistic.Status(r, http.StatusSwitchingProtocols)
	pool.Register <- client
	go client.Read()

}
