package kith_eu

import (
	"github.com/Copped-Inc/aio-types/console"
	"monitor-api/monitor/kith_eu_stock"
	"strings"
)

func analyse(res []product) {

	newProducts := make(map[int64]string)
	for _, re := range res {
		newProducts[re.Id] = func() string {
			sizes := ""
			for j := 0; j < len(re.Variants); j++ {
				if re.Variants[j].Available {
					sizes += "+" + re.Variants[j].Title
				}
			}
			return sizes
		}()

		if len(products) != 0 && products[re.Id] != newProducts[re.Id] {
			for _, keyword := range kith_eu_stock.Keywords {
				if strings.Contains(strings.ToLower(re.Title), keyword) {
					console.Log("Backup Found " + re.Title)
					go re.send()
					break
				}
			}
		}
	}

	products = newProducts

}
