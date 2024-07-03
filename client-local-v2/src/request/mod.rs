use reqwest::{Client, RequestBuilder, StatusCode};
use reqwest::redirect::Policy;
use crate::{api, log, session};
use async_recursion::async_recursion;
use crate::request::cookie::ReqCookie;
use crate::request::expected::Expected;
use crate::request::proxy::ReqProxy;
use crate::request::response::ResStruct;

pub mod cookie;
pub mod proxy;
pub mod expected;
mod response;

pub enum ContentType {
    FormUrlEncoded,
    JSON,
    GRCP,
    None,
}

pub struct ReqStruct {
    url: String,
    form: Vec<(String, String)>,
    body: Vec<u8>,
    cookies: Vec<ReqCookie>,
    content_type: ContentType,
    proxy: ReqProxy,
    retry: i32,
    tls: bool,
    custom_header: Vec<(String, String)>,
    status_whitelist: Vec<StatusCode>,
}

#[allow(dead_code)]
impl ReqStruct {
    pub fn new(proxy: ReqProxy) -> Self {
        ReqStruct {
            url: String::new(),
            form: Vec::new(),
            body: Vec::new(),
            cookies: Vec::new(),
            content_type: ContentType::None,
            proxy,
            retry: 0,
            tls: false,
            custom_header: vec![],
            status_whitelist: vec![],
        }
    }

    pub fn new_from_session(session: &session::Session) -> Self {
        let mut r = ReqStruct::new(session.proxy());
        r.cookies = session.cookies.clone();
        r
    }

    pub fn url(mut self, url: String) -> Self {
        self.url = url;
        self
    }

    pub fn form(mut self, form: Vec<(String, String)>) -> Self {
        self.form = form;
        self.content_type = ContentType::FormUrlEncoded;
        self
    }

    pub fn body(mut self, body: Vec<u8>, content_type: ContentType) -> Self {
        self.body = body;
        self.content_type = content_type;
        self
    }

    pub fn body_json(self, body: String) -> Self {
        self.body(body.into_bytes(), ContentType::JSON)
    }

    pub fn body_grcp(self, body: Vec<u8>) -> Self {
        self.body(body, ContentType::GRCP)
    }

    pub fn cookies(mut self, cookies: Vec<ReqCookie>) -> Self {
        self.cookies = cookies;
        self
    }

    pub fn tls(mut self) -> Self {
        self.tls = true;
        self
    }

    pub fn add_header(mut self, key: String, value: String) -> Self {
        self.custom_header.push((key, value));
        self
    }

    pub fn add_status_whitelist(mut self, status: StatusCode) -> Self {
        self.status_whitelist.push(status);
        self
    }

