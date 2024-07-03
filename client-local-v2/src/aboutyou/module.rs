use std::ops::Add;
use reqwest::StatusCode;
use crate::{error, jig, log, payments, request, session};
use crate::request::expected;

impl session::Session {
    pub async fn preload_aboutyou(&mut self) {
        self.create_account().await;
        if self.state == 0 {
            return;
        }

        self.cart_aboutyou().await;
        if self.state == 1 {
            return;
        }

        self.update().await;
        if self.state == 2 {
            return;
        }

        self.state().await;
        if self.state == 3 {
            return;
        }

        self.shipping().await;
        if self.state == 4 {
            return;
        }

        self.payment().await;
        if self.state == 5 {
            return;
        }

        self.execute().await;
        if self.state == 6 {
            return;
        }

        let r = self.payments("aboutyou".to_string()).await;
        if let Err(e) = r {
            error!("{}", e);
            return;
        }

        let r = payments::check(r.unwrap());
        if let Err(e) = r {
            error!("{}", e);
            return;
        }

        self.state = 0;
    }

    pub async fn start_aboutyou(&mut self) {
        self.cart_aboutyou().await;
        if self.state == 0 {
            return;
        }

        self.execute().await;
        if self.state == 1 {
            return;
        }

        self.state = 0;
    }

    async fn create_account(&mut self) {
        let password = jig::password();
        self.add_data.insert("password", password.clone());
        log!("{} {}", self.clone().task.shipping.email, &password);
        let req = request::ReqStruct::new(self.proxy())
            .url("https://grips-web.aboutyou.com/checkout.CheckoutV1/registerWithEmail".to_string())
            .body_grcp(crate::aboutyou::register_body(self.clone().task.shipping.email, password, jig::random_name(), jig::name(self.clone().task.shipping.last)))
            .tls();

        let res = req.post(expected::Expected::new_status(StatusCode::OK)).await;
        if let Ok(res) = res {
            if res.status_code == StatusCode::OK {
                self.add_state();
                self.cookies = res.cookies;

                let mut session = "eyJ".to_string().add(res.body.split("eyJ").collect::<Vec<&str>>()[1].to_string().split("\u{12}").collect::<Vec<&str>>()[0].to_string().as_str());
                log!("{}: {}", session.len(), session);

                if session.len() == 328 {
                    session = session.add("%3D%3D");
                } else if session.len() == 330 {
                    session = session.add("0%3D");
                } else if session.len() == 332 {
                    session = session.add("0=");
                } else if session.len() == 336 {
                    session = session.replace("%3D%3D", "0%3D");
                } else if session.len() == 338 {
                    session = session.replace("0%3D", "");
                }

                log!("{}: {}", session.len(), session);

                self.add_data.insert("session", session.clone());
                self.add_data.insert("auth", res.body.split("$").collect::<Vec<&str>>()[1].to_string().split("\u{1a}").collect::<Vec<&str>>()[0].to_string());
                self.add_data.insert("secret", res.body.split("\u{12}d").collect::<Vec<&str>>()[1].to_string().split("\u{12}").collect::<Vec<&str>>()[0].to_string());

                let cookie = request::cookie::ReqCookie::new("checkout_sid".to_string(), session);
                self.cookies.append(&mut vec![cookie]);
            }
        }
    }

    async fn cart_aboutyou(&mut self) {
        let body = "{\"variantId\":".to_string()
            .add(self.task.size.as_str())
            .add(",\"quantity\":1,\"customData\":{\"variantId\":")
            .add(self.task.size.as_str())
            .add(",\"productId\":")
            .add(self.task.product_id.as_str())
            .add(",\"sponsoredType\":null,\"sortingChannel\":null,\"linkedContentType\":null,\"linkedContentId\":null,\"originDevice\":\"desktop\"},\"shopId\":688}");

        let req = request::ReqStruct::new(self.proxy())
            .url("https://api.aboutyou.com/user/me/basket/bapi".to_string())
            .cookies(self.clone().cookies)
            .add_header("x-auth-token".to_string(), self.clone().add_data.get("auth").unwrap().to_string())
            .body_json(body)
            .tls();

        let res = req.post(expected::Expected::new_body_contains("\"items\":[{\"key\":\"")).await;
        if let Ok(res) = res {
            if res.status_code == StatusCode::CREATED && res.body.contains("\"items\":[{\"key\":\"") {
                self.add_state();
                self.add_data.insert("basket", res.body.split("\"items\":[{\"key\":\"").collect::<Vec<&str>>()[1].to_string().split("\"").collect::<Vec<&str>>()[0].to_string());
                self.cookies = res.cookies;
            }
        }
    }

