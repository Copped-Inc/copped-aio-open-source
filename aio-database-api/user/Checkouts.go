package user

import (
	"database-api/database"
	"errors"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

func (d *Database) GetCheckouts() (products []Product, err error) {
	products = []Product{}
	iter := database.GetDatabase().Collection("checkouts").Where("user", "==", d.User.ID).OrderBy("date", firestore.Desc).Limit(50).Documents(database.GetContext())

	for {
		var item *firestore.DocumentSnapshot
		item, err = iter.Next()
		if err != nil {
			if errors.Is(err, iterator.Done) {
				err = nil
				break
			}
			return
		}

		var product Product

		if err = item.DataTo(&product); err != nil {
			return
		}

		products = append(products, product)
	}

	return
}

func (d *Database) AddCheckout(p Product) error {
	p.User = d.User.ID
	if _, _, err := database.GetDatabase().Collection("checkouts").Add(database.GetContext(), p); err != nil {
		return err
	}

	d.CheckoutAmount++

	return d.Update()
}
