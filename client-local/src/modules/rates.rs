use std::ops::Add;
use country_emoji::code;
use hyper::{http, StatusCode};
use crate::data::values::{get_store_setting, HIGH_SECURE, Shipping};
use serde::{Deserialize, Serialize};
use crate::jig::{random_name, random_number, xyz};
use crate::ReqProxy;
use crate::request::cookie::ReqCookie;
use crate::request::ReqStruct;
use crate::request::expected::Expected;

impl Shipping {
    pub(crate) async fn get_rate(self, s: String) -> String {
        let (sku, handle) = self.clone().sku(s.clone()).await;
        if sku.is_empty() {
            return String::new();
        }

        let (checkout, cookies) = self.clone().checkout(s.clone(), sku, handle).await;
        if checkout.is_empty() {
            return String::new();
        }

        let (cookies, shipping) = self.clone().address(checkout.clone(), cookies).await;
        if shipping.is_empty() {
            return String::new();
        }

        let rate = self.clone().rate(checkout.clone(), cookies).await;
        if rate.is_empty() {
            return String::new();
        }

        rate
    }

    async fn sku(self, s: String) -> (String, String) {
        let req = ReqStruct::new_none_body(s.add("products.json"), vec![], ReqProxy::get());
        let res = req.get(false, None).await;
        if let Ok(res) = res {
            let body = serde_json::from_str::<Shopify>(&res.body).unwrap();
            for pr in body.products {
                if pr.tags.contains(&"CustomStockStatus:Online Release".to_string()) || pr.tags.contains(&"CustomStockStatus:Release via App".to_string())  || pr.tags.contains(&"CustomStockStatus:Sold Out".to_string()) {
                    continue;
                }

                for vr in pr.variants {
                    let price = vr.price.parse::<f32>().unwrap();
                    if vr.available && price > 100.0 {
                        return (vr.id.clone().to_string(), pr.handle.clone().to_string());
                    }
                }
            }
        }
        (String::new(), String::new())
    }

    async fn checkout(self, s: String, p: String, h: String) -> (String, Vec<ReqCookie>) {
        let (checkout, cookies) = if get_store_setting(HIGH_SECURE, s.clone()) == "0" {
            self.clone().checkout_from_id(s.clone(), p).await
        } else {
            let (token, cookies) = self.clone().token(s.clone(), h.clone()).await;
            if token.is_empty() {
                return (String::new(), vec![]);
            }

            let cookies = self.clone().add_to_cart(s.clone(), p, token, cookies).await;
            return if cookies.is_empty() {
                (String::new(), vec![])
            } else {
                self.clone().checkout_from_cart(s.clone(), cookies).await
            }
        };

        (checkout, cookies)
/*        return if checkout.contains("queue") {
            self.clone().queue(checkout, cookies).await
        } else {
            (checkout, cookies)
        }*/
    }

    async fn checkout_from_id(self, s: String, p: String) -> (String, Vec<ReqCookie>) {
        let req = ReqStruct::new_none_body(s.add("cart/").add(p.as_str()).add(":1").to_string(), vec![], ReqProxy::get());
        let res = req.get(false, Expected::new_status(http::StatusCode::FOUND)).await;
        if let Ok(res) = res {
            if res.header.get("Location").is_some() {
                return (res.header.get("Location").unwrap().to_str().unwrap().to_string(), res.cookies);
            }
        }
        (String::new(), vec![])
    }

    async fn token(self, s: String, h: String) -> (String, Vec<ReqCookie>) {
        let req = ReqStruct::new_none_body(s.clone().add("products/").add(h.as_str()).add("?variant=").add(s.as_str()), vec![], ReqProxy::get());
        let res = req.get(false, Expected::new_body_contains("data-token=\"")).await;
        if let Ok(res) = res {
            if res.body.clone().contains("data-token=\"") {
                let token = res.body.split("data-token=\"").nth(1).unwrap().split("\"").next().unwrap().to_string();
                return (token, res.cookies);
            }
        }
        (String::new(), vec![])
    }

    pub(crate) async fn add_to_cart(self, s: String, p: String, t: String, c: Vec<ReqCookie>) -> Vec<ReqCookie> {
        let body = vec![
            ("properties[_token]".to_owned(), t.to_owned()),
            ("id".to_owned(), p.to_owned()),
            ("quantity".to_owned(), "1".to_owned()),
        ];

        let req = ReqStruct::new_with_form(s.add("cart/add.js"), body, c, ReqProxy::get());
        let res = req.post(Expected::new_status(StatusCode::OK)).await;
        if let Ok(res) = res {
            return res.cookies;
        }
        vec![]
    }

