package captcha

import (
	"bytes"
	"encoding/json"
	"github.com/Copped-Inc/aio-types/console"
	"monitor-api/api"
	"monitor-api/request"
	"strings"
)

func Start() error {

	res, err := request.PostForm("https://gem-fs.global-e.com/1/Checkout/GetCartToken?merchantUniqueId=708", bytes.NewBuffer([]byte("merchantUniqueId: 708\nMerchantCartToken=a5f47c33e8469795375691d8c3c77254&CountryCode=DE&CurrencyCode=EUR&CultureCode=de&MerchantId=708&GetCartTokenUrl=https%3A%2F%2Fgem-fs.global-e.com%2F1&ClientCartContent=%7B%22token%22%3A%22a5f47c33e8469795375691d8c3c77254%22%2C%22note%22%3Anull%2C%22attributes%22%3A%7B%7D%2C%22original_total_price%22%3A22500%2C%22total_price%22%3A22500%2C%22total_discount%22%3A0%2C%22total_weight%22%3A1360.7771%2C%22item_count%22%3A1%2C%22items%22%3A%5B%7B%22id%22%3A39663837839437%2C%22properties%22%3A%7B%22upsell%22%3A%22mens%22%7D%2C%22quantity%22%3A1%2C%22variant_id%22%3A39663837839437%2C%22key%22%3A%2239663837839437%3A244f6ca7cbe78615be2fd8a2557e6b69%22%2C%22title%22%3A%22Kith%20Nylon%20Fulton%20Kimono%20Track%20Jacket%20-%20Canvas%20-%20M%22%2C%22price%22%3A22500%2C%22original_price%22%3A22500%2C%22discounted_price%22%3A22500%2C%22line_price%22%3A22500%2C%22original_line_price%22%3A22500%2C%22total_discount%22%3A0%2C%22discounts%22%3A%5B%5D%2C%22sku%22%3A%2214287967%22%2C%22grams%22%3A1361%2C%22vendor%22%3A%22Kith%22%2C%22taxable%22%3Atrue%2C%22product_id%22%3A6738789859405%2C%22product_has_only_default_variant%22%3Afalse%2C%22gift_card%22%3Afalse%2C%22final_price%22%3A22500%2C%22final_line_price%22%3A22500%2C%22url%22%3A%22%2Fproducts%2Fkhm030433-210%3Fvariant%3D39663837839437%22%2C%22featured_image%22%3A%7B%22aspect_ratio%22%3A1%2C%22alt%22%3A%22Kith%20Nylon%20Fulton%20Kimono%20Track%20Jacket%20-%20Canvas%22%2C%22height%22%3A2000%2C%22url%22%3A%22https%3A%2F%2Fcdn.shopify.com%2Fs%2Ffiles%2F1%2F0274%2F7469%2F0125%2Fproducts%2FKHM030433-210-FRONT.jpg%3Fv%3D1675929025%22%2C%22width%22%3A2000%7D%2C%22image%22%3A%22https%3A%2F%2Fcdn.shopify.com%2Fs%2Ffiles%2F1%2F0274%2F7469%2F0125%2Fproducts%2FKHM030433-210-FRONT.jpg%3Fv%3D1675929025%22%2C%22handle%22%3A%22khm030433-210%22%2C%22requires_shipping%22%3Atrue%2C%22product_type%22%3A%22Outerwear%22%2C%22product_title%22%3A%22Kith%20Nylon%20Fulton%20Kimono%20Track%20Jacket%20-%20Canvas%22%2C%22product_description%22%3A%22Wrinkle%20nylon%20fabric%5CnButton-out%20interior%20track%20neck%20construction%5CnTrapunto%20stitch%20design%20details%20at%20front%20placket%5CnKith%20branded%20excella%20zipper%20at%20front%5CnHidden%20snap%20closure%20at%20patch%20pockets%5CnEmbroidered%20Kith%20serif%20logo%5Cn%20%5CnStyle%3A%20khm030433-210%5CnColor%3A%20Canvas%5CnMaterial%3A%20Nylon%22%2C%22variant_title%22%3A%22M%22%2C%22variant_options%22%3A%5B%22M%22%5D%2C%22options_with_values%22%3A%5B%7B%22name%22%3A%22Size%22%2C%22value%22%3A%22M%22%7D%5D%2C%22line_level_discount_allocations%22%3A%5B%5D%2C%22line_level_total_discount%22%3A0%7D%5D%2C%22requires_shipping%22%3Atrue%2C%22currency%22%3A%22EUR%22%2C%22items_subtotal_price%22%3A22500%2C%22cart_level_discount_applications%22%3A%5B%5D%7D&AdditionalCartData=%255B%255D")))
	if err != nil {
		return err
	}

	var resp response
	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		return err
	}

	if resp.IsCaptcha {
		sitekey := strings.Split(strings.Split(resp.HtmlInject, "data-sitekey=\"")[1], "\"")[0]
		if sitekey != api.KithEuSiteKey {
			api.KithEuSiteKey = sitekey
			console.Log("KithEU new Sitekey", sitekey)
		}
	}

	return err

}
