use std::ops::Add;
use rand::Rng;
use reqwest::Proxy;
use crate::{api, cache};

pub async fn init() {
    let r = api::proxies().await;
    if let Ok(p) = r {
        let mut proxies = Vec::new();
        for proxy in p.proxies {
            proxies.push(ReqProxy::new(proxy.ip, proxy.port, proxy.username, proxy.password));
        }
        cache::set_proxies(proxies);
    }
}

#[derive(Clone)]
pub struct ReqProxy {
    pub proxy_str: String,
    pub proxy: Proxy,
}

impl ReqProxy {
    pub fn new(ip: String, port: String, username: String, password: String) -> Self {
        let proxy_uri = "http://".to_string()
            .add(username.as_str()).add(":").add(password.as_str())
            .add("@")
            .add(ip.as_str()).add(":").add(port.as_str())
            .to_string();

        let proxy = Proxy::all(proxy_uri).unwrap();

        Self {
            proxy,
            proxy_str: format!("{}:{}:{}:{}", ip, port, username, password),
        }
    }

    pub fn from_str(proxy_str: &String) -> Self {
        let proxy_split: Vec<&str> = proxy_str.split(":").collect();
        let ip = proxy_split[0].to_string();
        let port = proxy_split[1].to_string();
        let username = proxy_split[2].to_string();
        let password = proxy_split[3].to_string();
        Self::new(ip, port, username, password)
    }

    pub(crate) fn get() -> Self {
        let mut rng = rand::thread_rng();
        let index = rng.gen_range(0..cache::proxies().len() );
        cache::proxies()[index].clone()
    }
}