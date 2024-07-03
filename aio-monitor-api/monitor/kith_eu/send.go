package kith_eu

import (
	"monitor-api/api"
	"strconv"
	"time"
)

func (p product) send() {

	var skus = make(map[string]string)
	for i := 0; i < len(p.Variants); i++ {
		if p.Variants[i].Available {
			skus[p.Variants[i].Title] = strconv.Itoa(int(p.Variants[i].Id))
		}
	}

	if len(p.Variants) == 0 {
		return
	}

	price, err := strconv.ParseFloat(p.Variants[0].Price, 64)
	if err != nil {
		return
	}

	if len(p.Images) == 0 {
		return
	}

	api.InstockReq{
		Name:  p.Title,
		Sku:   p.Handle[2:],
		Skus:  skus,
		Date:  time.Now(),
		Link:  "https://eu.kith.com/products/" + p.Handle,
		Image: p.Images[0].Src,
		Price: price,
	}.Send()

}
