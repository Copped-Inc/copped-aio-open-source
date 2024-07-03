package link

import (
	"database-api/database"

	"github.com/Copped-Inc/aio-types/subscriptions"
)

func New() *Link {
	return &Link{}
}

func (l *Link) SetPlan(plan subscriptions.Plan) *Link {
	l.Plan = plan
	return l
}

func (l *Link) SetStock(stock int) *Link {
	l.Stock = stock
	return l
}

func (l *Link) Use() *Link {
	l.Stock--
	return l
}

func (l *Link) SetInstanceLimit(limit int) *Link {
	l.InstanceLimit = limit
	return l
}

func (l *Link) Create() error {
	doc, _, err := database.GetDatabase().Collection("link").Add(database.GetContext(), l)
	if err == nil {
		l.ID = doc.ID
	}
	return err
}

func (l *Link) Update() error {
	_, err := database.GetDatabase().Collection("link").Doc(l.ID).Set(database.GetContext(), l)
	return err
}
