package aboutyou

import (
	"encoding/json"
	"errors"
	"fmt"
	"monitor-api/request"
)

var lastProduct product

func Start() error {

	res, err := request.Get("https://api-cloud.aboutyou.de/v1/products?with=attributes%3Akey%28name%7Cbrand%29%2Cvariants&perPage=100&sortBy=updatedAt", nil, nil)
	if err != nil {
		return err
	}

	var newProducts response
	err = json.NewDecoder(res.Body).Decode(&newProducts)
	if err != nil {
		return errors.New(fmt.Sprintf("aboutyou error: %s", err.Error()))
	}

	analyse(newProducts.Entities)
	return err

}
