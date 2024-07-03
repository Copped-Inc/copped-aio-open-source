package product

import (
	"database-api/database"
)

func GlobalWhitelisted() ([]string, error) {

	doc, err := database.GetDatabase().Collection("products").Doc("cache").Get(database.GetContext())
	if err != nil {
		return nil, err
	}

	var c Cache
	err = doc.DataTo(&c)
	if err != nil {
		return nil, err
	}

	var skus []string
	for _, v := range c.Whitelist {
		if v.State == Whitelisted {
			skus = append(skus, v.SKU)
		}
	}

	return skus, err

}
