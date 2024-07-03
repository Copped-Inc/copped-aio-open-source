package product

import "strings"

func New(sku string, state State) Product {

	return Product{
		SKU:   strings.ToLower(sku),
		State: state,
	}

}
