package kith_eu_stock

type response struct {
	Items []Item `json:"items"`
}

type Item struct {
	stock           int
	ProductId       string `json:"product_id"`
	Title           string `json:"title"`
	Link            string `json:"link"`
	Price           string `json:"price"`
	ImageLink       string `json:"image_link"`
	ShopifyVariants []struct {
		VariantId string `json:"variant_id"`
		Options   struct {
			Size string `json:"Size"`
		} `json:"options"`
		QuantityTotal string `json:"quantity_total"`
	} `json:"shopify_variants"`
}
