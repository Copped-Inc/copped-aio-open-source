use chrono::{DateTime, Local};
use hyper::body::Bytes;
use serde_json::Result;
use crate::{cache, log, settings};

#[allow(dead_code)]

pub fn parse(b: Bytes) {
    let d: OpFour = serde_json::from_slice(&b).unwrap();
    let mut data = d.data.clone();

    if data.instances.is_none() { data.instances = Some(Vec::new()); }
    if data.session.is_none() { data.session = Some(Session::default()); }
    if data.billing.is_none() { data.billing = Some(Vec::new()); }
    if data.shipping.is_none() { data.shipping = Some(Vec::new()); }
    if data.settings.is_none() { data.settings = Some(Settings::default()); }

    if data.settings.clone().unwrap().stores.is_none() {
        data.settings = Some(Settings {
            stores: Some(Stores::default()),
            ..data.settings.unwrap()
        });
    }

    if data.settings.clone().unwrap().stores.unwrap().kith_eu.is_none() {
        data.settings = Some(Settings {
            stores: Some(Stores {
                kith_eu: Some(false),
                ..data.settings.clone().unwrap().stores.unwrap()
            }),
            ..data.settings.clone().unwrap()
        });
    }

    cache::set_data(data);
}

pub fn update(b: Bytes) {
    cache::lock();
    if cache::data().instances.is_none() {
        cache::unlock();
        return;
    }

    let t: Result<OpTwo> = serde_json::from_slice(&b);
    if let Ok(t) = t {
        handle_op_two(t);
        cache::unlock();
        return;
    }
    
    let t: Result<OpTwoInstances> = serde_json::from_slice(&b);
    if let Ok(t) = t {
        handle_op_two_instances(t);
        cache::unlock();
        return;
    }

    cache::unlock();
}

fn handle_op_two(t: OpTwo) {
    let d = cache::data();
    match t.data.action {
        UPDATE_STORES => {
            d.clone().update_stores(t.data.body.store.unwrap(), t.data.body.value.unwrap())
        }
        UPDATE_SESSION => {
            d.clone().update_session(t.data.body.session.unwrap());
        }
        UPDATE_BILLING => {
            d.clone().update_billing(t.data.body.billing.unwrap())
        }
        UPDATE_SHIPPING => {
            d.clone().update_shipping(t.data.body.shipping.unwrap())
        }
        _ => {}
    }
}

fn handle_op_two_instances(t: OpTwoInstances) {
    let d = cache::data();
    match t.data.action {
        UPDATE_INSTANCES => {
            d.clone().update_instances(t.data.body)
        }
        _ => {}
    }
}

const ADD_WEBHOOK: i32 = 1;
const DELETE_WEBHOOK: i32 = 2;
const UPDATE_STORES: i32 = 3;
const UPDATE_INSTANCES: i32 = 4;
const UPDATE_SESSION: i32 = 5;
const UPDATE_CHECKOUTS: i32 = 6;
const UPDATE_BILLING: i32 = 7;
const UPDATE_SHIPPING: i32 = 8;
const ADD_WHITELIST: i32 = 12;
const REMOVE_WHITELIST: i32 = 13;

#[derive(serde::Deserialize)]
#[derive(serde::Serialize)]
pub(crate) struct OpTwo {
    pub op: i32,
    pub data: OpTwoData,
}

#[derive(serde::Deserialize)]
#[derive(serde::Serialize)]
pub(crate) struct OpTwoData {
    pub action: i32,
    pub body: Body,
}

#[derive(serde::Deserialize)]
#[derive(serde::Serialize)]
pub(crate) struct Body {
    pub webhook: Option<String>,
    pub store: Option<String>,
    pub value: Option<bool>,
    pub session: Option<Session>,
    pub billing: Option<Vec<Billing>>,
    pub shipping: Option<Vec<Shipping>>,
    pub checkouts: Option<Vec<Product>>,
    pub product: Option<String>,
}

#[derive(serde::Deserialize)]
#[derive(serde::Serialize)]
pub(crate) struct OpTwoInstances {
    pub op: i32,
    pub data: OpTwoDataInstances,
}

