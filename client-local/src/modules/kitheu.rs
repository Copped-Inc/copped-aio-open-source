use regex::Regex;
use country_emoji::code;
use reqwest::StatusCode;
use serde::Deserialize;
use std::str;
use async_recursion::async_recursion;
use crate::captcha::{Captcha, get_captcha, NEED_CAPTCHA};
use crate::modules::Task;
use crate::request::ReqStruct;
use crate::jig::{country_id, random_name, random_number};
use crate::log;
use crate::modules::session::Session;
use crate::request::cookie::ReqCookie;
use crate::request::expected::Expected;

impl Task {
    pub(crate) async fn kith(self) {
        unsafe {
            NEED_CAPTCHA = true;
        }

        let s = &mut Session::get_from_task(self.clone());

        if s.state == 0 {
            let (checkout, cookies) = self.clone().cart().await;
            if checkout.is_empty() {
                s.save();
                return;
            }

            s.state = 1;
            s.checkout = checkout;
            s.cookies = cookies;
        }

        let c = get_captcha();
        if c.is_err() {
            s.save();
            return;
        }
        let c = c.unwrap();

        if s.state == 1 {
            self.clone().cart_token(s, c.clone()).await;
            if s.state == 1 {
                s.save();
                return;
            }
        }

        if s.state == 2 {
            self.clone().address(s).await;
            if s.state == 2 {
                s.save();
                return;
            }
        }

        if s.state == 3 {
            self.clone().card(s, c).await;
            if s.state == 3 {
                s.save();
                return;
            }
        }

        self.checked_out(s.clone()).await;
    }

    async fn cart(self) -> (String, Vec<ReqCookie>) {
        let body = vec![
            ("form_type".to_owned(), "product".to_owned()),
            ("utf8".to_owned(), "✓".to_owned()),
            ("properties[upsell]".to_owned(), "mens".to_owned()),
            ("properties[Size]".to_owned(), self.size.to_owned()),
            ("id".to_owned(), self.prod_id.to_owned()),
            ("quantity".to_owned(), "1".to_owned()),
        ];

        let req = ReqStruct::new_with_form("https://eu.kith.com/cart/add.js".to_string(), body, vec![], self.proxy.clone());
        let res = req.post(Expected::new_status(StatusCode::OK)).await;
        if let Ok(res) = res {
            return (res.body, res.cookies);
        }
        (String::new(), vec![])
    }

