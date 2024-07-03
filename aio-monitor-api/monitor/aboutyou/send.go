package aboutyou

import (
	"monitor-api/api"
	"strconv"
	"time"
)

func (p product) send() {

	var skus = make(map[string]string)
	stock := 0
	for i := 0; i < len(p.Variants); i++ {
		if p.Variants[i].Stock.Quantity > 0 {
			skus[strconv.Itoa(p.Variants[i].Id)] = strconv.Itoa(p.Variants[i].Id)
			stock += p.Variants[i].Stock.Quantity
		}
	}

	if len(p.Variants) == 0 {
		return
	}

	price := float64(p.Variants[0].Price.WithTax) / 100

	if len(p.Images) == 0 {
		return
	}

	api.InstockReq{
		Name:  p.Attributes.Name.Values.Label + " [Stock: " + strconv.Itoa(stock) + "]",
		Sku:   strconv.Itoa(p.Id),
		Skus:  skus,
		Date:  time.Now(),
		Link:  "https://www.aboutyou.com/p/ci/ci-" + strconv.Itoa(p.Id),
		Image: "https://cdn.aboutyou.de/" + p.Images[0].Hash,
		Price: price,
	}.Send()

}
