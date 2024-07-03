use std::ops::Add;
use std::time::SystemTime;
use chrono::{DateTime, Local};
use rand::Rng;
use crate::{cache, data, session, websocket};
use crate::cache::Captcha;
use crate::request::proxy::ReqProxy;
use crate::threads::rt;

pub mod captcha;
mod module;

#[derive(Clone)]
#[derive(Default)]
pub struct Cache {
    captchas: Vec<Captcha>,
    captcha_required: i64,
    pub sessions: Vec<session::Session>,
}

fn captchas() -> Vec<Captcha> {
    cache::kitheu().captchas.clone()
}

fn captcha() -> Option<Captcha> {
    let mut kitheu = cache::kitheu();
    let mut captchas = kitheu.captchas.clone();
    if captchas.len() > 0 {
        let captcha = captchas.remove(0);
        kitheu.captchas = captchas;
        cache::set_kitheu(kitheu);
        return Some(captcha);
    }
    None
}

fn set_captchas(captchas: Vec<Captcha>) {
    let mut kitheu = cache::kitheu();
    kitheu.captchas = captchas;
    cache::set_kitheu(kitheu);
}

fn append_captcha(captcha: Captcha) {
    let mut kitheu = cache::kitheu();
    kitheu.captchas.push(captcha);
    cache::set_kitheu(kitheu);
}

fn captcha_required() -> bool {
    cache::kitheu().captcha_required > 0
}

fn add_captcha_required() {
    let mut kitheu = cache::kitheu();
    kitheu.captcha_required += 1;
    cache::set_kitheu(kitheu);
}

fn remove_captcha_required() {
    let mut kitheu = cache::kitheu();
    kitheu.captcha_required -= 1;
    cache::set_kitheu(kitheu);
}

fn sessions() -> Vec<session::Session> {
    cache::kitheu().sessions.clone()
}

fn set_sessions(sessions: Vec<session::Session>) {
    let mut kitheu = cache::kitheu();
    kitheu.sessions = sessions;
    cache::set_kitheu(kitheu);
}

impl websocket::Product {
    pub fn kith_eu(self) {
        cache::lock();
        for _ in cache::running()..cache::settings().task_max {
            cache::add_running();
            let product = self.clone();

            let mut rng = rand::thread_rng();
            let sizes: Vec<String> = product.skus.keys().map(|x| x.to_string()).collect();
            let skus: Vec<String> = product.skus.values().map(|x| x.to_string()).collect();
            let i = rng.gen_range(0..skus.len());

            let task = session::Task {
                shipping: data::Shipping::default(),
                billing: data::Billing::default(),
                size: sizes[i].clone(),
                product_id: skus[i].clone(),
                proxy: Some(ReqProxy::get())
            };

            let mut session = session::Session::get_kitheu(task, self.clone());

            session.parse_data();

            rt().spawn(async move {
                session.start_kith_eu().await;
                cache::remove_running();
            });
        }
        cache::unlock();
    }
}

impl session::Session {
    pub fn get_kitheu(task: session::Task, product: websocket::Product) -> Self {
        let mut sessions = sessions();
        let mut ret = session::Session::default();
        ret.product = product.clone();
        ret.task = task.clone();

        sessions.retain(|session| {
            if session.expires < Local::now() {
                return false;
            } else if session.task.product_id == task.product_id && ret.state == 0 {
                ret = session.clone();
                return false;
            }
            true
        });

        set_sessions(sessions);
        ret.expires = Local::now().add(chrono::Duration::days(30));
        ret
    }

    pub fn save_kitheu(&mut self) {
        self.expires = Local::now().add(chrono::Duration::days(30));

        cache::lock();
        let mut s = sessions();

        s.push(self.clone());
        set_sessions(s);
        cache::unlock();

        session::save_sessions();
    }

    pub async fn checked_out_kitheu(self) {
        data::Product {
            date: DateTime::from(SystemTime::now()),
            name: self.product.name.clone(),
            link: "https://eu.kith.com/variants/".to_owned().add(self.task.product_id.as_str()).to_string(),
            image: self.product.image,
            store: "kitheu".to_string(),
            size: self.task.size,
            price: self.product.price,
        }.checkout(self.add_data.get("location").unwrap().as_str()).await;
    }
}