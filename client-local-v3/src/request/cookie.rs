use reqwest::Response;

#[derive(serde::Deserialize, serde::Serialize, Clone)]
pub struct Cookie {
    pub name: String,
    pub value: String,
}

#[allow(dead_code)]
impl Cookie {
    pub fn new(name: String, value: String) -> Self {
        Cookie {
            name,
            value,
        }
    }

    pub fn from_string(s: String) -> Self {
        let ns = s.split(";").next().unwrap();
        let mut parts = ns.split('=');
        if let (Some(name), Some(value)) = (parts.next(), parts.next()) {
            return Cookie::new(name.to_string(), value.to_string());
        }

        Cookie::new(s.to_string(), "".to_string())
    }

    pub fn to_string(&self) -> String {
        format!("{}={}; ", self.name, self.value)
    }

    pub fn from_response(res: &Response, old_cookies: Vec<Cookie>) -> Vec<Self> {
        let mut old = old_cookies.clone();
        let set = res.headers().get_all("Set-Cookie");
        for s in set {
            let ns = Cookie::from_string(s.to_str().unwrap().to_string());
            if let Some(i) = old.iter().position(|x| x.name == ns.name) {
                old[i] = ns;
            } else {
                old.push(ns);
            }
        }
        old
    }
}

pub trait Find {
    fn find(&self, name: &'static str) -> String;
}

impl Find for Vec<Cookie> {
    fn find(&self, name: &'static str) -> String {
        for cookie in self {
            if cookie.name == name.to_string() {
                return cookie.clone().value;
            }
        }
        String::from("")
    }
}