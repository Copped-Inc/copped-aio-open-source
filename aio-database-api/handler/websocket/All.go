package websocket

import (
	"database-api/user"
)

var b = Body{id: "all"}

func Monitor(p interface{}, l []string) {

	b.Op = NewProduct
	b.Data = p
	b.list = l

	Broadcast <- b
}

func Update() {

	b.Op = NewUpdate

	Broadcast <- b
}

func NotificationCreate(d user.Notification) {

	b.Op = NewNotification
	b.Data = d

	Broadcast <- b
}

func NotificationUpdate(d user.Notification) {

	b.Op = UpdatedNotification
	b.Data = d

	Broadcast <- b
}

func NotificationDelete(d user.Notification) {

	b.Op = DeletedNotification

	if d.ID != "" {
		b.Data = d
	}

	Broadcast <- b
}
