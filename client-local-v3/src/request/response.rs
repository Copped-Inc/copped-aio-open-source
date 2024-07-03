use headers::HeaderMap;
use reqwest::{Error, StatusCode};
use crate::request::cookie::Cookie;
use crate::request::Request;

#[derive(Clone)]
pub struct Response {
    pub status: StatusCode,
    pub body: String,
    pub cookies: Vec<Cookie>,
    pub header: HeaderMap,
}

#[allow(dead_code)]
impl Response {
    pub fn new(status_code: StatusCode, body: String, cookies: Vec<Cookie>, header: HeaderMap) -> Self {
        Response {
            status: status_code,
            body,
            cookies,
            header,
        }
    }

    pub async fn from_response(res: Result<reqwest::Response, Error>, req: &Request) -> Self {
        let status = res.as_ref().unwrap().status();
        let cookies = Cookie::from_response(res.as_ref().unwrap(), req.cookies.clone());
        let header = res.as_ref().unwrap().headers().clone();
        let body = if status == StatusCode::FOUND { String::from("") }
                                                    else { res.unwrap().text().await.unwrap() };
        Response::new(status, body, cookies, header)
    }

    pub fn json<'a, T>(&'a self) -> Option<T>
    where T: serde::Deserialize<'a> {
        serde_json::from_str(self.body.as_str()).ok()
    }
}