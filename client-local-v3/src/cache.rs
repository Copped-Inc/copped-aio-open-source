use std::env;
use std::sync::{Arc, Mutex};
use crate::{data, log, session};
use crate::request::proxy;
use crate::settings::Settings;
use crate::threads::WaitLock;

lazy_static! {
    static ref CACHE: Arc<Mutex<Cache>> = Arc::new(Mutex::new(Cache::default()));
}

#[derive(Default)]
pub struct Cache {
    dev: bool,
    connected: bool,
    instance_status: bool,
    global_status: bool,
    stores: Stores,
    settings: Settings,
    sessions: Vec<session::Sessions>,
    proxies: Vec<proxy::Proxy>,
    billing: Vec<data::Billing>,
    shipping: Vec<data::Shipping>,
}

#[derive(Default)]
pub struct Stores {
    pub kith_eu: bool,
    pub queue_it: bool,
}

#[allow(dead_code)]
impl Cache {
    pub fn init() {
        let mut c = CACHE.wait_lock();
        *c = Cache::default();

        let args: Vec<String> = env::args().collect();
        if args.len() > 1 && args[1] == "dev" {
            log!("Running in dev mode");
            c.dev = true;
        }
    }

    pub fn dev() -> bool {
        CACHE.wait_lock().dev
    }

    pub fn set_connected() {
        let mut c = CACHE.wait_lock();
        c.connected = true;
    }

    pub fn connected() -> bool {
        CACHE.wait_lock().connected
    }

    pub fn set_instance_status(status: bool) {
        let mut c = CACHE.wait_lock();
        c.instance_status = status;
    }

    pub fn set_global_status(status: bool) {
        let mut c = CACHE.wait_lock();
        c.global_status = status;
    }

    pub fn set_kith_eu(b: bool) {
        let mut c = CACHE.wait_lock();
        c.stores.kith_eu = b;
    }

    pub fn kith_eu() -> bool {
        CACHE.wait_lock().stores.kith_eu && CACHE.wait_lock().instance_status && CACHE.wait_lock().global_status
    }

    pub fn set_queue_it(b: bool) {
        let mut c = CACHE.wait_lock();
        c.stores.queue_it = b;
    }

    pub fn queue_it() -> bool {
        CACHE.wait_lock().stores.queue_it && CACHE.wait_lock().instance_status && CACHE.wait_lock().global_status
    }

    pub fn set_settings(settings: Settings) {
        let mut c = CACHE.wait_lock();
        c.settings = settings;
    }

    pub fn settings() -> Settings {
        CACHE.wait_lock().settings.clone()
    }

    pub fn set_sessions(session: Vec<session::Sessions>) {
        let mut c = CACHE.wait_lock();
        c.sessions = session;
    }

    pub fn add_task(mut task: session::Task) {
        let mut c = CACHE.wait_lock();
        for s in c.sessions.iter_mut() {
            if s.module == task.module.clone() {
                if s.fallback != 0 && s.fallback < task.state {
                    task.state = s.fallback;
                }

                for delete_value in s.delete_values.iter() {
                    task.data.remove(delete_value);
                }

                s.tasks.push(task);
                return;
            }
        }
    }

    pub fn sessions() -> Vec<session::Sessions> {
        CACHE.wait_lock().sessions.clone()
    }

    pub fn update_session(session: session::Sessions) {
        let mut c = CACHE.wait_lock();
        for s in c.sessions.iter_mut() {
            if s.module == session.module {
                *s = session;
                return;
            }
        }
        c.sessions.push(session);
    }

    pub fn set_proxies(proxies: Vec<proxy::Proxy>) {
        let mut c = CACHE.wait_lock();
        c.proxies = proxies;
    }

    pub fn proxies() -> Vec<proxy::Proxy> {
        CACHE.wait_lock().proxies.clone()
    }

    pub fn set_billing(billing: Vec<data::Billing>) {
        let mut c = CACHE.wait_lock();
        c.billing = billing;
    }

    pub fn billing() -> Vec<data::Billing> {
        CACHE.wait_lock().billing.clone()
    }

    pub fn set_shipping(shipping: Vec<data::Shipping>) {
        let mut c = CACHE.wait_lock();
        c.shipping = shipping;
    }

    pub fn shipping() -> Vec<data::Shipping> {
        CACHE.wait_lock().shipping.clone()
    }
}