    async fn update(&mut self) {
        let req = request::ReqStruct::new(self.proxy())
            .url("https://grips-web.aboutyou.com/checkout.CheckoutV1/getCheckoutUrl".to_string())
            .cookies(self.clone().cookies)
            .body_grcp(crate::aboutyou::checkout_body(self.add_data.get("session").unwrap().to_string(), self.add_data.get("secret").unwrap().to_string(), self.add_data.get("auth").unwrap().to_string()))
            .tls();

        let res = req.post(expected::Expected::new_body_contains("customer")).await;
        if let Ok(res) = res {
            if res.status_code == StatusCode::OK && res.body.contains("customer") {
                self.add_data.insert("id", res.body.split("customer_").collect::<Vec<&str>>()[1].to_string().split("&").collect::<Vec<&str>>()[0].to_string());
                self.add_state();
                self.cookies = res.cookies;
            }
        }
    }

    async fn shipping(&mut self) {
        let (street, number) = self.split_address();
        let body = "{\"country\":{\"iso2Code\":\"".to_string()
            .add(country_emoji::code(self.task.shipping.country.clone().as_str()).unwrap())
            .add("\"},\"gender\":\"f\",\"firstName\":\"")
            .add(jig::random_name().as_str())
            .add("\",\"lastName\":\"")
            .add(jig::name(self.clone().task.shipping.last).as_str())
            .add("\",\"street\":\"")
            .add(street.as_str())
            .add("\",\"houseNumber\":\"")
            .add(number.as_str())
            .add("\",\"zipCode\":\"")
            .add(self.clone().task.shipping.zip.as_str())
            .add("\",\"city\":\"")
            .add(self.clone().task.shipping.city.as_str())
            .add("\"}");

        let req = request::ReqStruct::new(self.proxy())
            .url("https://checkout-v3.aboutyou.de/api/co/v3/state/order/addresses/shipping".to_string())
            .cookies(self.clone().cookies)
            .add_header("x-shop-id".to_string(), "688".to_string())
            .add_header("x-signature".to_string(), self.get_secret(body.clone()))
            .body_json(body)
            .tls();

        let res = req.put(expected::Expected::new_status(StatusCode::CREATED)).await;
        if let Ok(res) = res {
            if res.status_code == StatusCode::CREATED {
                self.add_state();
                self.cookies = res.cookies;
            }
        }
    }

    async fn state(&mut self) {
        let body = "{\"device\":\"desktop\",\"basketId\":\"aboutyou_customer_".to_string()
            .add(self.add_data.get("id").unwrap().as_str())
            .add("\",\"campaignKey\":\"00\",\"customData\":{\"device\":\"desktop-web\"}}");

        let req = request::ReqStruct::new(self.proxy())
            .url("https://checkout-v3.aboutyou.de/api/co/v3/state".to_string())
            .cookies(self.cookies.clone())
            .add_header("x-session".to_string(), self.get_session())
            .add_header("x-shop-id".to_string(), "688".to_string())
            .add_header("x-signature".to_string(), self.get_secret(body.clone()))
            .body_json(body)
            .tls();

        let res = req.put(expected::Expected::new_status(StatusCode::OK)).await;
        if let Ok(res) = res {
            if res.status_code == StatusCode::OK {
                self.add_state();
                self.cookies = res.cookies;
            }
        }
    }

