package websocket

import (
	"database-api/log"
	"database-api/user"
	"encoding/json"
	"reflect"
	"time"

	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/webhook"
)

func (c *Client) Read() {
	defer func() {
		console.Log("Client", "logged out User: "+c.User.ID)
		c.Pool.Unregister <- c
		_ = c.Conn.Close()
	}()

	for {
		var message Body
		err := c.Conn.ReadJSON(&message)
		if err != nil {
			console.ErrorLog(err)
			return
		}

		switch message.Op {
		case Ping:
			go func() {
				if c.Instance != nil {
					time.Sleep(time.Second * 20)
				}

				err = c.Conn.WriteJSON(message)
				if err != nil {
					console.ErrorLog(err)
				}
			}()

		case NewQueueRequest:
			if message.Data == nil || reflect.TypeOf(message.Data).Kind() != reflect.String {
				continue
			}

			d, db, err := user.DataFromWebsocket(c.User.ID, message.Data.(string))
			if err != nil {
				continue
			}

			data, err := db.ToDataResp(d)
			if err != nil {
				continue
			}

			if db.PassNeeded {
				b := Body{
					id:   db.User.ID,
					Op:   SendData,
					Data: data,
				}
				Broadcast <- b
				err = db.NeedPassword(false).Update()
				if err != nil {
					console.ErrorLog(err)
				}

				if db.Settings != nil && db.Settings.Webhooks != nil {
					go webhook.New().AddEmbed(webhook.DataReceived).SendMultiple(db.Settings.Webhooks)
				}
			}

		case InstanceLog:
			b, err := json.Marshal(message.Data)
			if err != nil {
				console.ErrorLog(err)
				continue
			}

			var l []log.Log
			err = json.Unmarshal(b, &l)
			if err != nil {
				console.ErrorLog(err)
				continue
			}

			go log.Handle(c.User.ID, c.Instance.ID, l)

		default:
			console.Log(message.Data.(string))
		}
	}
}