    #[async_recursion]
    async fn cart_token(self, s: &mut Session, c: Captcha) {
        let mut cart = "".to_string();
        for ck in &s.cookies {
            if ck.name == "cart" {
                cart = ck.value.clone();
            }
        }

        let re_price = Regex::new(r#""price":(.*?),"#).unwrap();
        let re_grams = Regex::new(r#""grams":(.*?),"#).unwrap();
        let price = re_price.find(&s.checkout);
        let weight = re_grams.find(&s.checkout);
        if price.is_none() || weight.is_none() {
            return;
        }

        let mut price = price.unwrap().as_str()[8..].to_string();
        price = price[..price.len() - 1].to_string();

        let mut weight = weight.unwrap().as_str()[8..].to_string();
        weight = weight[..weight.len() - 1].to_string() + ".7771";

        let body = vec![
            ("MerchantCartToken".to_owned(), cart.to_owned()),
            ("CaptchaResponseToken".to_owned(), c.captcha_one.clone()),
            ("CountryCode".to_owned(), code(self.shipping.country.clone().as_str()).unwrap().to_owned()),
            ("CurrencyCode".to_owned(), "EUR".to_owned()),
            ("CultureCode".to_owned(), code(self.shipping.country.clone().as_str()).unwrap().to_owned().to_lowercase()),
            ("MerchantId".to_owned(), "708".to_owned()),
            ("GetCartTokenUrl".to_owned(), "https://gem-fs.global-e.com/1".to_owned()),
            ("ClientCartContent".to_owned(), format!("{{\"token\":\"{}\",\"note\":null,\"attributes\":{{\"lang\":\"en\",\"Invoice Language\":\"en\"}},\"original_total_price\":{},\"total_price\":{},\"total_discount\":0,\"total_weight\":{},\"item_count\":1,\"items\":[{}],\"requires_shipping\":true,\"currency\":\"EUR\",\"items_subtotal_price\":{},\"cart_level_discount_applications\":[]}}", cart, price, price, weight, s.checkout, price).replace("\\/", "/").to_owned()),
            ("AdditionalCartData".to_owned(), "%5B%5D".to_owned()),
        ];

        let req = ReqStruct::new_with_form("https://gem-fs.global-e.com/1/Checkout/GetCartToken?merchantUniqueId=708".to_string(), body, vec![], self.proxy.clone());
        let res = req.post(Expected::new_status(StatusCode::OK)).await;
        if let Ok(res) = res {
            println!("{}", res.body.clone());
            let c:serde_json::Result<CartToken> = serde_json::from_slice((&res.body).as_ref());
            if c.is_ok() {
                s.state = 2;
                s.checkout = c.unwrap().cart_token;
            } else {
                let c2:serde_json::Result<CartToken2> = serde_json::from_slice((&res.body).as_ref());
                if c2.is_ok() {
                    s.state = 2;
                    s.checkout = c2.unwrap().cart_token;
                }
            }
        }
    }

    async fn address(self, s: &mut Session) {
        let id = country_id(code(self.shipping.country.clone().as_str()).unwrap().to_owned());
        let body = vec![
            ("CheckoutData.CartToken".to_owned(), s.checkout.clone().to_owned()),
            ("CheckoutData.CultureID".to_owned(), "2057".to_owned()),
            ("CheckoutData.GASessionsID".to_owned(), "632716701.625364308.708".to_owned()),
            ("CheckoutData.IsVirtualOrder".to_owned(), "False".to_owned()),
            ("CheckoutData.ExternalData.CurrentGatewayId".to_owned(), "2".to_owned()),
            ("CheckoutData.ForterToken".to_owned(), "c51d801bb024c74975c556bad260ee0c_____tt".to_owned()),
            ("CheckoutData.ExternalData.AllowedCharsRegex".to_owned(), "^[A-Za-z0-9,\"\"'`\\s@+&%$#\\*\\(\\)\\[\\]._\\-\\s\\/]*$".to_owned()),
            ("CheckoutData.ExternalData.UnsupportedCharactersErrorTipTimeout".to_owned(), "15000".to_owned()),
            ("CheckoutData.EnableUnsupportedCharactersValidation".to_owned(), "True".to_owned()),
            ("CheckoutData.BillingFirstName".to_owned(), random_name().to_owned()),
            ("CheckoutData.BillingLastName".to_owned(), self.shipping.last.to_owned()),
            ("CheckoutData.Email".to_owned(), self.shipping.email.to_owned()),
            ("CheckoutData.BillingCountryID".to_owned(), id.to_owned()),
            ("CheckoutData.BillingAddress1".to_owned(), self.shipping.address1.to_owned()),
            ("CheckoutData.BillingAddress2".to_owned(), self.shipping.address2.to_owned()),
            ("CheckoutData.BillingCity".to_owned(), self.shipping.city.to_owned()),
            ("CheckoutData.BillingCountyID".to_owned(), "".to_owned()),
            ("CheckoutData.BillingZIP".to_owned(), self.shipping.zip.to_owned()),
            ("CheckoutData.BillingStateID".to_owned(), "".to_owned()),
            ("CheckoutData.BillingPhone".to_owned(), "0157".to_owned() + random_number(8).as_str()),
            ("CheckoutData.OffersFromMerchant".to_owned(), "false".to_owned()),
            ("CheckoutData.ShippingType".to_owned(), "ShippingSameAsBilling".to_owned()),
            ("CheckoutData.ShippingFirstName".to_owned(), "".to_owned()),
            ("CheckoutData.ShippingLastName".to_owned(), "".to_owned()),
            ("CheckoutData.ShippingCountryID".to_owned(), id.to_owned()),
            ("CheckoutData.ShippingAddress1".to_owned(), "".to_owned()),
            ("CheckoutData.ShippingAddress2".to_owned(), "".to_owned()),
            ("CheckoutData.ShippingCity".to_owned(), "".to_owned()),
            ("CheckoutData.ShippingCountyID".to_owned(), "".to_owned()),
            ("CheckoutData.ShippingZIP".to_owned(), "".to_owned()),
            ("CheckoutData.ShippingStateID".to_owned(), "".to_owned()),
            ("CheckoutData.ShippingPhone".to_owned(), "".to_owned()),
            ("CheckoutData.SelectedShippingOptionID".to_owned(), "2728".to_owned()),
            ("CheckoutData.SelectedTaxOption".to_owned(), "3".to_owned()),
            ("ioBlackBox".to_owned(), "0400AKHeWb1CT4UXk1Rjuv1iJgWxIe7xNABi4fWLoKuCjDO1I7X1XkVbR56yHWIulRE2G351wfp+MZWAa+qm7VSS+5sZhQDshHSvL1nHQLC7Q6NpxQ17D0VyS32SOzubmcBOVtcqUhwQqm2+q1yRomE8gXHwdwKa7CjKYUIlPPcC6OJt7RPxTkFnpBRzmAkgrV84+7LojWr7qSyPkgaRwqbiPYz8/yX9hyTE1qBH7CuG2cYMgog1WtCs0txF6Ft8Ea5XAHtlrOPcolu5kVyEZ5Un1A1zM3XireXLCHex/Y96FFfSctR6oPvoc/HBQm916StC15tFwbb5NkVtrNHzt95ePgXs+oQQc60trrGto44dFZ/k0B48ux+V4SxwhwBmNydI6S6la0lc7CNSwejOX5dhKDBAKBIrRLbRVKNQ0ZFLTxhjYx7F0zh29JE75agXQtdj5mNV6Qxgbckwx2tcx5lGRWLQVG6Wu5dBA/z8Qi7K1AaNm0URxgdANnwB7SaH85JmyAfUMZ5LKFqFl81CC76WZKXakoCa7XCW4IhJrOz08Pzf+x4pw6W9esbGEpq7CoRYvFQ3SxgXyvTx/y4JwLTE4Nxv4KowoTq9ZHvy+/OUGsT6WhmJ6uI+FEzI+mPNi65XSGFq0WvDIDWnFBXTiQwVA5aGaf80/3iSPOYe85lXFCR4CjSd/Iq56FG3EcsazU/5If+SxsuUFwyKCknSJs1HLUgFh/6610IF77A/AHthzIQI07p4tGGHBe8mR20/oWqfLQ/2q0KAoA7AqvQx/6iBZvdaZ+9WVICZBBEpCAE1uTZr4w812j83DNT16TeM5diEHf7WUtiodkTjCCAKQfEFHFRukDiyEjh0DtPq/jlyS0vN1/EvKVgGJuwRUFY7xdAiqlr028/3Pq+IowQkPvsTj/7uzpZjbtz5zJmHSgKWcInIB9QxnksoWgDgFtbjU/i//xXxz5iHJ2Gi4wL5ZAuoseRsDHISKM19CXZjVj+Jpy/CWq368mRJ0sgH1DGeSyhahZfNQgu+lmSl2pKAmu1wluCISazs9PD83/seKcOlvXrGxhKauwqEWLxUN0sYF8r08f8uCcC0xODcb+CqMKE6vS38gBEcl7QToyIHPMsM3kmWMCVbBcdxP+fMxqv4WdGaMrHg9Btf04cRCL7CUahnxfbKRXey06cY4a7euFqiMUivlcetxenoDd7uW4fKrTjGYSqEk9HYTgjgklNKGT1jRHxfYV3WeA0UZK211nKAo9ysA7oZ0bYcjqyUg7sNjWRD/a8We6nhBVkkxpY6vV1ZmNpeo69fq+suLh7/vXNSSrf24S1uO6XpqrPHg4UJ3v1LARL9wkv36UGe9i9NNGwrvViY2TSbe4LeLD4mV3TQqnxCcKKadzVSLFfqn4tUjKmd1yI9NBoVDObuig/wKeM18GaUEcqinpriak1JwUr1GLDcE4nZ3meVnnb6dxDiokse7hxjoy31alFmlBHKop6a4mpNScFK9Riw3BOJ2d5nlZ52+ncQ4qJLHu4cY6Mt9WpRZpQRyqKemuLpRlyC5JQyFQQy/m2yvWeDf1XeQj/MbGILRpUuwgADh2omO4uhtmPNcLP6ot9/focoyCt7mIxvb/XIDmwAD2xQYvDt+zi/9Qw5qQdImDqMwPhRLooUyschaDxE3IkI9vGQLxya83WTD/EXGpaCP16YDLvcu2zV4SQS1Ab9fYJjDUkQzBSzKL9vnyo/skDCiZx0aYH7Fo+AaNLvnFBFh8P6er39ssWr96p0mUdkfWuYAO3AOArjJHNP".to_owned()),
            ("CheckoutData.StoreID".to_owned(), "0".to_owned()),
            ("CheckoutData.AddressVerified".to_owned(), "true".to_owned()),
            ("CheckoutData.SelectedPaymentMethodID".to_owned(), "2".to_owned()),
            ("CheckoutData.CurrentPaymentGayewayID".to_owned(), "2".to_owned()),
            ("CheckoutData.MerchantID".to_owned(), "708".to_owned()),
            ("CheckoutData.MultipleAddressesMode".to_owned(), "false".to_owned()),
            ("CheckoutData.MerchantSupportsAddressName".to_owned(), "false".to_owned()),
            ("CheckoutData.MultipleAddressesMode".to_owned(), "true".to_owned()),
            ("CheckoutData.MultipleAddressesMode".to_owned(), "true".to_owned()),
            ("CheckoutData.MultipleAddressesMode".to_owned(), "true".to_owned()),
            ("CheckoutData.MultipleAddressesMode".to_owned(), "true".to_owned()),
            ("CheckoutData.CollectionPointZip".to_owned(), "".to_owned()),
            ("CheckoutData.UseAvalara".to_owned(), "false".to_owned()),
            ("CheckoutData.IsAvalaraLoaded".to_owned(), "false".to_owned()),
            ("CheckoutData.IsUnsupportedRegion".to_owned(), "".to_owned()),
            ("CheckoutData.IsShowTitle".to_owned(), "false".to_owned()),
            ("CheckoutData.IsBillingSavedAddressUsed".to_owned(), "false".to_owned()),
            ("CheckoutData.IsShippingSavedAddressUsed".to_owned(), "false".to_owned()),
            ("CheckoutData.SaveBillingCountryOnChange".to_owned(), "false".to_owned()),
            ("CheckoutData.DisplayInternatioanlPrefixInCheckout".to_owned(), "false".to_owned()),
            ("CheckoutData.IsValidationMessagesV2".to_owned(), "false".to_owned()),
            ("CheckoutData.IgnoreBillingCityRegionValidation".to_owned(), "false".to_owned()),
            ("CheckoutData.IgnoreShippingCityRegionValidation".to_owned(), "false".to_owned()),
            ("CheckoutData.DoLightSave".to_owned(), "false".to_owned()),
        ];

        let req = ReqStruct::new_with_form("https://fs708.global-e.com/checkoutv2/save/8rxx/".to_string() + &s.checkout, body, vec![ReqCookie::new("GlobalE_Data".to_string(), "{\"countryISO\":\"DE\",\"currencyCode\":\"EUR\",\"cultureCode\":\"de\"}".to_string())], self.proxy.clone());
        let res = req.post(Expected::new_status(StatusCode::OK)).await;
        if let Ok(_) = res {
            s.state = 3;
        }
    }

    async fn card(self, s: &mut Session, c: Captcha) {
        let body = vec![
            ("PaymentData.HCaptchaEKey".to_owned(), "".to_owned()),
            ("PaymentData.HCaptchaResponse".to_owned(), c.captcha_two),
            ("PaymentData.cardNum".to_owned(), self.billing.ccnumber.to_owned()),
            ("PaymentData.cardExpiryMonth".to_owned(), self.billing.month.to_owned()),
            ("PaymentData.cardExpiryYear".to_owned(), self.billing.year.to_owned()),
            ("PaymentData.cvdNumber".to_owned(), self.billing.cvv.to_owned()),
            ("PaymentData.checkoutV2".to_owned(), "true".to_owned()),
            ("PaymentData.cartToken".to_owned(), s.checkout.to_owned()),
            ("PaymentData.gatewayId".to_owned(), "2".to_owned()),
            ("PaymentData.paymentMethodId".to_owned(), "2".to_owned()),
            ("PaymentData.machineId".to_owned(), "0400AKHeWb1CT4UXk1Rjuv1iJgWxIe7xNABi4fWLoKuCjDO1I7X1XkVbR56yHWIulRE2G351wfp+MZWAa+qm7VSS+5sZhQDshHSvL1nHQLC7Q6NpxQ17D0VyS32SOzubmcBOVtcqUhwQqm2+q1yRomE8gXHwdwKa7CjKACtSZQ+qebVaf+msKjlXGcSU8+ApUngzZCYnl/wtOFVlx0gLnkNaVYz8/yX9hyTE1qBH7CuG2cYMgog1WtCs0txF6Ft8Ea5XAHtlrOPcolu5kVyEZ5Un1A1zM3XireXLCHex/Y96FFfSctR6oPvoc/HBQm916StC15tFwbb5NkVtrNHzt95ePgXs+oQQc60trrGto44dFZ/k0B48ux+V4SxwhwBmNydI6S6la0lc7CNSwejOX5dhKDBAKBIrRLbRVKNQ0ZFLTxhjYx7F0zh29JE75agXQtdj5mNV6Qxgbckwx2tcx5lGRWLQVG6Wu5dBA/z8Qi7K1AaNm0URxgdANnwB7SaH85JmyAfUMZ5LKFqFl81CC76WZKXakoCa7XCW4IhJrOz08Pzf+x4pw6W9esbGEpq7CoRYvFQ3SxgXyvTx/y4JwLTE4Nxv4KowoTq9ZHvy+/OUGsT6WhmJ6uI+FEzI+mPNi65XSGFq0WvDIDWnFBXTiQwVA5aGaf80/3iSPOYe85lXFCR4CjSd/Iq56FG3EcsazU/5If+SxsuUFwyKCknSJs1HLUgFh/6610IF77A/AHthzIQI07p4tGGHBe8mR20/oWqfLQ/2q0KAoA7AqvQx/6iBZvdaZ+9WVICZeAqVsdwD6Narhe+Eg38xENT16TeM5diEHf7WUtiodkTjCCAKQfEFHFRukDiyEjh0DtPq/jlyS0vN1/EvKVgGJuwRUFY7xdAiqlr028/3Pq+IowQkPvsTj/7uzpZjbtz5zJmHSgKWcInIB9QxnksoWgDgFtbjU/i//xXxz5iHJ2Gi4wL5ZAuoseRsDHISKM19CXZjVj+Jpy/CWq368mRJ0sgH1DGeSyhahZfNQgu+lmSl2pKAmu1wluCISazs9PD83/seKcOlvXrGxhKauwqEWLxUN0sYF8r08f8uCcC0xODcb+CqMKE6vS38gBEcl7QToyIHPMsM3kmWMCVbBcdxP+fMxqv4WdGaMrHg9Btf04cRCL7CUahnxfbKRXey06cY4a7euFqiMUivlcetxenoDd7uW4fKrTjGYSqEk9HYTgjgklNKGT1jRHxfYV3WeA0UZK211nKAo9ysA7oZ0bYcjqyUg7sNjWRD/a8We6nhBVkkxpY6vV1ZmNpeo69fq+suLh7/vXNSSrf24S1uO6XpqrPHg4UJ3v1LARL9wkv36UGe9i9NNGwrvViY2TSbe4LeLD4mV3TQqnxCcKKadzVSLFfqn4tUjKmd1yI9NBoVDObuig/wKeM18GaUEcqinpriak1JwUr1GLDcE4nZ3meVnnb6dxDiokse7hxjoy31alFmlBHKop6a4mpNScFK9Riw3BOJ2d5nlZ52+ncQ4qJLHu4cY6Mt9WpRZpQRyqKemuLpRlyC5JQyFQQy/m2yvWeDf1XeQj/MbGILRpUuwgADh2omO4uhtmPNcLP6ot9/focoyCt7mIxvb/XIDmwAD2xQYvDt+zi/9Qw5qQdImDqMwPhRLooUyschaDxE3IkI9vGQLxya83WTD/EXGpaCP16YDLvcu2zV4SQS1Ab9fYJjDUkQzBSzKL9vnyo/skDCiZx0aYH7Fo+AaNLvnFBFh8P6er39ssWr96p0mUdkfWuYAO3AOArjJHNP".to_owned()),
            ("PaymentData.createTransaction".to_owned(), "true".to_owned()),
            ("PaymentData.checkoutCDNEnabled".to_owned(), "value".to_owned()),
            ("PaymentData.recapchaToken".to_owned(), "".to_owned()),
            ("PaymentData.recapchaTime".to_owned(), "".to_owned()),
            ("PaymentData.customerScreenColorDepth".to_owned(), "24".to_owned()),
            ("PaymentData.customerScreenWidth".to_owned(), "2560".to_owned()),
            ("PaymentData.customerScreenHeight".to_owned(), "1440".to_owned()),
            ("PaymentData.customerTimeZoneOffset".to_owned(), "-120".to_owned()),
            ("PaymentData.customerLanguage".to_owned(), "de-DE".to_owned()),
            ("PaymentData.UrlStructureTokenEncoded".to_owned(), "".to_owned()),
            ("PaymentData.IsValidationMessagesV2".to_owned(), "false".to_owned()),
        ];

        // Normal: https://secure-fs.global-e.com/1/Payments/HandleCreditCardRequestV2/8rxx/
        // Secret: https://securev2.global-e.com/1/Payments/HandleCreditCardRequestV2/8umv/
        let req = ReqStruct::new_with_form("https://securev2.global-e.com/1/Payments/HandleCreditCardRequestV2/8umv/".to_string() + &s.checkout + "?mode=13535", body, vec![ReqCookie::new("GlobalE_Data".to_string(), "{\"countryISO\":\"DE\",\"currencyCode\":\"EUR\",\"cultureCode\":\"de\"}".to_string())], self.proxy.clone());
        let res = req.post(Expected::new_status(StatusCode::FOUND)).await;
        if let Ok(res) = res {
            if res.header.clone().get("Location").is_some() {

                let location = res.header.get("Location").unwrap().clone();
                let token = str::split(str::split(location.to_str().unwrap(), "=").collect::<Vec<&str>>()[1], ".").collect::<Vec<&str>>()[1];

                let bytes = base64::decode(token.clone()).unwrap();
                let str = str::from_utf8(&bytes).unwrap();
                let success_split = str::split(&str, "Success\",\"Value\":\"").collect::<Vec<&str>>()[1];
                let success = str::split(success_split, "\"").collect::<Vec<&str>>()[0];
                println!("{}", success);

                if success.to_lowercase() == "false" {
                    log(str);
                    return;
                }

                s.state = 4;
                s.checkout = location.to_str().unwrap().to_string();
            }
        }
    }
}

#[derive(Deserialize)]
#[allow(dead_code)]
struct CartToken {
    #[serde(rename = "success")]
    pub success: bool,
    #[serde(rename = "CartToken")]
    pub cart_token: String,
}

#[derive(Deserialize)]
#[allow(dead_code)]
struct CartToken2 {
    #[serde(rename = "Success")]
    pub success: bool,
    #[serde(rename = "CartToken")]
    pub cart_token: String,
}