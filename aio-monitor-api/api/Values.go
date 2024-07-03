package api

import "time"

type InstockReq struct {
	Name  string            `json:"name"`
	Sku   string            `json:"sku"`
	Skus  map[string]string `json:"skus"`
	Date  time.Time         `json:"date"`
	Link  string            `json:"link"`
	Image string            `json:"image"`
	Price float64           `json:"price"`
}
