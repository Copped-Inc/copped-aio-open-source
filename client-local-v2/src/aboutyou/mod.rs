use std::collections::btree_map::BTreeMap;
use std::ops::Add;

use chrono::Local;
use hmac::{Hmac, Mac};
use jwt::SignWithKey;
use jwt::header::HeaderType;
use rand::Rng;
use sha2::Sha256;

use crate::{cache, error, log, session, websocket};
use crate::request::proxy::ReqProxy;
use crate::threads::rt;

mod module;

#[derive(Clone)]
#[derive(Default)]
pub struct Cache {
    pub sessions: Vec<session::Session>,
}

fn sessions() -> Vec<session::Session> {
    cache::aboutyou().sessions.clone()
}

fn set_sessions(sessions: Vec<session::Session>) {
    let mut aboutyou = cache::aboutyou();
    aboutyou.sessions = sessions;
    cache::set_aboutyou(aboutyou);
}

impl websocket::Product {
    pub fn aboutyou(self) {
        cache::lock();
        for _ in cache::running()..cache::settings().task_max {
            cache::add_running();
            let product = self.clone();

            let mut rng = rand::thread_rng();
            let product_id: Vec<String> = product.skus.keys().map(|x| x.to_string()).collect();
            let variant: Vec<String> = product.skus.values().map(|x| x.to_string()).collect();
            let i = rng.gen_range(0..variant.len());

            let mut session = session::Session::get_aboutyou(product_id[i].clone(), variant[i].clone(), self.clone());

            rt().spawn(async move {
                session.start_aboutyou().await;
                cache::remove_running();
            });
        }
        cache::unlock();
    }
}

pub fn preload() {
    let a = cache::aboutyou();

    let t = cache::settings().task_max - a.sessions.len() as i32;
    for _ in 0..t {
        rt().spawn(async move {
            let mut s = session::Session::default();
            s.task.proxy = Some(ReqProxy::get());
            s.parse_data();

            s.task.product_id = "5710945".to_string();
            s.task.size = "43268487".to_string();
            s.preload_aboutyou().await;

            s.task.product_id = "".to_string();
            s.task.size = "".to_string();
            if s.state != 0 {
                error!("failed to preload aboutyou session");
                return;
            }

            log!("preloaded aboutyou session");
            s.save_aboutyou();
        });
    }
}

impl session::Session {
    pub fn get_aboutyou(product_id: String, size: String, product: websocket::Product) -> session::Session {
        let mut sessions = sessions();
        sessions.retain(|x| x.expires > Local::now());

        let ret = sessions.pop();
        if ret.is_none() {
            panic!("no sessions available");
        }

        let mut ret = ret.unwrap();
        ret.product = product.clone();
        ret.task.product_id = product_id;
        ret.task.size = size;

        set_sessions(sessions);
        ret.expires = Local::now().add(chrono::Duration::days(30));
        ret
    }

    pub fn save_aboutyou(&mut self) {
        self.expires = Local::now().add(chrono::Duration::days(30));

        cache::lock();
        let mut s = sessions();

        s.push(self.clone());
        set_sessions(s);
        cache::unlock();

        session::save_sessions();
    }

    pub fn get_session(&mut self) -> String {
        for cookie in self.clone().cookies {
            if cookie.name == "checkout_sid" {
                return cookie.value.replace("%3D", "=");
            }
        }
        String::new()
    }

    pub fn get_secret(&mut self, body: String) -> String {
        let key: Hmac<Sha256> = Hmac::new_from_slice(self.add_data.get("secret").unwrap().as_bytes()).unwrap();
        let header = jwt::Header {
            algorithm: jwt::AlgorithmType::Hs256,
            type_: Option::from(HeaderType::JsonWebToken),
            ..Default::default()
        };

        let mut claims = BTreeMap::new();
        if body.clone() != "" {
            claims.insert("dta", sha256::digest(body));
        }

        jwt::Token::new(header, claims).sign_with_key(&key).unwrap().as_str().to_string()
    }
}

