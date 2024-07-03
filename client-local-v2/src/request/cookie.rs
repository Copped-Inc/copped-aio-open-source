use reqwest::Response;

#[derive(Clone)]
#[derive(serde::Serialize)]
#[derive(serde::Deserialize)]
pub struct ReqCookie {
    pub name: String,
    pub value: String,
}

impl ReqCookie {
    pub fn new(name: String, value: String) -> Self {
        ReqCookie {
            name,
            value,
        }
    }

    pub fn from_string(s: String) -> Self {
        let ns = s.split(";").next().unwrap();
        let mut parts = ns.split('=');
        if let (Some(name), Some(value)) = (parts.next(), parts.next()) {
            ReqCookie::new(name.to_string(), value.to_string())
        } else {
            ReqCookie::new(s.to_string(), "".to_string())
        }
    }

    pub fn to_string(&self) -> String {
        format!("{}={}; ", self.name, self.value)
    }

    pub fn from_response(res: &Response, old: Vec<ReqCookie>) -> Vec<Self> {
        let set = res.headers().get_all("Set-Cookie");
        let mut new = Vec::new();
        for s in set {
            let ns = ReqCookie::from_string(s.to_str().unwrap().to_string());
            new.push(ns);
        }

        let mut cookies = old;
        for c in new {
            if let Some(i) = cookies.iter().position(|x| x.name == c.name) {
                cookies[i] = c;
            } else {
                cookies.push(c);
            }
        }

        cookies
    }
}