package websocket

func (w Websocket) Send(id string) {
	b := Body{
		id:   id,
		Op:   DataUpdate,
		Data: w,
	}
	Broadcast <- b
}

func UserMonitor(p interface{}, id string) {

	b := Body{
		id:   id,
		Op:   NewProduct,
		Data: p,
	}

	Broadcast <- b
}

func UserPayments(p interface{}, id string) {

	b := Body{
		id:   id,
		Op:   Payments,
		Data: p,
	}

	Broadcast <- b
}

type Websocket struct {
	Action WsAction    `json:"action"`
	Body   interface{} `json:"body"`
}

type WsAction int

const (
	AddWebhook WsAction = iota + 1
	DeleteWebhook
	UpdateStores
	UpdateInstances
	UpdateSession
	AddCheckout
	UpdateBilling
	UpdateShipping
	CreateNotification
	UpdateNotification
	DeleteNotification
	AddWhitelist
	RemoveWhitelist
	UpdateNotificationReadstate
)
