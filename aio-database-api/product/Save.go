package product

import (
	"database-api/database"
)

func (p *Product) Save() error {

	go UpdateCache(*p)
	_, err := database.GetDatabase().Collection("products").Doc(p.SKU).Set(database.GetContext(), p)
	return err

}
