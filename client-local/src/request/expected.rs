use hyper::StatusCode;
use headers::HeaderMap;
use crate::console::log;

#[derive(Clone)]
pub(crate) struct Expected {
    pub status: Vec<StatusCode>,
    pub body_contains: Vec<String>,
}

#[allow(dead_code)]
impl Expected {
    pub(crate) fn new_status(status: StatusCode) -> Option<Self> {
        Some(Expected {
            status: vec![status],
            body_contains: vec![],
        })
    }

    pub(crate) fn new_body_contains(body_contains: &str) -> Option<Self> {
        Some(Expected {
            status: vec![],
            body_contains: vec![body_contains.to_string()],
        })
    }

    pub(crate) fn new_multiple_status(status: Vec<StatusCode>) -> Option<Self> {
        Some(Expected {
            status,
            body_contains: vec![],
        })
    }

    pub(crate) fn new_multiple_body_contains(body_contains: Vec<&str>) -> Option<Self> {
        Some(Expected {
            status: vec![],
            body_contains: body_contains.iter().map(|s| s.to_string()).collect(),
        })
    }
}

pub(crate) fn print_status(s: StatusCode, url: String, body: String, header: HeaderMap, expected: Option<Expected>) {
    if expected.is_none() {
        log(format!("{}: {}", url, s).as_str());
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
            log(format!("{}: {}, {}", url, s, header.get("Location").unwrap().to_str().unwrap()).as_str());
        } else{
            log(format!("{}: {}", url, s).as_str());
        }
    } else {
        if s == StatusCode::FOUND {
            log(format!("{}: {}, {}", url, s, header.get("Location").unwrap().to_str().unwrap()).as_str());
        } else{
            let body = body.replace("\n", "").replace("\r", "");
            log(format!("{}: {}\n{}", url, s, body).as_str());
        }
    }
}
