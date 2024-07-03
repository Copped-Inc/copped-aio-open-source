use std::ops::Add;
use async_recursion::async_recursion;
use country_emoji::code;
use hyper::{http, StatusCode};
use serde::{Serialize, Deserialize};
use crate::data::values::{get_rates, get_session, get_store_setting, HIGH_SECURE, PAYMENT};
use crate::modules::Task;
use crate::jig::{random_name, random_number};
use crate::request::cookie::ReqCookie;
use crate::request::ReqStruct;
use crate::modules::session::Session;
use crate::request::expected::Expected;

impl Task {
    #[async_recursion]
    pub(crate) async fn shopify(self) {
        let s = &mut Session::get_from_task(self.clone());

        if s.state == 0 {
            let (checkout, cookies) = self.clone().checkout().await;
            if checkout.is_empty() {
                s.save();
                return;
            }

            s.cookies = cookies;
            s.checkout = checkout;
            s.state = 1;
        }

        if s.state == 1 {
            self.clone().cookies(s).await;
            if s.state != 2 {
                s.save();
                return;
            }
        }

        if s.state == 2 {
            self.clone().submit(s).await;
            if s.state != 3 {
                s.save();
                return;
            }
        }

        self.checked_out(s.clone()).await;
    }

    async fn checkout(self) -> (String, Vec<ReqCookie>) {
        return if get_store_setting(HIGH_SECURE, self.store.clone()) == "0" {
            self.clone().checkout_from_id().await
        } else {
            let (token, cookies) = self.clone().token().await;
            if token.is_empty() {
                return (String::new(), vec![]);
            }

            let cookies = self.clone().add_to_cart(token, cookies).await;
            return if cookies.is_empty() {
                (String::new(), vec![])
            } else {
                self.clone().checkout_from_cart(cookies).await
            }
        }
    }

    pub(crate) async fn checkout_from_id(self) -> (String, Vec<ReqCookie>) {
        let req = ReqStruct::new_none_body("https://".to_owned().add(self.store.as_str()).add("/cart/").add(self.prod_id.as_str()).add(":1").to_string(), vec![], self.proxy.clone());
        let res = req.get(false, Expected::new_status(http::StatusCode::FOUND)).await;
        if let Ok(res) = res {
            if res.header.clone().get("Location").is_some() {
                return (res.header.get("Location").unwrap().to_str().unwrap().to_string(), res.cookies);
            }
        }
        (String::new(), vec![])
    }

    async fn token(self) -> (String, Vec<ReqCookie>) {
        let req = ReqStruct::new_none_body("https://".to_owned().add(self.store.as_str()).add("/variants/").add(self.prod_id.as_str()), vec![], self.proxy.clone());
        let res = req.get(true, Expected::new_body_contains("data-token=\"")).await;
        if let Ok(res) = res {
            if res.body.clone().contains("data-token=\"") {
                let token = res.body.split("data-token=\"").nth(1).unwrap().split("\"").next().unwrap().to_string();
                return (token, res.cookies);
            }
        }
        (String::new(), vec![])
    }

    pub(crate) async fn add_to_cart(self, t: String, c: Vec<ReqCookie>) -> Vec<ReqCookie> {
        let body = vec![
            ("properties[_token]".to_owned(), t.to_owned()),
            ("id".to_owned(), self.prod_id.to_owned()),
            ("quantity".to_owned(), "1".to_owned()),
        ];

        let req = ReqStruct::new_with_form("https://".to_owned().add(self.store.as_str()).add("/cart/add.js"), body, c, self.proxy.clone());
        let res = req.post(Expected::new_status(StatusCode::OK)).await;
        if let Ok(res) = res {
            return res.cookies;
        }
        vec![]
    }

    pub(crate) async fn checkout_from_cart(self, c: Vec<ReqCookie>) -> (String, Vec<ReqCookie>) {
        let req = ReqStruct::new_none_body("https://".to_owned().add(self.store.as_str()).add("/checkout").to_string(), c, self.proxy.clone());
        let res = req.get(false, Expected::new_status(http::StatusCode::FOUND)).await;
        if let Ok(res) = res {
            if res.header.clone().get("Location").is_some() {
                return (res.header.get("Location").unwrap().to_str().unwrap().to_string(), res.cookies);
            }
        }
        (String::new(), vec![])
    }

    pub(crate) async fn cookies(self, s: &mut Session) {
        let req = ReqStruct::new_none_body(s.checkout.clone(), s.cookies.clone(), self.proxy.clone());
        let res = req.get(false, Expected::new_body_contains("Shopify.Checkout.estimatedPrice = ")).await;
        if let Ok(res) = res {
            s.cookies = res.cookies;

            if res.body.clone().contains("Shopify.Checkout.estimatedPrice = ") {
                let price = res.body.split("Shopify.Checkout.estimatedPrice = ").nth(1).unwrap().split(";").next().unwrap().to_string();

                s.price = price;
                s.state = 2;
            }
        }
    }

