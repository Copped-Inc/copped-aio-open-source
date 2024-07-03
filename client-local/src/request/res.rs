use reqwest::header::HeaderMap;
use reqwest::StatusCode;
use crate::request::cookie::ReqCookie;

#[allow(dead_code)]
pub(crate) struct ResStruct {
    pub status_code: StatusCode,
    pub body: String,
    pub cookies: Vec<ReqCookie>,
    pub header: HeaderMap,
}

impl ResStruct {
    pub fn new(status_code: StatusCode, body: String, cookies: Vec<ReqCookie>, header: HeaderMap) -> Self {
        ResStruct {
            status_code,
            body,
            cookies,
            header,
        }
    }
}