pub fn register_body(mail: String, password: String, first: String, last: String) -> Vec<u8> {

    let mail_bytes = mail.as_bytes();
    let first_bytes = first.as_bytes();
    let last_bytes = last.as_bytes();
    let password_bytes = password.as_bytes();

    let mut bytes = vec![0, 0, 0, 0, (mail_bytes.len() + first_bytes.len() + last_bytes.len() + password_bytes.len() + 19) as u8, 8, 176, 5, 16, 4, 26, mail_bytes.len() as u8];
    bytes.append(&mut mail_bytes.to_vec());
    bytes.append(&mut vec![42, first_bytes.len() as u8]);
    bytes.append(&mut first_bytes.to_vec());
    bytes.append(&mut vec![50, last_bytes.len() as u8]);
    bytes.append(&mut last_bytes.to_vec());
    bytes.append(&mut vec![56, 1, 74, (password_bytes.len() + 2) as u8, 10, password_bytes.len() as u8]);
    bytes.append(&mut password_bytes.to_vec());
    bytes.append(&mut vec![80, 1]);

    bytes

}

pub fn uncart_body(auth: String, basket: String) -> Vec<u8> {

    let auth_bytes = auth.as_bytes();
    let basket_bytes = basket.as_bytes();

    let mut bytes = vec![0, 0, 0, 0, 191, 10, 104, 8, 1, 16, 16, 24, 1, 34, 14, 98, 117, 105, 108, 100, 45, 97, 48, 49, 48, 50, 52, 48, 99, 50, 23, 10, 14, 97, 98, 95, 98, 97, 95, 115, 117, 112, 114, 101, 99, 95, 100, 18, 5, 115, 114, 100, 45, 49, 50, 17, 10, 8, 97, 98, 95, 97, 121, 111, 95, 100, 18, 5, 97, 121, 111, 45, 49, 50, 17, 10, 8, 97, 98, 95, 98, 102, 49, 95, 100, 18, 5, 98, 102, 49, 45, 49, 50, 17, 10, 8, 97, 98, 95, 118, 112, 111, 95, 100, 18, 5, 118, 112, 111, 45, 51, 18, 40, 10, 38, 10, 36];
    bytes.append(&mut auth_bytes.to_vec());
    bytes.append(&mut vec![130, 1, 32]);
    bytes.append(&mut basket_bytes.to_vec());
    bytes.append(&mut vec![136, 1, 1, 146, 1, 2, 10, 0]);

    bytes

}

pub fn checkout_body(session: String, secret: String, auth: String) -> Vec<u8> {

    let session_bytes = session.as_bytes();
    let secret_bytes = secret.as_bytes();
    let auth_bytes = auth.as_bytes();

    let mut bytes = vec![0, 0, 0, 1, 237, 16, 176, 5, 26, 38, 10, 36];
    bytes.append(&mut auth_bytes.to_vec());
    bytes.append(&mut vec![42, 183, 3, 10, 206, 2]);
    bytes.append(&mut session_bytes.to_vec());
    bytes.append(&mut vec![18, 100]);
    bytes.append(&mut secret_bytes.to_vec());
    bytes.append(&mut vec![48, 4, 58, 2, 10, 0, 74, 0]);

    bytes

}

pub fn orders_body(session: String, secret: String, auth: String) -> Vec<u8> {

    let session_bytes = session.as_bytes();
    let secret_bytes = secret.as_bytes();
    let auth_bytes = auth.as_bytes();

    let mut bytes = vec![0, 0, 0, 1, 233, 8, 176, 5, 18, 183, 3, 10, 206, 2];
    bytes.append(&mut session_bytes.to_vec());
    bytes.append(&mut vec![18, 100]);
    bytes.append(&mut secret_bytes.to_vec());
    bytes.append(&mut vec![26, 0, 40, 4, 50, 38, 10, 36]);
    bytes.append(&mut auth_bytes.to_vec());

    bytes

}

pub fn cancel_body(session: String, secret: String, auth: String) -> Vec<u8> {

    let session_bytes = session.as_bytes();
    let secret_bytes = secret.as_bytes();
    let auth_bytes = auth.as_bytes();

    let mut bytes = vec![0, 0, 0, 1, 242, 8, 176, 5, 18, 183, 3, 10, 206, 2];
    bytes.append(&mut session_bytes.to_vec());
    bytes.append(&mut vec![18, 100]);
    bytes.append(&mut secret_bytes.to_vec());
    bytes.append(&mut vec![24, 248, 196, 232, 106, 32, 154, 1, 40, 4, 50, 38, 10, 36]);
    bytes.append(&mut auth_bytes.to_vec());
    bytes.append(&mut vec![66, 1, 1]);

    bytes

}