    pub(crate) async fn checkout_from_cart(self, s: String, c: Vec<ReqCookie>) -> (String, Vec<ReqCookie>) {
        let req = ReqStruct::new_none_body(s.add("checkout"), c, ReqProxy::get());
        let res = req.get(false, Expected::new_status(StatusCode::FOUND)).await;
        if let Ok(res) = res {
            if res.header.clone().get("Location").is_some() {
                return (res.header.get("Location").unwrap().to_str().unwrap().to_string(), res.cookies);
            }
        }
        (String::new(), vec![])
    }
/*
    #[async_recursion]
    pub(crate) async fn queue(self, s: String, c: Vec<ReqCookie>) -> (String, Vec<ReqCookie>) {
        let req = ReqStruct::new_none_body(s.clone(), c, ReqProxy::from_file());
        let res = req.get(false).await;
        if let Ok(res) = res {
            return if res.header.clone().get("Location").is_some() {
                (res.header.get("Location").unwrap().to_str().unwrap().to_string(), res.cookies)
            } else {
                sleep(Duration::from_secs(5));
                self.clone().queue(s.clone(), res.cookies).await
            }
        }
        (String::new(), vec![])
    }*/

    async fn address(self, l: String, c: Vec<ReqCookie>) -> (Vec<ReqCookie>, String) {
        let mut email = self.email.clone();
        if !email.contains("@") {
            email = email.clone().add(xyz().as_str()).add("@").add(email.as_str());
        }
        let body =  vec![
            ("_method".to_owned(), "patch".to_owned()),
            ("authenticity_token".to_owned(), "ID3RQe0VmzvcWcCWh4C3d67nUKqzmuIn6ij018cMTBPUiV3g9RbTBlAHpHAokd9HGD8mkKFhEG6uTEUqIIzWKQ".to_owned()),
            ("previous_step".to_owned(), "contact_information".to_owned()),
            ("step".to_owned(), "shipping_method".to_owned()),
            ("checkout[email]".to_owned(), email),
            ("checkout[shipping_address][first_name]".to_owned(), random_name().to_owned()),
            ("checkout[shipping_address][last_name]".to_owned(), self.last.to_owned()),
            ("checkout[shipping_address][company]".to_owned(), "".to_owned()),
            ("checkout[shipping_address][address1]".to_owned(), self.address1.to_owned()),
            ("checkout[shipping_address][address2]".to_owned(), self.address2.to_owned()),
            ("checkout[shipping_address][city]".to_owned(), self.city.to_owned()),
            ("checkout[shipping_address][country]".to_owned(), code(self.country.clone().as_str()).unwrap().to_owned()),
            ("checkout[shipping_address][province]".to_owned(), "".to_owned()),
            ("checkout[shipping_address][zip]".to_owned(), self.zip.to_owned()),
            ("checkout[shipping_address][phone]".to_owned(), "0157".to_owned() + random_number(8).as_str()),
            ("checkout[shipping_address][country]".to_owned(), self.country.to_owned()),
            ("checkout[shipping_address][first_name]".to_owned(), random_name().to_owned()),
            ("checkout[shipping_address][last_name]".to_owned(), self.last.to_owned()),
            ("checkout[shipping_address][company]".to_owned(), "".to_owned()),
            ("checkout[shipping_address][address1]".to_owned(), self.address1.to_owned()),
            ("checkout[shipping_address][address2]".to_owned(), self.address2.to_owned()),
            ("checkout[shipping_address][zip]".to_owned(), self.zip.to_owned()),
            ("checkout[shipping_address][city]".to_owned(), self.city.to_owned()),
            ("checkout[remember_me]".to_owned(), "".to_owned()),
            ("checkout[remember_me]".to_owned(), "0".to_owned()),
            ("checkout[client_details][browser_width]".to_owned(), "2543".to_owned()),
            ("checkout[client_details][browser_height]".to_owned(), "1289".to_owned()),
            ("checkout[client_details][javascript_enabled]".to_owned(), "1".to_owned()),
            ("checkout[client_details][color_depth]".to_owned(), "24".to_owned()),
            ("checkout[client_details][java_enabled]".to_owned(), "false".to_owned()),
            ("checkout[client_details][browser_tz]".to_owned(), "-120".to_owned()),
        ];

        let req = ReqStruct::new_with_form(l, body, c, ReqProxy::get());
        let res = req.post(Expected::new_status(StatusCode::FOUND)).await;
        if let Ok(res) = res {
            if res.header.get("Location").is_some() {
                return (res.cookies, res.header.get("Location").unwrap().to_str().unwrap().to_string());
            }
        }
        (vec![], String::new())
    }

    async fn rate(self, l: String, c: Vec<ReqCookie>) -> String {
        let req = ReqStruct::new_none_body(l, c, ReqProxy::get());
        let res = req.get(false, Expected::new_body_contains("data-shipping-method=\"")).await;
        if let Ok(res) = res {
            if res.body.clone().contains("data-shipping-method=\"") {
                let rate = res.body.split("data-shipping-method=\"").nth(1).unwrap().split("\"").nth(0).unwrap().to_string();
                if rate.to_lowercase().contains("express") {
                    panic!("express shipping");
                }
                return res.body.split("data-shipping-method=\"").nth(1).unwrap().split("\"").nth(0).unwrap().to_string();
            }
        }
        String::new()
    }
}

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Shopify {
    pub products: Vec<Product>,
}

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Product {
    pub handle: String,
    pub variants: Vec<Variant>,
    pub tags: Vec<String>,
}

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Variant {
    pub id: i64,
    pub available: bool,
    pub price: String,
}