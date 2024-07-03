use std::time::SystemTime;
use crate::data::Data;
use crate::{aboutyou, kitheu, payments};
use crate::request::proxy;
use crate::settings::Settings;
use crate::threads::rt;

static mut CACHE: Option<Cache> = None;
static mut IN_USE: bool = false;

pub fn lock() {
    unsafe {
        let start = SystemTime::now();
        while IN_USE {
            spin_sleep::sleep(std::time::Duration::from_millis(10));
            if start.elapsed().unwrap().as_secs() > 2 {
                panic!("Cache lock timeout");
            }
        }
        IN_USE = true;
    }
}

pub fn unlock() {
    unsafe {
        IN_USE = false;
    }
}

#[derive(Clone)]
#[derive(Default)]
struct Cache {
    dev: bool,
    kitheu: kitheu::Cache,
    aboutyou: aboutyou::Cache,
    settings: Settings,
    connected: bool,
    running: i32,
    data: Data,
    proxy: Vec<proxy::ReqProxy>,
    payment: Vec<payments::Payment>,
}

#[derive(Clone)]
#[derive(Default)]
pub struct Captcha {
    pub tokens: Vec<String>,
    pub expire: i64,
}

pub fn check_thread() {
    rt().spawn(async {
        loop {
            rt().spawn(async {
                kitheu::captcha::start().await;
            });
            spin_sleep::sleep(std::time::Duration::from_secs(5));
        }
    });
}

impl Captcha {
    pub fn is_expired(&self) -> bool {
        let now = chrono::Utc::now().timestamp_millis();
        now > self.expire
    }
}

pub fn init() {
    unsafe {
        CACHE = Some(Cache::default());
    }
}

fn get() -> Cache {
    unsafe {
        CACHE.clone().unwrap()
    }
}

fn set(cache: Cache) {
    unsafe {
        CACHE = Some(cache);
    }
}

pub fn set_dev(dev: bool) {
    let mut c = get();
    c.dev = dev;
    set(c);
}

pub fn dev() -> bool {
    get().dev
}

pub fn set_settings(settings: Settings) {
    let mut c = get();
    c.settings = settings;
    set(c);
}

pub fn settings() -> Settings {
    get().settings.clone()
}

pub fn set_connected() {
    let mut c = get();
    c.connected = true;
    set(c);
}

pub fn connected() -> bool {
    get().connected
}

pub fn running() -> i32 {
    get().running
}

pub fn add_running() {
    let mut c = get();
    c.running += 1;
    set(c);
}

pub fn remove_running() {
    let mut c = get();
    c.running -= 1;
    set(c);
}

pub fn set_data(data: Data) {
    let mut c = get();
    c.data = data;
    set(c);
}

pub fn data() -> Data {
    get().data.clone()
}

pub fn set_kitheu(cache: kitheu::Cache) {
    let mut c = get();
    c.kitheu = cache;
    set(c);
}

pub fn kitheu() -> kitheu::Cache {
    get().kitheu.clone()
}

pub fn set_aboutyou(cache: aboutyou::Cache) {
    let mut c = get();
    c.aboutyou = cache;
    set(c);
}

pub fn aboutyou() -> aboutyou::Cache {
    get().aboutyou.clone()
}

pub fn set_proxies(proxies: Vec<proxy::ReqProxy>) {
    let mut c = get();
    c.proxy = proxies;
    set(c);
}

pub fn proxies() -> Vec<proxy::ReqProxy> {
    get().proxy.clone()
}

pub fn set_payments(payments: Vec<payments::Payment>) {
    let mut c = get();
    c.payment = payments;
    set(c);
}

pub fn payments() -> Vec<payments::Payment> {
    get().payment.clone()
}
