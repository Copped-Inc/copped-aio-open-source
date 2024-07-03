package aboutyou

import "time"

type response struct {
	Entities []product `json:"entities"`
}

type product struct {
	Id         int       `json:"id"`
	IsActive   bool      `json:"isActive"`
	IsSoldOut  bool      `json:"isSoldOut"`
	IsNew      bool      `json:"isNew"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	Attributes struct {
		Brand struct {
			Values struct {
				Label string `json:"label"`
				Value string `json:"value"`
			} `json:"values"`
		} `json:"brand"`
		Name struct {
			Values struct {
				Label string `json:"label"`
				Value string `json:"value"`
			} `json:"values"`
		} `json:"name"`
	} `json:"attributes"`
	Images []struct {
		Hash string `json:"hash"`
	} `json:"images"`
	Variants []struct {
		Id    int `json:"id"`
		Stock struct {
			Quantity int `json:"quantity"`
		} `json:"stock"`
		Price struct {
			CurrencyCode      string        `json:"currencyCode"`
			WithTax           int           `json:"withTax"`
			AppliedReductions []interface{} `json:"appliedReductions"`
		} `json:"price"`
		CustomData struct {
		} `json:"customData"`
	} `json:"variants"`
	CustomData struct {
	} `json:"customData"`
}
