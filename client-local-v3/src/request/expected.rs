use reqwest::StatusCode;
use crate::log;
use crate::request::response::Response;

#[derive(Clone, Default)]
pub struct Expected {
    pub status: Vec<StatusCode>,
    pub body_contains: Vec<String>,
}

#[allow(dead_code)]
impl Expected {
    pub fn status(status: StatusCode) -> Self {
        Expected {
            status: vec![status],
            body_contains: vec![],
        }
    }

    pub fn body_contains(body_contains: &str) -> Self {
        Expected {
            status: vec![],
            body_contains: vec![body_contains.to_string()],
        }
    }

    pub fn multiple_status(status: Vec<StatusCode>) -> Self {
        Expected {
            status,
            body_contains: vec![],
        }
    }

    pub fn multiple_body_contains(body_contains: Vec<&str>) -> Self {
        Expected {
            status: vec![],
            body_contains: body_contains.iter().map(|s| s.to_string()).collect(),
        }
    }

    pub fn check(self, res: Response, url: String, proxy_ip: String) -> Result<(), &'static str> {
        let status = if self.status.len() > 0 { self.status.contains(&res.status) }
                                              else { true };
        let body = if self.body_contains.len() > 0 { self.body_contains.iter().any(|e| res.body.contains(e)) }
                                                   else { true };

        if status && body {
            if res.status == StatusCode::FOUND {
                log!("{}: {} {}, {}", url, &res.status, proxy_ip, res.header.get("Location").unwrap().to_str().unwrap());
            } else{
                log!("{}: {} {}", url, &res.status, proxy_ip);
            }
            Ok(())
        } else {
            if res.status == StatusCode::FOUND {
                log!("Error {}: {} {}, {}", url, &res.status, proxy_ip, res.header.get("Location").unwrap().to_str().unwrap());
            } else{
                let b = res.body.replace("\n", "").replace("\r", "");
                log!("Error {}: {} {}\n{}", url, &res.status, proxy_ip, b);
            }
            Err("Expected check failed")
        }
    }
}