    async fn payment(&mut self) {
        let body = "{\"data\":{\"paymentOptionKey\":\"paypal_instant\"}}".to_string();

        let req = request::ReqStruct::new(self.proxy())
            .url("https://checkout-v3.aboutyou.de/api/co/v3/state/order/payment/option".to_string())
            .cookies(self.cookies.clone())
            .add_header("x-signature".to_string(), self.get_secret(body.clone()))
            .body_json(body)
            .tls();

        let res = req.post(expected::Expected::new_status(StatusCode::CREATED)).await;
        if let Ok(res) = res {
            if res.status_code == StatusCode::CREATED {
                self.add_state();
                self.cookies = res.cookies;
            }
        }
    }

    async fn execute(&mut self) {
        let body = "{\"data\":{\"paymentOptionKey\":\"paypal_instant\"},\"referer\":\"https://checkout-v3.aboutyou.de/api/co/v3/state/order/addresses/shipping\"}".to_string();

        let req = request::ReqStruct::new(self.proxy())
            .url("https://checkout-v3.aboutyou.de/api/co/v3/state/order/confirmation/execute".to_string())
            .cookies(self.cookies.clone())
            .add_header("x-signature".to_string(), self.get_secret(body.to_string()))
            .body_json(body)
            .tls();

        let res = req.post(expected::Expected::new_body_contains("approve?ba_token=")).await;
        if let Ok(res) = res {
            if res.status_code == StatusCode::OK {
                self.add_data.insert("paypal", res.body.split("approve?ba_token=").collect::<Vec<&str>>()[1].to_string().split("&").collect::<Vec<&str>>()[0].to_string());
                self.add_state();
                self.cookies = res.cookies;
            }
        }
    }

    async fn uncart(&mut self) {
        let req = request::ReqStruct::new(self.proxy())
            .url("https://tadarida-web.aboutyou.com/aysa_api.services.basket.v1.page.BasketPageService/RemoveItem".to_string())
            .cookies(self.clone().cookies)
            .body_grcp(crate::aboutyou::uncart_body(self.add_data.get("auth").unwrap().to_string(), self.add_data.get("basket").unwrap().to_string()))
            .tls();

        let res = req.post(expected::Expected::new_status(StatusCode::OK)).await;
        if let Ok(res) = res {
            if res.status_code == StatusCode::OK {
                self.add_data.remove("basket");
                self.add_state();
                self.cookies = res.cookies;
            }
        }
    }

    async fn orders(&mut self) {
        let req = request::ReqStruct::new(self.proxy())
            .url("https://grips-web.aboutyou.com/orders.v1.OrdersService/getOrders".to_string())
            .cookies(self.clone().cookies)
            .body_grcp(crate::aboutyou::orders_body(self.add_data.get("session").unwrap().to_string(), self.add_data.get("secret").unwrap().to_string(), self.add_data.get("auth").unwrap().to_string()))
            .tls();

        let res = req.post(expected::Expected::new_status(StatusCode::OK)).await;
        if let Ok(res) = res {
            if res.status_code == StatusCode::OK {
                for header in res.header.iter() {
                    println!("{:?}: {:?}", header.0, header.1);
                }
                println!("{:?}", res.body);
                println!("{:?}", res.body.bytes());

                self.add_data.remove("basket");
                self.add_state();
                self.cookies = res.cookies;
            }
        }
    }

    async fn cancel(&mut self) {
        let req = request::ReqStruct::new(self.proxy())
            .url("https://grips-web.aboutyou.com/orders.v1.OrdersService/cancelOrder".to_string())
            .cookies(self.clone().cookies)
            .body_grcp(crate::aboutyou::cancel_body(self.add_data.get("session").unwrap().to_string(), self.add_data.get("secret").unwrap().to_string(), self.add_data.get("auth").unwrap().to_string()))
            .tls();

        let res = req.post(expected::Expected::new_status(StatusCode::OK)).await;
        if let Ok(res) = res {
            if res.status_code == StatusCode::OK {
                self.add_data.remove("basket");
                self.add_state();
                self.cookies = res.cookies;
            }
        }
    }
}