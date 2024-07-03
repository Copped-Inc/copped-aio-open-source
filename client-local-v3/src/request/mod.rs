use reqwest::redirect::Policy;
use reqwest::{Client, StatusCode};
use crate::api::Method;
use crate::{log, session};
use crate::request::cookie::Cookie;
use crate::request::expected::Expected;
use crate::request::proxy::Proxy;
use crate::request::response::Response;
use crate::threads::handle;

pub mod cookie;
pub mod proxy;
pub mod expected;
pub mod response;

#[allow(dead_code)]
pub enum ContentType {
    FormUrlEncoded,
    JSON,
    GRCP,
    None,
}

impl Default for ContentType {
    fn default() -> Self {
        ContentType::None
    }
}

#[derive(Default)]
#[allow(dead_code)]
pub struct Request {
    proxy: Proxy,
    url: String,
    form: Vec<(String, String)>,
    body: Vec<u8>,
    cookies: Vec<Cookie>,
    content_type: ContentType,
    tls: bool,
    header: Vec<(String, String)>,
    redirect: bool,
    expected: Expected,
    retry: i32,
}

impl session::Task {
    pub fn request(&mut self) -> Request {
        Request::from_task(self)
    }
}

#[allow(dead_code)]
impl Request {
    pub fn new(proxy: Proxy) -> Self {
        Request {
            proxy,
            ..Default::default()
        }
    }

    pub fn from_task(t: &session::Task) -> Self {
        let mut r = Request::new(t.proxy.clone());
        r.cookies = t.cookies.clone();
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

    pub fn body_json<T>(self, body: &T) -> Self
    where T: serde::Serialize {
        self.body(serde_json::to_string(body).unwrap().into_bytes(), ContentType::JSON)
    }

    pub fn body_grcp(self, body: Vec<u8>) -> Self {
        self.body(body, ContentType::GRCP)
    }

    pub fn cookies(mut self, cookies: Vec<Cookie>) -> Self {
        self.cookies = cookies;
        self
    }

    pub fn tls(mut self) -> Self {
        self.tls = true;
        self
    }

    pub fn add_header(mut self, key: String, value: String) -> Self {
        self.header.push((key, value));
        self
    }

    pub fn redirect(mut self) -> Self {
        self.redirect = true;
        self
    }

    pub fn expected(mut self, expected: Expected) -> Self {
        self.expected = expected;
        self
    }

    pub fn get(self) -> Result<Response, &'static str> {
        self.request(Method::Get)
    }

    pub fn post(self) -> Result<Response, &'static str> {
        self.request(Method::Post)
    }

    pub fn put(self) -> Result<Response, &'static str> {
        self.request(Method::Put)
    }

    pub fn delete(self) -> Result<Response, &'static str> {
        self.request(Method::Delete)
    }

    fn request(self, method: Method) -> Result<Response, &'static str> {
        let redirect = if self.redirect { all_redirect() }
                                          else { no_redirect() };

        let client_builder = Client::builder()
            .proxy(self.proxy())
            .redirect(redirect);

        let client = if self.tls { client_builder.use_rustls_tls() }
                                   else { client_builder }.danger_accept_invalid_certs(true).build().unwrap();

        let mut request_builder = match method {
            Method::Get => client.get(self.url.as_str()),
            Method::Post => client.post(self.url.as_str()),
            Method::Put => client.put(self.url.as_str()),
            Method::Delete => client.delete(self.url.as_str()),
        };

        request_builder = request_builder.header("Host", (self.url.clone().split("://").nth(1).unwrap()).split("/").nth(0).unwrap())
            .header("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36")
            .header("Accept", "*/*")
            .header("Accept-Encoding", "gzip, deflate, br")
            .header("sec-ch-ua", "\"Google Chrome\";v=\"107\", \"Chromium\";v=\"107\", \";Not A Brand\";v=\"24\"")
            .header("sec-ch-ua-mobile", "?0")
            .header("sec-ch-ua-platform", "\"Windows\"")
            .header("sec-fetch-dest", "empty")
            .header("sec-fetch-mode", "cors")
            .header("sec-fetch-site", "same-origin");

        for (key, value) in &self.header {
            request_builder = request_builder.header(key, value);
        }

        if !self.cookies.is_empty() {
            let mut cookie = String::new();
            for c in &self.cookies {
                cookie.push_str(c.to_string().as_str());
            }
            request_builder = request_builder.header("Cookie", cookie);
        }

        if method == Method::Post || method == Method::Put {
            match self.content_type {
                ContentType::FormUrlEncoded => request_builder = request_builder.header("Content-Type", "application/x-www-form-urlencoded")
                    .form(&self.form),
                ContentType::JSON => request_builder = request_builder.header("Content-Type", "application/json")
                    .body(self.body.clone()),
                ContentType::GRCP => request_builder = request_builder.header("Content-Type", "application/grpc-web+proto")
                    .header("x-grpc-web", "1")
                    .body(self.body.clone()),
                _ => {}
            }
        }

        let res = handle().block_on(request_builder.send());
        if res.is_err() { return self.return_err(res.unwrap_err().to_string(), method); }

        let res = handle().block_on(Response::from_response(res, &self));
        let check = self.expected.clone().check(res.clone(), self.url.clone(), self.proxy.ip.clone());
        if check.is_err() { return self.return_err(check.unwrap_err().to_string(), method); }

        return Ok(res);
    }

    fn return_err(mut self, e: String, method: Method) -> Result<Response, &'static str> {
        log!("Request error: {}", e);
        if self.retry >= 3 {
            return Err("Request failed");
        }

        self.retry += 1;
        spin_sleep::sleep(std::time::Duration::from_secs(3));
        self.proxy = Proxy::get();
        return self.request(method);
    }
}

fn no_redirect() -> Policy {
    Policy::custom(|attempt| {
        return if attempt.status() == StatusCode::MOVED_PERMANENTLY {
            log!("Redirecting to {}", attempt.url().as_str());
            attempt.follow()
        } else {
            attempt.stop()
        }
    })
}

fn all_redirect() -> Policy {
    Policy::custom(|attempt| {
        if attempt.previous().len() > 10 {
            attempt.error("too many redirects")
        } else {
            log!("Redirecting to {}", attempt.url().as_str());
            attempt.follow()
        }
    })

}
