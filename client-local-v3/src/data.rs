use crate::cache::Cache;
use crate::{log, modules, settings};

const UPDATE_STORES: i32 = 3;
const UPDATE_INSTANCES: i32 = 4;
const UPDATE_SESSION: i32 = 5;
const UPDATE_BILLING: i32 = 7;
const UPDATE_SHIPPING: i32 = 8;

#[derive(serde::Deserialize, Clone)]
pub struct WebsocketData {
    pub data: Data,
}

#[derive(serde::Deserialize, Clone)]
pub struct WebsocketUpdate {
    pub data: WebsocketUpdateData,
}

#[derive(serde::Deserialize, Clone)]
pub struct WebsocketUpdateData {
    pub action: i32,
    pub body: WebsocketUpdateBody,
}

#[derive(Clone, serde::Deserialize)]
#[serde(untagged)]
pub enum WebsocketUpdateBody {
    UpdateData(UpdateData),
    Instances(Vec<Instance>),
}

#[derive(Clone, Default, serde::Deserialize)]
pub struct UpdateData {
    pub store: Option<String>,
    pub value: Option<bool>,
    pub session: Option<Session>,
    pub billing: Option<Vec<Billing>>,
    pub shipping: Option<Vec<Shipping>>,
}

#[derive(serde::Deserialize, Clone)]
pub struct Data {
    pub instances: Option<Vec<Instance>>,
    pub settings: Option<Settings>,
    pub session: Option<Session>,
    pub billing: Option<Vec<Billing>>,
    pub shipping: Option<Vec<Shipping>>,
}

#[derive(serde::Deserialize, Clone)]
pub struct Instance {
    pub price: f64,
    pub provider: String,
    pub id: String,
    pub status: String,
    pub task_max: String,
    pub region: String,
}

#[derive(serde::Deserialize, Clone, Default)]
pub struct Settings {
    pub stores: Option<Stores>,
}

#[derive(serde::Deserialize, Clone, Default)]
pub struct Stores {
    pub kith_eu: Option<bool>,
    pub queue_it: Option<bool>,
}

#[derive(serde::Deserialize, Clone, Default)]
pub struct Session {
    pub status: String,
}

#[derive(serde::Deserialize, serde::Serialize, Clone, Default)]
pub struct Billing {
    pub ccnumber: String,
    pub month: String,
    pub year: String,
    pub cvv: String,
}

#[derive(serde::Deserialize, serde::Serialize, Clone, Default)]
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

pub fn parse(body: String) {
    let d: WebsocketData = serde_json::from_str(&body.as_str()).unwrap();
    let data = d.data.clone();

    if let Some(s) = data.instances { s.parse(); }
    if let Some(s) = data.settings { s.parse(); }
    if let Some(s) = data.session { s.parse(); }
    if let Some(s) = data.billing { s.parse(); }
    if let Some(s) = data.shipping { s.parse(); }

    Cache::connected();
    return;
}

pub fn update(body: String) {
    let d = serde_json::from_str(&body.as_str());
    if d.is_err() {
        log!("Skipping not handled event: {}", body);
        return;
    }

    let d: WebsocketUpdate = d.unwrap();
    let mut update_data: UpdateData = UpdateData::default();
    let mut instances: Vec<Instance> = Vec::new();

    match d.data.body {
        WebsocketUpdateBody::UpdateData(u) => update_data = u,
        WebsocketUpdateBody::Instances(i) => instances = i,
    }

    match d.data.action {
        UPDATE_STORES => {
            match update_data.store.unwrap().as_str() {
                "kith_eu" => Cache::set_kith_eu(update_data.value.unwrap()),
                "queue_it" => Cache::set_queue_it(update_data.value.unwrap()),
                _ => {}
            }
        }
        UPDATE_INSTANCES => instances.parse(),
        UPDATE_SESSION => update_data.session.unwrap().parse(),
        UPDATE_BILLING => update_data.billing.unwrap().parse(),
        UPDATE_SHIPPING => update_data.shipping.unwrap().parse(),
        _ => {}
    }

    if d.data.action == UPDATE_BILLING || d.data.action == UPDATE_SHIPPING {
        modules::build(true);
    }
}

trait Parse {
    fn parse(self);
}

impl Parse for Vec<Instance> {
    fn parse(self) {
        let id = Cache::settings().id;
        for i in self {
            if i.id == id {
                log!("Instance status: {}", i.status);
                Cache::set_instance_status(i.status == "Running");
                return;
            }
        }

        log!("Instance deleted");
        settings::delete();
    }
}

impl Parse for Settings {
    fn parse(self) {
        if let Some(s) = self.stores {
            s.parse();
        }
    }
}

impl Parse for Stores {
    fn parse(self) {
        if let Some(k) = self.kith_eu {
            Cache::set_kith_eu(k);
        }
        if let Some(k) = self.queue_it {
            Cache::set_queue_it(k);
        }
    }
}

impl Parse for Session {
    fn parse(self) {
        Cache::set_global_status(self.status == "Running");
    }
}

impl Parse for Vec<Billing> {
    fn parse(self) {
        Cache::set_billing(self);
    }
}

impl Parse for Vec<Shipping> {
    fn parse(self) {
        Cache::set_shipping(self);
    }
}
