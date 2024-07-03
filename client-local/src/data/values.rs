use std::collections::HashMap;
use std::ops::Add;
use std::sync::mpsc;
use std::sync::mpsc::{Receiver, Sender};
use std::thread::sleep;
use chrono::{DateTime, Local};
use serde::{Serialize, Deserialize};
use tokio::runtime;
use crate::jig::random_name;
use crate::modules::shopify::{CreditCard, SessionReq};
use crate::modules::Task;
use crate::{log, ReqProxy};
use crate::request::ReqStruct;

#[derive(Deserialize)]
#[derive(Serialize)]
pub(crate) struct OpTwoResp {
    pub op: i32,
    pub data: OpTwoData,
}

#[derive(Deserialize)]
#[derive(Serialize)]
pub(crate) struct OpTwoData {
    pub action: i32,
    pub body: Body,
}

#[derive(Deserialize)]
#[derive(Serialize)]
pub(crate) struct Body {
    pub webhook: Option<String>,
    pub store: Option<String>,
    pub value: Option<bool>,
    pub session: Option<Session>,
    pub billing: Option<Vec<Billing>>,
    pub shipping: Option<Vec<Shipping>>,
    pub checkouts: Option<Vec<Product>>,
}

#[derive(Deserialize)]
#[derive(Serialize)]
pub(crate) struct OpTwoRespInstances {
    pub op: i32,
    pub data: OpTwoDataInstances,
}

#[derive(Deserialize)]
#[derive(Serialize)]
pub(crate) struct OpTwoDataInstances {
    pub action: i32,
    pub body: Vec<Instance>,
}

#[derive(Deserialize)]
#[derive(Serialize)]
pub(crate) struct OpFourResp {
    pub op: i32,
    pub data: Data,
}

#[derive(Deserialize)]
#[derive(Serialize)]
#[derive(Clone)]
pub(crate) struct Data {
    pub checkouts: Option<Vec<Product>>,
    pub declines: Option<Vec<Product>>,
    pub instances: Option<Vec<Instance>>,
    pub settings: Settings,
    pub session: Session,
    pub billing: Option<Vec<Billing>>,
    pub shipping: Option<Vec<Shipping>>,
    pub user: User,
}

#[derive(Deserialize)]
#[derive(Serialize)]
#[derive(Clone)]
pub(crate) struct Product {
    pub date: DateTime<Local>,
    pub name: String,
    pub link: String,
    pub image: String,
    pub store: String,
    pub size: String,
    pub price: f64,
    pub est_sell: f64,
}

#[derive(Deserialize)]
#[derive(Serialize)]
#[derive(Clone)]
pub(crate) struct Instance {
    pub price: f64,
    pub provider: String,
    pub id: String,
    pub status: String,
    pub task_max: String,
    pub region: String,
}

#[derive(Deserialize)]
#[derive(Serialize)]
#[derive(Clone)]
pub(crate) struct Settings {
    pub stores: Stores,
}

#[derive(Deserialize)]
#[derive(Serialize)]
#[derive(Clone)]
pub(crate) struct Stores {
    pub kith_eu: bool,
    pub shopify: bool,
}

#[derive(Deserialize)]
#[derive(Serialize)]
#[derive(Clone)]
pub(crate) struct Session {
    pub status: String,
    pub checkouts: Option<Vec<Product>>,
    pub declines: Option<Vec<Product>>,
    pub tasks: i32,
}

#[derive(Deserialize)]
#[derive(Serialize)]
#[derive(Clone)]
pub(crate) struct Billing {
    pub ccnumber: String,
    pub month: String,
    pub year: String,
    pub cvv: String,
}

#[derive(Deserialize)]
#[derive(Serialize)]
#[derive(Clone)]
pub(crate) struct Shipping {
    pub last: String,
    pub address1: String,
    pub address2: String,
    pub email: String,
    pub city: String,
    pub country: String,
    pub state: String,
    pub zip: String,
}

#[derive(Deserialize)]
#[derive(Serialize)]
#[derive(Clone)]
pub(crate) struct User {
    pub name: String,
    pub email: String,
    pub id: String,
    pub picture: String,
    pub plan: Plan,
    pub instance_limit: i32,
    pub code_expire: String,
    pub webhooks: Option<Vec<String>>,
}

#[derive(Deserialize)]
#[derive(Serialize)]
#[derive(Clone)]
pub(crate) struct Plan {
    pub id: i32,
    pub name: String,
    pub monthly_fee: f64,
}

pub(crate) static mut SESSION: Vec<String> = Vec::new();
pub(crate) static mut RATES: Option<HashMap<String, HashMap<String, String>>> = None;
static STORES: [&str; 3] = [
    "de.slamjam.com",
    "en.afew-store.com",
    "www.asphaltgold.com"
];

pub(crate) static PAYMENT: [&str; 3] = [
    "66119860408",
    "36440014891",
    "58970603684"
];

pub(crate) static HIGH_SECURE: [&str; 3] = [
    "0",
    "0",
    "1"
];

pub(crate) fn get_store_setting(a: [&str; 3], s: String) -> String {
    return if s.contains("slamjam") {
        a[0].to_string()
    } else if s.contains("afew") {
        a[1].to_string()
    } else if s.contains("asphaltgold") {
        a[2].to_string()
    } else {
        "error".to_string()
    };
}

