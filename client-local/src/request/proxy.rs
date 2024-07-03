use std::ops::Add;
use rand::Rng;
use reqwest::Proxy;
use crate::api::proxies;

pub(crate) static mut PROXIES: Vec<ReqProxy> = vec![];

#[derive(Clone)]
pub(crate) struct ReqProxy {
    pub(crate) proxy: Proxy,
}

impl ReqProxy {
    pub(crate) fn new(ip: String, port: String, username: String, password: String) -> Self {
        let proxy_uri = "http://".to_string()
            .add(username.as_str()).add(":").add(password.as_str())
            .add("@")
            .add(ip.as_str()).add(":").add(port.as_str())
            .to_string();

        let proxy = Proxy::all(proxy_uri).unwrap();

        Self {
            proxy,
        }
    }

    pub(crate) fn clone(&self) -> Self {
        Self {
            proxy: self.proxy.clone(),
        }
    }

    pub(crate) fn get() -> Self {
        let mut rng = rand::thread_rng();
        let index = rng.gen_range(0..unsafe { PROXIES.len() });
        unsafe { PROXIES[index].clone() }
    }

    pub(crate) async unsafe fn get_from_api(auth: String) {
        let proxy_resp = proxies(auth).await;
        if let Ok(p) = proxy_resp {
            PROXIES = vec![];
            let proxies = p.proxies;
            for proxy in proxies {
                PROXIES.push(ReqProxy::new(proxy.ip, proxy.port, proxy.username, proxy.password));
            }
        }
    }
}