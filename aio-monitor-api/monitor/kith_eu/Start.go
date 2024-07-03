package kith_eu

import (
	"encoding/json"
	"errors"
	"fmt"
	"monitor-api/request"
)

var products map[int64]string

func Start() error {

	res, err := request.Get("https://eu.kith.com/products.json?limit=250", nil, nil)
	if err != nil {
		return err
	}

	var newProducts response
	err = json.NewDecoder(res.Body).Decode(&newProducts)
	if err != nil {
		return errors.New(fmt.Sprintf("shopify error: %s", err.Error()))
	}

	analyse(newProducts.Products)
	return err

}
