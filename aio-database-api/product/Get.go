package product

import (
	"database-api/database"
	"errors"
	"strings"
)

func Get(sku string) (product Product, err error) {

	doc, err := database.GetDatabase().Collection("products").Doc(strings.ToLower(sku)).Get(database.GetContext())
	if err != nil {
		return GetFromHandle(sku)
	}

	err = doc.DataTo(&product)
	return product, err

}

func GetFromHandle(handle string) (product Product, err error) {

	docs, err := database.GetDatabase().Collection("products").Where("handles", "array-contains", strings.ToLower(handle)).Documents(database.GetContext()).GetAll()
	if err != nil || len(docs) == 0 {
		return product, errors.New("product not found")
	}

	err = docs[0].DataTo(&product)
	return product, err

}
