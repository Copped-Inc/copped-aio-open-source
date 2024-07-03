package kith_eu_stock

import (
	"encoding/json"
	"errors"
	"fmt"
	"monitor-api/request"
)

func Start(keyword string) func() error {
	var products map[string]Item
	return func() error {
		res, err := request.Get("https://searchserverapi.com/getwidgets?api_key=3c7s6k4F2C&q="+keyword+"&items=true&_=1676966966392&maxResults=550&ci=1", nil, nil)
		if err != nil {
			return err
		}

		var newProducts response
		err = json.NewDecoder(res.Body).Decode(&newProducts)
		if err != nil {
			return errors.New(fmt.Sprintf("kith_eu_stock error: %s", err.Error()))
		}

		products = analyse(newProducts.Items, products)
		return err
	}
}
