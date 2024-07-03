use std::process::{Command, exit};
use std::sync::mpsc::Receiver;
use hyper::body::Bytes;
use serde_json::Result;
use crate::data::values::{Data, OpFourResp, OpTwoResp, OpTwoRespInstances};
use crate::log;
use crate::settings::delete;
use crate::settings::values::Settings;

pub(crate) mod values;
pub(crate) mod checkout;

#[allow(dead_code)]
pub(crate) static mut PAY_SESSIONS: Option<Receiver<String>> = None;

impl Data {
    pub(crate) unsafe fn parse_data(&mut self, b: Bytes) {
        let d: OpFourResp = serde_json::from_slice(&b).unwrap();
        *self = d.data;
    }

    pub(crate) fn update(&mut self, b: Bytes, s: &Settings) -> Self {
        let t: Result<OpTwoResp> = serde_json::from_slice(&b);
        if let Ok(t) = t {
            match t.data.action {
                1 => /* AddWebhook */ {
                    let mut v = self.user.webhooks.clone().unwrap();
                    v.push(t.data.body.webhook.unwrap());
                    self.user.webhooks = Option::from(v);
                }
                2 => /* DeleteWebhook */ {
                    let mut v = self.user.webhooks.clone().unwrap();
                    let mut i = 0;
                    for w in &v {
                        if w == t.data.body.webhook.as_ref().unwrap() {
                            v.remove(i);
                            break;
                        }
                        i += 1;
                    }
                    self.user.webhooks = Option::from(v);
                }
                3 => /* UpdateStores */ {
                    let store = t.data.body.store.unwrap();
                    if store.clone() == String::from("kith_eu") {
                        self.settings.stores.kith_eu = t.data.body.value.unwrap()
                    } else if store.clone() == String::from("shopify") {
                        self.settings.stores.shopify = t.data.body.value.unwrap()
                    }
                }
                5 => /* UpdateSession */ {
                    self.session = t.data.body.session.unwrap();
                }
                6 => /* UpdateCheckouts */ {
                    self.checkouts = t.data.body.checkouts;
                }
                7|8 => /* UpdateBilling | UpdateShipping */ {
                    let name = std::env::current_exe()
                            .ok()
                            .and_then(|pb| pb.file_name().map(|s| s.to_os_string()))
                            .and_then(|s| s.into_string().ok());

                    let _ = Command::new("cmd.exe")
                        .args(&["/C", "start", name.unwrap().as_str()])
                        .status()
                        .expect("failed to execute process");

                    exit(0);
                }
                _ => {}
            }
        } else {
            let t: Result<OpTwoRespInstances> = serde_json::from_slice(&b);
            if let Ok(t) = t {
                match t.data.action {
                    4 => /* UpdateInstances*/ {
                        self.instances = Option::from(t.data.body);
                        for i in self.instances.as_ref().unwrap() {
                            if i.id == s.id {
                                log(&format!("Updated instance: {}", i.status));
                                return self.clone();
                            }
                        }
                        log("Instance deleted");
                        delete()
                    }
                    _ => {}
                }
            }
        }

        self.clone()
    }
}