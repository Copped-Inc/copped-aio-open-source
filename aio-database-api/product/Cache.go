package product

import (
	"database-api/database"

	"github.com/Copped-Inc/aio-types/console"
	"github.com/Copped-Inc/aio-types/webhook"
)

func UpdateCache(p Product) {

	doc, err := database.GetDatabase().Collection("products").Doc("cache").Get(database.GetContext())
	if err != nil {
		console.ErrorLog(err)
		return
	}

	var c Cache
	err = doc.DataTo(&c)
	if err != nil {
		console.ErrorLog(err)
		return
	}

	for i, v := range c.Whitelist {
		if v.SKU == p.SKU {
			c.Whitelist = append(c.Whitelist[:i], c.Whitelist[i+1:]...)
			break
		}
	}

	if p.State == Whitelisted {
		c.Whitelist = append(c.Whitelist, p)
	}

	_, err = database.GetDatabase().Collection("products").Doc("cache").Set(database.GetContext(), c)
	if err != nil {
		console.ErrorLog(err)
		return
	}

	wh := webhook.New().AddEmbed(webhook.Whitelist)
	for _, product := range c.Whitelist {
		add := func() (s string) {
			s += "["
			if product.Name != "" {
				s += product.Name + " - " + product.SKU
			} else {
				s += product.SKU
			}

			s += "]("
			if product.StockX != "" {
				s += product.StockX
			} else {
				s += "https://stockx.com/search?s=" + product.SKU
			}
			s += ")\n"

			return s
		}()

		if len(wh.Embeds[len(wh.Embeds)-1].Description)+len(add) > 1024 {
			wh.AddEmbed(webhook.Whitelist)
		}

		wh.Embeds[len(wh.Embeds)-1].Description += add
	}

	err = wh.Update("") // INSERT Webhook URL here
	if err != nil {
		console.ErrorLog(err)
		return
	}

}

type Cache struct {
	Whitelist []Product `json:"whitelist"`
}
