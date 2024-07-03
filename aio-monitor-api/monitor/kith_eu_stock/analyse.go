package kith_eu_stock

import (
	"github.com/Copped-Inc/aio-types/console"
	"strconv"
	"strings"
)

var Keywords = []string{"nike", "adidas", "jordan", "yeezy"}

func analyse(res []Item, products map[string]Item) map[string]Item {

	newProducts := make(map[string]Item)
	for _, re := range res {
		stock := 0
		for _, variant := range re.ShopifyVariants {
			if variant.QuantityTotal != "" {
				s, err := strconv.Atoi(variant.QuantityTotal)
				if err != nil {
					continue
				}

				stock += s
			}
		}

		if stock == 0 {
			continue
		}

		re.stock = stock
		newProducts[re.ProductId] = re
		if products[re.ProductId].stock == re.stock || len(products) == 0 {
			continue
		}

		for _, keyword := range Keywords {
			if strings.Contains(strings.ToLower(re.Title), keyword) {
				console.Log("Kith Found " + re.Title)
				go re.send()
				break
			}
		}
	}
	return newProducts

}