impl Data {
    pub(crate) fn new() -> Data {
        let c: [Product; 0] = [];
        let d: [Product; 0] = [];
        let i: [Instance; 0] = [];
        let b: [Billing; 0] = [];
        let s: [Shipping; 0] = [];
        Data {
            checkouts: Option::from(Vec::from(c)),
            declines: Option::from(Vec::from(d)),
            instances: Option::from(Vec::from(i)),
            settings: Settings::new(),
            session: Session::new(),
            billing: Option::from(Vec::from(b)),
            shipping: Option::from(Vec::from(s)),
            user: User::new(),
        }
    }

    pub(crate) async unsafe fn get_shipping_rates(&mut self) {
        let (_, max_tokio_blocking_threads) = (num_cpus::get(), 512);
        let rt = runtime::Builder::new_multi_thread()
            .enable_all()
            .thread_stack_size(8 * 1024 * 1024)
            .worker_threads(STORES.clone().len())
            .max_blocking_threads(max_tokio_blocking_threads)
            .build();

        let rt = rt.unwrap();
        let (tx, rx): (Sender<bool>, Receiver<bool>) = mpsc::channel();

        if RATES.is_none() {
            RATES = Some(HashMap::new());
        }

        for s in STORES.clone() {
            let shipping_obj = self.shipping.clone();
            let sender = tx.clone();

            rt.spawn(async move {
                let mut zip_rates = HashMap::new();
                for shipping in shipping_obj.clone().unwrap() {
                    if zip_rates.get(&shipping.zip).is_none() {
                        let rate = shipping.clone().get_rate("https://".to_owned().add(s).add("/")).await;
                        if rate.is_empty() {
                            log(format!("No rate for {} on {}", shipping.zip, s).as_str());
                        }
                        zip_rates.insert(shipping.zip.clone(), rate);
                    }
                }

                let mut rates = RATES.clone().unwrap();
                rates.insert(s.to_string(), zip_rates.clone());
                RATES = Some(rates);
                let _ = sender.send(true);
            });
        }

        for _ in 0..STORES.clone().len() {
            rx.recv().unwrap();
        }

        rt.shutdown_background();
    }

    pub(crate) async unsafe fn gen_session(self, t: usize) {
        let mut last = Local::now();
        loop {
            if Local::now().signed_duration_since(last).num_minutes() >= 5 {
                SESSION = Vec::new();
                last = Local::now();
            }

            if SESSION.len() >= t {
                sleep(std::time::Duration::from_secs(1));
                continue;
            }

            let b = self.billing.as_ref().unwrap()[0].clone();

            let m = b.month.parse::<i64>().unwrap();
            let y = b.year.parse::<i64>().unwrap();

            let req = SessionReq {
                credit_card: CreditCard {
                    number: b.ccnumber.clone(),
                    name: random_name().add(" ").add(random_name().as_str()),
                    start_month: m.clone(),
                    start_year: y.clone(),
                    month: m,
                    year: y,
                    verification_value: b.clone().cvv,
                    issue_number: b.clone().ccnumber
                }
            };

            let req = ReqStruct::new_with_json(
                "https://deposit.us.shopifycs.com/sessions".to_string(),
                serde_json::to_string(&req).unwrap(),
                vec![],
                ReqProxy::get()
            );

            let res = req.post(None).await;

            if let Ok(res) = res {
                let id = res.body[7..res.body.len() - 2].to_string();
                SESSION.push(id);
            }
        }

    }
}

pub(crate) fn get_session() -> String {
    unsafe {
        if SESSION.len() > 0 {
            let s = SESSION.pop().unwrap();
            return s;
        }

        "".to_string()
    }
}

pub(crate) fn get_rates(t: Task) -> String {
    unsafe {
        let r = RATES.as_ref().unwrap().get(&t.store.to_string());
        if r.is_some() {
            let r = r.unwrap();
            let r = r.get(&t.shipping.zip.to_string());
            if r.is_some() {
                return r.unwrap().to_string();
            }
        }

        "".to_string()
    }
}

impl Settings {
    pub(crate) fn new() -> Settings {
        let s: Stores = Stores::new();
        Settings {
            stores: s,
        }
    }
}


impl Stores {
    pub(crate) fn new() -> Stores {
        Stores {
            kith_eu: false,
            shopify: false
        }
    }
}

impl Session {
    pub(crate) fn new() -> Session {
        let c: [Product; 0] = [];
        let d: [Product; 0] = [];
        Session {
            status: String::new(),
            checkouts: Option::from(Vec::from(c)),
            declines: Option::from(Vec::from(d)),
            tasks: 0,
        }
    }
}

impl Shipping {
    pub(crate) fn new() -> Shipping {
        Shipping {
            last: String::new(),
            address1: String::new(),
            address2: String::new(),
            email: String::new(),
            city: String::new(),
            country: String::new(),
            state: String::new(),
            zip: String::new(),
        }
    }
}

impl Billing {
    pub(crate) fn new() -> Billing {
        Billing {
            ccnumber: String::new(),
            month: String::new(),
            year: String::new(),
            cvv: String::new(),
        }
    }
}

impl User {
    pub(crate) fn new() -> User {
        User {
            name: String::new(),
            email: String::new(),
            id: String::new(),
            picture: String::new(),
            plan: Plan::new(),
            instance_limit: 0,
            code_expire: String::new(),
            webhooks: Option::from(Vec::from(vec![])),
        }
    }
}

impl Plan {
    pub(crate) fn new() -> Plan {
        Plan {
            id: 0,
            name: String::new(),
            monthly_fee: 0.0,
        }
    }
}