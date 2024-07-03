package websocket

import (
	"changeme/payments"
	"encoding/json"
	"time"
)

func Read(c *Client) {
	println("Reading")
	for {
		var message Body
		err := c.Conn.ReadJSON(&message)
		if err != nil {
			c.Connect()
			return
		}

		switch message.Op {
		case 1:
			go func() {
				time.Sleep(time.Second * 20)
				if c.Conn == nil {
					return
				}

				err = c.Conn.WriteJSON(message)
				if err != nil {
					c.Connect()
					return
				}
			}()
		case 5:
			var b []byte
			b, err = json.Marshal(message.Data)
			if err != nil {
				continue
			}

			var p *payments.Payments
			err = json.Unmarshal(b, &p)
			if err != nil {
				continue
			}

			println("Payment Received", p.Id, p.State)
			if p.State != payments.Created {
				continue
			}

			go p.Handle(c.Header["Cookie"][0])
		}
	}
}
