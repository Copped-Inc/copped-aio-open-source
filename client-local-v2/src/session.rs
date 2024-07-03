use std::collections::HashMap;
use std::fs::File;
use std::io::Write;
use chrono::{DateTime, Local};
use platform_dirs::AppDirs;
use crate::{cache, data, error, log, websocket};
use crate::request::cookie::ReqCookie;
use crate::request::proxy;

#[derive(Clone)]
#[derive(Default)]
pub struct Session {
    pub product: websocket::Product,
    pub state: i32,
    pub cookies: Vec<ReqCookie>,
    pub add_data: HashMap<&'static str, String>,
    pub expires: DateTime<Local>,
    pub task: Task,
}

#[derive(Clone)]
#[derive(Default)]
pub struct Task {
    pub shipping: data::Shipping,
    pub billing: data::Billing,
    pub size: String,
    pub product_id: String,
    pub proxy: Option<proxy::ReqProxy>,
}

impl Session {
    pub fn add_state(&mut self) -> Self {
        self.state += 1;
        self.clone()
    }

    pub fn proxy(&self) -> proxy::ReqProxy {
        self.task.proxy.clone().unwrap()
    }

    fn to_savable(&self) -> SavableSession {
        let mut ad = String::new();
        for (key, value) in self.add_data.iter() {
            ad.push_str(key);
            ad.push_str(":");
            ad.push_str(value);
            ad.push_str(";");
        }

        let t = SavableTask {
            shipping: self.task.shipping.clone(),
            billing: self.task.billing.clone(),
            size: self.task.size.clone(),
            product_id: self.task.product_id.clone(),
            proxy: self.task.proxy.as_ref().unwrap().proxy_str.clone(),
        };

        SavableSession {
            product: self.product.clone(),
            state: self.state,
            cookies: self.cookies.clone(),
            add_data: ad,
            expires: self.expires,
            task: t,
        }
    }
}

#[derive(Clone)]
#[derive(serde::Deserialize)]
#[derive(serde::Serialize)]
struct Sessions {
    pub kitheu: Vec<SavableSession>,
    pub aboutyou: Vec<SavableSession>,
}

#[derive(Clone)]
#[derive(serde::Deserialize)]
#[derive(serde::Serialize)]
struct SavableSession {
    pub product: websocket::Product,
    pub state: i32,
    pub cookies: Vec<ReqCookie>,
    pub add_data: String,
    pub expires: DateTime<Local>,
    pub task: SavableTask,
}

#[derive(Clone)]
#[derive(serde::Deserialize)]
#[derive(serde::Serialize)]
struct SavableTask {
    pub shipping: data::Shipping,
    pub billing: data::Billing,
    pub size: String,
    pub product_id: String,
    pub proxy: String,
}

impl SavableSession {
    fn to_session(self) -> Session {
        let mut ad = HashMap::new();
        for y in self.add_data.split(";") {
            let mut z = y.split(":");
            if z.clone().count() != 2 {
                continue;
            }
            let key = z.next().unwrap();
            let value = z.next().unwrap();

            match key {
                "checkout" => ad.insert("checkout", value.to_string()),
                "cart_token" => ad.insert("cart_token", value.to_string()),
                "location" => ad.insert("location", value.to_string()),
                "email" => ad.insert("email", value.to_string()),
                "password" => ad.insert("password", value.to_string()),
                "token" => ad.insert("token", value.to_string()),
                "item" => ad.insert("item", value.to_string()),
                "secret" => ad.insert("secret", value.to_string()),
                "session" => ad.insert("session", value.to_string()),
                "auth" => ad.insert("auth", value.to_string()),
                "id" => ad.insert("idc", value.to_string()),
                "basket" => ad.insert("basket", value.to_string()),
                "paypal" => ad.insert("paypal", value.to_string()),
                _ => {
                    panic!("Unknown key in add_data: {}", key);
                }
            };
        }

        let t = Task {
            shipping: self.task.shipping.clone(),
            billing: self.task.billing.clone(),
            size: self.task.size.clone(),
            product_id: self.task.product_id.clone(),
            proxy: None,
        };

        Session {
            product: self.product.clone(),
            state: self.state,
            cookies: self.cookies.clone(),
            add_data: ad,
            expires: self.expires,
            task: t,
        }
    }
}

pub fn save_sessions() {
    cache::lock();
    let kitheu = cache::kitheu().sessions.clone();
    let aboutyou = cache::aboutyou().sessions.clone();
    cache::unlock();

    let kitheu_savable: Vec<SavableSession> = kitheu.iter().map(|x| x.to_savable()).collect();
    let aboutyou_savable: Vec<SavableSession> = aboutyou.iter().map(|x| x.to_savable()).collect();

    if kitheu_savable.len() == 0 && aboutyou_savable.len() == 0 {
        return;
    }

    let s = Sessions {
        kitheu: kitheu_savable,
        aboutyou: aboutyou_savable,
    };

    let path = AppDirs::new(Some("Copped AIO"), true).unwrap().
        config_dir.join(format!("sessions{}.json", if cache::dev() { "-dev" } else { "" }));

    let mut file = File::create(path).unwrap();
    file.write_all(serde_json::to_string(&s).unwrap().as_bytes()).unwrap();
}

pub fn load_sessions() {
    let path = AppDirs::new(Some("Copped AIO"), true).unwrap().
        config_dir.join(format!("sessions{}.json", if cache::dev() { "-dev" } else { "" }));

    if !path.exists() {
        return;
    }

    let file = File::open(path).unwrap();
    let s: serde_json::Result<Sessions> = serde_json::from_reader(file);

    if s.is_err() {
        error!("Failed to load sessions.json - updating format");
        cache::unlock();
        return;
    }

    let s = s.unwrap();

    let mut kitheu = Vec::new();
    for x in s.kitheu {
        let mut session = x.clone().to_session();
        if session.expires < Local::now() {
            continue;
        }
        session.task.proxy = Some(proxy::ReqProxy::from_str(&x.task.proxy));
        kitheu.push(session.clone());
    }
    log!("Loaded {} Kitheu sessions", kitheu.len());

    let mut aboutyou = Vec::new();
    for x in s.aboutyou {
        let mut session = x.clone().to_session();
        if session.expires < Local::now() {
            continue;
        }
        session.task.proxy = Some(proxy::ReqProxy::from_str(&x.task.proxy));
        aboutyou.push(session.clone());
    }
    log!("Loaded {} AboutYou sessions", aboutyou.len());

    cache::lock();
    let mut k = cache::kitheu();
    let mut a = cache::aboutyou();

    k.sessions = kitheu;
    a.sessions = aboutyou;

    cache::set_kitheu(k);
    cache::set_aboutyou(a);
    cache::unlock();
}