package kith_eu_stock

import (
	"github.com/Copped-Inc/aio-types/console"
	"monitor-api/api"
	"strconv"
	"strings"
	"time"
)

func (item Item) send() {

	console.Log("Sending product " + item.Link)
	var skus = make(map[string]string)
	for i := 0; i < len(item.ShopifyVariants); i++ {
		if item.ShopifyVariants[i].QuantityTotal != "" {
			skus[item.ShopifyVariants[i].Options.Size] = item.ShopifyVariants[i].VariantId
		}
	}

	if len(item.ShopifyVariants) == 0 {
		return
	}

	price, err := strconv.ParseFloat(item.Price, 64)
	if err != nil {
		return
	}

	if len(item.ImageLink) == 0 {
		return
	}
	split := strings.Split(item.Link, "/")

	api.InstockReq{
		Name:  item.Title + " [Stock: " + strconv.Itoa(item.stock) + "]",
		Sku:   split[len(split)-1][2:],
		Skus:  skus,
		Date:  time.Now(),
		Link:  "https://eu.kith.com" + item.Link,
		Image: item.ImageLink,
		Price: price,
	}.Send()

}
