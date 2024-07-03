use headers::HeaderMap;
use reqwest::StatusCode;
use crate::{error, log};

#[derive(Clone)]
pub struct Expected {
    pub status: Vec<StatusCode>,
    pub body_contains: Vec<String>,
}

#[allow(dead_code)]
impl Expected {
    pub fn new_status(status: StatusCode) -> Option<Self> {
        Some(Expected {
            status: vec![status],
            body_contains: vec![],
        })
    }

    pub fn new_body_contains(body_contains: &str) -> Option<Self> {
        Some(Expected {
            status: vec![],
            body_contains: vec![body_contains.to_string()],
        })
    }

    pub fn new_multiple_status(status: Vec<StatusCode>) -> Option<Self> {
        Some(Expected {
            status,
            body_contains: vec![],
        })
    }

    pub fn new_multiple_body_contains(body_contains: Vec<&str>) -> Option<Self> {
        Some(Expected {
            status: vec![],
            body_contains: body_contains.iter().map(|s| s.to_string()).collect(),
        })
    }
}

pub fn print_status(s: StatusCode, url: String, body: String, header: HeaderMap, expected: Option<Expected>) {
    if expected.is_none() {
        log!("{}: {}", url, s);
        return;
    }
    let expected = expected.unwrap();
    let mut found = false;

    for e in expected.body_contains {
        if body.contains(e.as_str()) {
            found = true;
            break;
        }
    }

    for e in expected.status {
        if s == e {
            found = true;
            break;
        }
    }

    if found {
        if s == StatusCode::FOUND {
            log!("{}: {}, {}", url, s, header.get("Location").unwrap().to_str().unwrap());
        } else{
            log!("{}: {}", url, s);
        }
    } else {
        if s == StatusCode::FOUND {
            error!("{}: {}, {}", url, s, header.get("Location").unwrap().to_str().unwrap());
        } else{
            let body = body.replace("\n", "").replace("\r", "");
            error!("{}: {}\n{}", url, s, body);
        }
    }
}

