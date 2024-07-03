package websocket

import (
	"database-api/user"
	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/webhook"
	"github.com/gorilla/websocket"
	"golang.org/x/exp/slices"
	"time"
)

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client, 32),
		Unregister: make(chan *Client, 32),
		Clients:    make(map[string]*Client),
	}
}

var Broadcast = make(chan Body, 32)

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			if c, ok := pool.Clients[client.User.ID]; ok {
				c.Conns[client.Conn] = client.Instance
			} else {
				client.Conns = make(map[*websocket.Conn]*user.Instance)
				client.Conns[client.Conn] = client.Instance
				pool.Clients[client.User.ID] = client
			}
		case client := <-pool.Unregister:
			if c, ok := pool.Clients[client.User.ID]; ok {
				delete(c.Conns, client.Conn)
				if len(c.Conns) == 0 {
					delete(pool.Clients, client.User.ID)
				}

				if client.Instance == nil {
					break
				}

				go func() {
					time.Sleep(time.Second * 5)
					if c, ok := pool.Clients[client.User.ID]; ok {
						for _, instance := range c.Conns {
							if instance != nil && instance.ID == client.Instance.ID {
								return
							}
						}
					}

					d, err := user.FromId(client.User.ID)
					if err != nil {
						console.ErrorLog(err)
						return
					}

					for _, instance := range d.Instances {
						id := instance.ID
						if id != client.Instance.ID {
							continue
						}

						client.Instance.Status = "Offline"
						err = d.UpdateInstance(*client.Instance).Update()
						if err != nil {
							console.ErrorLog(err)
							break
						}

						Websocket{
							Action: UpdateInstances,
							Body:   d.Instances,
						}.Send(client.User.ID)

						Websocket{
							Action: UpdateSession,
							Body: struct {
								Session user.Session `json:"session"`
							}{
								Session: *d.Session,
							},
						}.Send(client.User.ID)

						webhook.New().AddEmbed(webhook.InstanceLogout, id).SendMultiple(d.Settings.Webhooks)
					}
				}()
			}
		case message := <-Broadcast:
			if message.id == "all" {
				for _, client := range pool.Clients {
					if slices.Contains(message.list, client.User.ID) {
						continue
					}

					for conn := range client.Conns {
						_ = conn.WriteJSON(message)
					}
				}
			} else if c, ok := pool.Clients[message.id]; ok {
				for conn := range c.Conns {
					if err := conn.WriteJSON(message); err != nil {
						delete(c.Conns, conn)
						if len(c.Conns) == 0 {
							delete(pool.Clients, message.id)
						}
					}
				}
			}
		}
	}
}
