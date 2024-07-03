package websocket

import (
	"github.com/gorilla/websocket"
	"net/http"
)

type Client struct {
	Header http.Header
	Conn   *websocket.Conn
}

type Body struct {
	Op   int         `json:"op"`
	Data interface{} `json:"data"`
}