#[derive(serde::Deserialize)]
#[derive(serde::Serialize)]
pub(crate) struct OpTwoDataInstances {
    pub action: i32,
    pub body: Vec<Instance>,
}

#[derive(serde::Deserialize)]
#[derive(serde::Serialize)]
pub struct OpFour {
    pub op: i32,
    pub data: Data,
}

#[derive(serde::Deserialize)]
#[derive(serde::Serialize)]
#[derive(Clone)]
#[derive(Default)]
pub struct Data {
    pub instances: Option<Vec<Instance>>,
    pub settings: Option<Settings>,
    pub session: Option<Session>,
    pub billing: Option<Vec<Billing>>,
    pub shipping: Option<Vec<Shipping>>,
}

impl Data {
    fn update_stores(&mut self, store: String, value: bool) {
        let mut k = self.settings.as_mut().unwrap().stores.clone().unwrap().kith_eu.unwrap();
        match store.as_str() {
            "kith_eu" => k = value,
            _ => {}
        }
        self.settings = Some(Settings {
            stores: Some(Stores {
                kith_eu: Some(k),
            }),
            ..self.settings.as_ref().unwrap().clone()
        });
        cache::set_data(self.clone());
    }

    fn update_instances(&mut self, instances: Vec<Instance>) {
        self.instances = Some(instances);
        cache::set_data(self.clone());
        for i in self.instances.as_ref().unwrap() {
            if i.id == cache::settings().id {
                log!("Updated instance: {}", i.status);
                return;
            }
        }

        log!("Instance deleted");
        settings::delete();
    }

    fn update_session(&mut self, session: Session) {
        self.session = Some(session);
        cache::set_data(self.clone());
    }

    fn update_billing(&mut self, billing: Vec<Billing>) {
        self.billing = Some(billing);
        cache::set_data(self.clone());
    }

    fn update_shipping(&mut self, shipping: Vec<Shipping>) {
        self.shipping = Some(shipping);
        cache::set_data(self.clone());
    }
}

#[derive(serde::Deserialize)]
#[derive(serde::Serialize)]
#[derive(Clone)]
#[derive(Default)]
#[derive(Debug)]
pub struct Product {
    pub date: DateTime<Local>,
    pub name: String,
    pub link: String,
    pub image: String,
    pub store: String,
    pub size: String,
    pub price: f64,
}

#[derive(serde::Deserialize)]
#[derive(serde::Serialize)]
#[derive(Clone)]
#[derive(Default)]
pub struct Instance {
    pub price: f64,
    pub provider: String,
    pub id: String,
    pub status: String,
    pub task_max: String,
    pub region: String,
}

#[derive(serde::Deserialize)]
#[derive(serde::Serialize)]
#[derive(Clone)]
#[derive(Default)]
pub struct Settings {
    pub stores: Option<Stores>,
}

#[derive(serde::Deserialize)]
#[derive(serde::Serialize)]
#[derive(Clone)]
#[derive(Default)]
pub struct Stores {
    pub kith_eu: Option<bool>,
}

#[derive(serde::Deserialize)]
#[derive(serde::Serialize)]
#[derive(Clone)]
#[derive(Default)]
pub struct Session {
    pub status: Option<String>,
    pub checkouts: Option<Vec<Product>>,
    pub declines: Option<Vec<Product>>,
    pub tasks: Option<i32>,
}

#[derive(serde::Deserialize)]
#[derive(serde::Serialize)]
#[derive(Clone)]
#[derive(Default)]
pub struct Billing {
    pub ccnumber: String,
    pub month: String,
    pub year: String,
    pub cvv: String,
}

#[derive(serde::Deserialize)]
#[derive(serde::Serialize)]
#[derive(Clone)]
#[derive(Default)]
pub struct Shipping {
    pub last: String,
    pub address1: String,
    pub address2: String,
    pub email: String,
    pub city: String,
    pub country: String,
    pub state: String,
    pub zip: String,
}

#[derive(serde::Deserialize)]
#[derive(serde::Serialize)]
#[derive(Clone)]
#[derive(Default)]
pub struct User {
    pub name: String,
    pub email: String,
    pub id: String,
    pub picture: String,
    pub plan: i32,
    pub instance_limit: i32,
    pub code_expire: String,
}
