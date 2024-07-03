use std::io;
use std::thread::sleep;
use reqwest::{Body, Client, StatusCode};
use async_recursion::async_recursion;
use reqwest::redirect::Policy;
use expected::Expected;
use crate::log;
use crate::request::cookie::ReqCookie;
use crate::request::proxy::ReqProxy;
use crate::request::res::ResStruct;

pub(crate) mod res;
pub(crate) mod cookie;
pub(crate) mod proxy;
pub(crate) mod expected;

pub(crate) enum ContentType {
    FormUrlEncoded,
    JSON,
    None,
}

pub(crate) struct ReqStruct {
    pub url: String,
    pub form: Vec<(String, String)>,
    pub body: String,
    pub cookies: Vec<ReqCookie>,
    pub content_type: ContentType,
    pub proxy: ReqProxy,
    retry: i32,
}

impl ReqStruct {
    pub(crate) fn new_with_form(url: String, form: Vec<(String, String)>, cookies: Vec<ReqCookie>, proxy: ReqProxy) -> Self {
        ReqStruct {
            url,
            form,
            body: "".to_string(),
            cookies,
            content_type: ContentType::FormUrlEncoded,
            proxy,
            retry: 0
        }
    }

    pub(crate) fn new_with_json(url: String, body: String, cookies: Vec<ReqCookie>, proxy: ReqProxy) -> Self {
        ReqStruct {
            url,
            form: vec![],
            body,
            cookies,
            content_type: ContentType::JSON,
            proxy,
            retry: 0
        }
    }

    pub(crate) fn new_none_body(url: String, cookies: Vec<ReqCookie>, proxy: ReqProxy) -> Self {
        ReqStruct {
            url,
            form: vec![],
            body: "".to_string(),
            cookies,
            content_type: ContentType::None,
            proxy,
            retry: 0
        }
    }

    #[async_recursion]
    pub(crate) async fn post(mut self, expected: Option<Expected>) -> Result<ResStruct, Box<dyn std::error::Error + Send + Sync>> {
        let custom = custom_redirect();
        let mut client = Client::builder()
            .proxy(self.proxy.proxy).brotli(true).gzip(true)
            .redirect(custom)
            .build()?
            .post(self.url.as_str())
            .header("Host", (self.url.clone().split("://").nth(1).unwrap()).split("/").nth(0).unwrap())
            .header("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36")
            .header("Accept", "*/*")
            .header("Accept-Encoding", "gzip, deflate, br");

        if !self.cookies.is_empty() {
            let mut cookie_str = "".to_string();
            for cookie in &self.cookies {
                cookie_str.push_str(cookie.to_string().as_str());
            }
            client = client.header("Cookie", cookie_str);
        }

        match self.content_type {
            ContentType::FormUrlEncoded => {
                client = client.header("Content-Type", "application/x-www-form-urlencoded")
                    .form(&self.form);
            },
            ContentType::JSON => {
                client = client.header("Content-Type", "application/json")
                    .body(Body::from(self.body.clone()));
            }
            _ => {}
        }

        let res = client.send().await;
        if res.is_err() {
            println!("{}", res.as_ref().unwrap_err());
            if self.retry < 3 {
                self.retry += 1;
                sleep(std::time::Duration::from_secs(3));
                self.proxy = ReqProxy::get();
                return self.post(expected).await;
            }
            return Err(Box::new(io::Error::new(io::ErrorKind::Other, res.unwrap_err())));
        }

        let status = res.as_ref().unwrap().status();
        let cookies = ReqCookie::from_response(res.as_ref().unwrap(), self.cookies.clone());
        let header = res.as_ref().unwrap().headers().clone();
        let body = res.unwrap().text().await.unwrap();

        expected::print_status(status, self.url.clone(), body.clone(), header.clone(), expected.clone());

        match status {
            StatusCode::OK => Ok(ResStruct::new(status, body, cookies, header)),
            StatusCode::FOUND => Ok(ResStruct::new(status, "".to_string(), cookies, header)),
            _ => {
                if self.retry < 3 {
                    self.retry += 1;
                    self.proxy = ReqProxy::get();
                    sleep(std::time::Duration::from_secs(3));
                    return self.post(expected.clone()).await;
                }
                Err(Box::new(io::Error::new(io::ErrorKind::Other, format!("Status: {}", status))))
            },
        }
    }

    #[async_recursion]
    pub(crate) async fn get(mut self, r: bool, expected: Option<Expected>) -> Result<ResStruct, Box<dyn std::error::Error + Send + Sync>> {
        let mut custom = custom_redirect();
        if r {
            custom = Policy::default();
        }

        let mut client = Client::builder()
            .proxy(self.proxy.proxy).brotli(true).gzip(true)
            .redirect(custom)
            .build()?
            .get(self.url.as_str())
            .header("Host", (self.url.clone().split("://").nth(1).unwrap()).split("/").nth(0).unwrap().split("?").nth(0).unwrap())
            .header("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36")
            .header("Accept", "*/*")
            .header("Accept-Encoding", "gzip, deflate, br");

        if !self.cookies.is_empty() {
            let mut cookie_str = "".to_string();
            for cookie in &self.cookies {
                cookie_str.push_str(cookie.to_string().as_str());
            }
            client = client.header("Cookie", cookie_str);
        }
        let res = client.send().await;

        if res.is_err() {
            println!("{}", res.as_ref().unwrap_err());
            if self.retry < 3 {
                self.retry += 1;
                sleep(std::time::Duration::from_secs(3));
                self.proxy = ReqProxy::get();
                return self.get(false, expected).await;
            }
            return Err(Box::new(io::Error::new(io::ErrorKind::Other, res.unwrap_err())));
        }

        let status = res.as_ref().unwrap().status();
        let cookies = ReqCookie::from_response(res.as_ref().unwrap(), self.cookies.clone());
        let header = res.as_ref().unwrap().headers().clone();
        let body = res.unwrap().text().await.unwrap();

        expected::print_status(status, self.url.clone(), body.clone(), header.clone(), expected.clone());

        match status {
            StatusCode::OK => Ok(ResStruct::new(status, body, cookies, header)),
            StatusCode::FOUND => Ok(ResStruct::new(status, "".to_string(), cookies, header)),
            _ => {
                if self.retry < 3 {
                    self.retry += 1;
                    self.proxy = ReqProxy::get();
                    sleep(std::time::Duration::from_secs(3));
                    return self.get(false, expected.clone()).await;
                }
                Err(Box::new(io::Error::new(io::ErrorKind::Other, format!("Status: {}", status))))
            },
        }
    }
}

fn custom_redirect() -> Policy {
    Policy::custom(|attempt| {
        return if attempt.status() == StatusCode::MOVED_PERMANENTLY {
            log(format!("Redirecting to {}", attempt.url().as_str()).as_str());
            attempt.follow()
        } else {
            attempt.stop()
        }
    })
}