    pub(crate) async fn submit(self, s: &mut Session) {
        let session = get_session();
        let r = get_rates(self.clone());
        if session.is_empty() || r.is_empty() {
            return;
        }

        let price = s.price.parse::<f32>().unwrap();
        let rate_split = r.split("-").collect::<Vec<&str>>();
        let rate = rate_split[rate_split.len() -1].parse::<f32>().unwrap();

        s.price = (price + rate).to_string();
        if !s.price.contains(".") {
            s.price = s.price.clone().add(".00");
        }

        let body = vec![
            ("_method".to_owned(), "patch".to_owned()),
            ("authenticity_token".to_owned(), "ID3RQe0VmzvcWcCWh4C3d67nUKqzmuIn6ij018cMTBPUiV3g9RbTBlAHpHAokd9HGD8mkKFhEG6uTEUqIIzWKQ".to_owned()),
            ("previous_step".to_owned(), "contact_information".to_owned()),
            ("step".to_owned(), "shipping_method".to_owned()),
            ("checkout[email]".to_owned(), self.shipping.email.to_owned()),
            ("checkout[shipping_address][first_name]".to_owned(), random_name().to_owned()),
            ("checkout[shipping_address][last_name]".to_owned(), self.shipping.last.to_owned()),
            ("checkout[shipping_address][company]".to_owned(), "".to_owned()),
            ("checkout[shipping_address][address1]".to_owned(), self.shipping.address1.to_owned()),
            ("checkout[shipping_address][address2]".to_owned(), self.shipping.address2.to_owned()),
            ("checkout[shipping_address][city]".to_owned(), self.shipping.city.to_owned()),
            ("checkout[shipping_address][country]".to_owned(), code(self.shipping.country.clone().as_str()).unwrap().to_owned()),
            ("checkout[shipping_address][province]".to_owned(), "".to_owned()),
            ("checkout[shipping_address][zip]".to_owned(), self.shipping.zip.to_owned()),
            ("checkout[shipping_address][phone]".to_owned(), "0157".to_owned() + random_number(8).as_str()),
            ("checkout[shipping_address][country]".to_owned(), self.shipping.country.to_owned()),
            ("checkout[shipping_address][first_name]".to_owned(), random_name().to_owned()),
            ("checkout[shipping_address][last_name]".to_owned(), self.shipping.last.to_owned()),
            ("checkout[shipping_address][company]".to_owned(), "".to_owned()),
            ("checkout[shipping_address][address1]".to_owned(), self.shipping.address1.to_owned()),
            ("checkout[shipping_address][address2]".to_owned(), self.shipping.address2.to_owned()),
            ("checkout[shipping_address][zip]".to_owned(), self.shipping.zip.to_owned()),
            ("checkout[shipping_address][city]".to_owned(), self.shipping.city.to_owned()),
            ("checkout[remember_me]".to_owned(), "".to_owned()),
            ("checkout[remember_me]".to_owned(), "0".to_owned()),
            ("checkout[client_details][browser_width]".to_owned(), "2543".to_owned()),
            ("checkout[client_details][browser_height]".to_owned(), "1289".to_owned()),
            ("checkout[client_details][javascript_enabled]".to_owned(), "1".to_owned()),
            ("checkout[client_details][color_depth]".to_owned(), "24".to_owned()),
            ("checkout[client_details][java_enabled]".to_owned(), "false".to_owned()),
            ("checkout[client_details][browser_tz]".to_owned(), "-120".to_owned()),
            ("checkout[shipping_rate][id]".to_owned(), r.to_owned()),
            ("s".to_owned(), session.to_owned()),
            ("checkout[payment_gateway]".to_owned(), get_store_setting(PAYMENT, self.store).to_owned()),
            ("checkout[credit_card][vault]".to_owned(), "false".to_owned()),
            ("checkout[different_billing_address]".to_owned(), "false".to_owned()),
            ("checkout[total_price]".to_owned(), s.price.replace(".", "").to_owned()),
            ("complete".to_owned(), "1".to_owned()),
        ];

        let req = ReqStruct::new_with_form(s.checkout.clone(), body, s.cookies.clone(), self.proxy.clone());
        let res = req.post(Expected::new_status(StatusCode::FOUND)).await;
        if let Ok(res) = res {
            s.cookies = res.cookies;

            if res.header.clone().get("Location").is_some() &&
                res.header.get("Location").unwrap().to_str().unwrap().contains("processing") {

                s.state = 3;
                s.checkout = res.header.get("Location").unwrap().to_str().unwrap().to_string();
            }
        }
    }
}

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct SessionReq {
    #[serde(rename = "credit_card")]
    pub credit_card: CreditCard,
}

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct CreditCard {
    pub number: String,
    pub name: String,
    #[serde(rename = "start_month")]
    pub start_month: i64,
    #[serde(rename = "start_year")]
    pub start_year: i64,
    pub month: i64,
    pub year: i64,
    #[serde(rename = "verification_value")]
    pub verification_value: String,
    #[serde(rename = "issue_number")]
    pub issue_number: String,
}
