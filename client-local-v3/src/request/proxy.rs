use std::ops::Add;
use rand::Rng;
use reqwest::Proxy as reqProxy;
use crate::{api, log};
use crate::cache::Cache;
use crate::request::Request;

pub async fn init() {
    let r = api::proxies().await;
    if r.is_err() {
        return;
    }

    let mut proxies = Vec::new();
    for proxy in r.unwrap().proxies {
        proxies.push(Proxy::new(proxy.ip, proxy.port, proxy.username, proxy.password));
    }

    log!("Got Proxies: {}", proxies.len());
    Cache::set_proxies(proxies);
}


#[derive(serde::Deserialize, serde::Serialize, Clone, Default)]
pub struct Proxy {
    pub ip: String,
    pub port: String,
    pub username: String,
    pub password: String,
}

#[derive(serde::Deserialize)]
pub struct ResProxy {
    pub ip: String,
    pub port: String,
    pub username: String,
    pub password: String,
}

#[allow(dead_code)]
impl Proxy {
    pub fn new(ip: String, port: String, username: String, password: String) -> Self {
        Self {
            ip,
            port,
            username,
            password,
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
        let index = rng.gen_range(0..Cache::proxies().len() );
        Cache::proxies()[index].clone()
    }
}

impl Request {
    pub fn proxy(&self) -> reqProxy {
        let proxy_uri = "http://".to_string()
            .add(self.proxy.username.as_str()).add(":").add(self.proxy.password.as_str())
            .add("@")
            .add(self.proxy.ip.as_str()).add(":").add(self.proxy.port.as_str())
            .to_string();

        reqProxy::all(proxy_uri).unwrap()
    }
}