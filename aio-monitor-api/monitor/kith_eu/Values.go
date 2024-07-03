package kith_eu

type response struct {
	Products []product `json:"products"`
}

type product struct {
	Id       int64     `json:"id"`
	Title    string    `json:"title"`
	BodyHtml string    `json:"body_html"`
	Handle   string    `json:"handle"`
	Tags     []string  `json:"tags"`
	Variants []variant `json:"variants"`
	Images   []image   `json:"images"`
}

type variant struct {
	Id             int64  `json:"id"`
	Title          string `json:"title"`
	Available      bool   `json:"available"`
	Price          string `json:"price"`
	CompareAtPrice string `json:"compare_at_price"`
}

type image struct {
	Id  int64  `json:"id"`
	Src string `json:"src"`
}
