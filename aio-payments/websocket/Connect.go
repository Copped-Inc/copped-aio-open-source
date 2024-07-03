package websocket

import (
	"github.com/gorilla/websocket"
	"time"
)

func (c *Client) Connect() {

	conn, _, err := websocket.DefaultDialer.Dial("wss://database.copped-inc.com/websocket", c.Header)
	if err != nil {
		c.Conn = nil
		time.Sleep(time.Second * 5)
		c.Connect()
		return
	}

	c.Conn = conn
	go Read(c)

}
