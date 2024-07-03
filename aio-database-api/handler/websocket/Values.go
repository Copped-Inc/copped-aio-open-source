package websocket

import (
	"database-api/user"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	Login    time.Time
	Conn     *websocket.Conn
	Conns    map[*websocket.Conn]*user.Instance
	Pool     *Pool
	User     user.User
	Instance *user.Instance
}

type Body struct {
	id   string
	list []string
	Op   WsOp        `json:"op"`
	Data interface{} `json:"data"`
}

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[string]*Client
}

type WsOp int

const (
	Ping WsOp = iota + 1
	DataUpdate
	NewQueueRequest
	SendData
	Payments
	NewProduct
	NewUpdate
	NewNotification
	UpdatedNotification
	DeletedNotification
	InstanceLog
)