    pub async fn get(self, r: bool, expected: Option<Expected>) -> Result<ResStruct, &'static str> {
        self.request(expected, r, api::METHOD_GET).await
    }

    pub async fn post(self, expected: Option<Expected>) -> Result<ResStruct, &'static str> {
        self.request(expected, false, api::METHOD_POST).await
    }

    pub async fn put(self, expected: Option<Expected>) -> Result<ResStruct, &'static str> {
        self.request(expected, false, api::METHOD_PUT).await
    }

    pub async fn delete(self, expected: Option<Expected>) -> Result<ResStruct, &'static str> {
        self.request(expected, false, api::METHOD_DELETE).await
    }

    #[async_recursion]
    async fn request(mut self, expected: Option<Expected>, r: bool, method: i8) -> Result<ResStruct, &'static str> {
        let mut custom = custom_redirect();
        if r {
            custom = Policy::default();
        }

        let client_builder = Client::builder()
            .proxy(self.proxy.proxy).brotli(true).gzip(true)
            .redirect(custom);

        let client = if self.tls {
            client_builder.use_rustls_tls().build().unwrap()
        } else {
            client_builder.build().unwrap()
        };

        let mut builder: RequestBuilder;
        match method {
            api::METHOD_GET => {
                builder = client.get(self.url.as_str());
            }
            api::METHOD_POST => {
                builder = client.post(self.url.as_str());
            }
            api::METHOD_PUT => {
                builder = client.put(self.url.as_str());
            }
            api::METHOD_DELETE => {
                builder = client.delete(self.url.as_str());
            }
            _ => {
                return Err("Invalid method");
            }
        }

        builder = builder.header("Host", (self.url.clone().split("://").nth(1).unwrap()).split("/").nth(0).unwrap())
            .header("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36")
            .header("Accept", "*/*")
            .header("Accept-Encoding", "gzip, deflate, br")
            .header("sec-ch-ua", "\"Google Chrome\";v=\"107\", \"Chromium\";v=\"107\", \";Not A Brand\";v=\"24\"")
            .header("sec-ch-ua-mobile", "?0")
            .header("sec-ch-ua-platform", "\"Windows\"")
            .header("sec-fetch-dest", "empty")
            .header("sec-fetch-mode", "cors")
            .header("sec-fetch-site", "same-origin");

        for (key, value) in &self.custom_header {
            builder = builder.header(key, value);
        }

        if !self.cookies.is_empty() {
            let mut cookie_str = "".to_string();
            for cookie in &self.cookies {
                cookie_str.push_str(cookie.to_string().as_str());
            }
            builder = builder.header("Cookie", cookie_str);
        }

        if method == api::METHOD_POST || method == api::METHOD_PUT {
            match self.content_type {
                ContentType::FormUrlEncoded => {
                    builder = builder.header("Content-Type", "application/x-www-form-urlencoded")
                        .form(&self.form);
                },
                ContentType::JSON => {
                    builder = builder.header("Content-Type", "application/json")
                        .body(self.body.clone());
                }
                ContentType::GRCP => {
                    // println!("{:?}", self.body);
                    builder = builder.header("Content-Type", "application/grpc-web+proto")
                        .header("x-grpc-web", "1")
                        .body(self.body.clone());
                }
                _ => {}
            }
        }

        let res = builder.send().await;
        if res.is_err() {
            println!("{}", res.as_ref().unwrap_err());
            if self.retry < 3 {
                self.retry += 1;
                spin_sleep::sleep(std::time::Duration::from_secs(3));
                self.proxy = ReqProxy::get();
                return self.request(expected, r, method).await;
            }
            return Err("Request failed");
        }

        let status = res.as_ref().unwrap().status();
        let cookies = ReqCookie::from_response(res.as_ref().unwrap(), self.cookies.clone());
        let header = res.as_ref().unwrap().headers().clone();
        let body = res.unwrap().text().await.unwrap();

        expected::print_status(status, self.url.clone(), body.clone(), header.clone(), expected.clone());
        log!("Cookies got: {} -> Cookies returning: {}", self.cookies.len(), cookies.len());

        match status {
            StatusCode::OK => Ok(ResStruct::new(status, body, cookies, header)),
            StatusCode::CREATED => Ok(ResStruct::new(status, body, cookies, header)),
            StatusCode::FOUND => Ok(ResStruct::new(status, "".to_string(), cookies, header)),
            _ => {
                if self.status_whitelist.contains(&status) {
                    return Ok(ResStruct::new(status, body, cookies, header));
                }

                if self.retry < 3 {
                    self.retry += 1;
                    self.proxy = ReqProxy::get();
                    spin_sleep::sleep(std::time::Duration::from_secs(3));
                    return self.request(expected, r, method).await;
                }
                Err("Request failed")
            },
        }
    }
}

fn custom_redirect() -> Policy {
    Policy::custom(|attempt| {
        return if attempt.status() == StatusCode::MOVED_PERMANENTLY {
            log!("Redirecting to {}", attempt.url().as_str());
            attempt.follow()
        } else {
            attempt.stop()
        }
    })
}