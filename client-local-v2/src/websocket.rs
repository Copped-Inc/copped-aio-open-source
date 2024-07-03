use std::collections::HashMap;
use std::net;
use chrono::{DateTime, Local};
use tungstenite::{Message, stream};
use crate::{aboutyou, api, cache, console, data, error, log, payments};

pub async fn websocket() {
    let mut socket = connect();

    loop {
        let msg = socket.read_message();
        if let Ok(msg) = msg {
            match msg {
                Message::Text(text) => {
                    let body = hyper::body::to_bytes(text).await.unwrap();
                    let l: Basic = serde_json::from_slice(&body).unwrap();

                    match l.op {
                        PING => {
                            socket.write_message(Message::Text(serde_json::to_string(&l).unwrap())).unwrap();
                        }
                        DATA_UPDATE => {
                            data::update(body);
                        }
                        SEND_DATA => {
                            cache::lock();
                            data::parse(body);

                            if !cache::connected() && cache::data().billing.unwrap().len() > 0 && cache::data().shipping.unwrap().len() > 0 {
                                log!("Data received, preloading if needed...");
                                if cache::dev() {
                                    aboutyou::preload();
                                }

                                cache::set_connected();
                            }
                            cache::unlock();
                        }
                        PAYMENTS => {
                            let l: Payment = serde_json::from_slice(&body).unwrap();
                            l.data.update();
                        }
                        NEW_PRODUCT => {
                            let l: Monitor = serde_json::from_slice(&body).unwrap();

                            cache::lock();
                            let settings = cache::settings();
                            let data = cache::data();
                            if data.instances.is_none() {
                                cache::unlock();
                                log!("Skipping Product, no Data received yet");
                                continue;
                            }

                            let running = settings.clone().running();
                            let status = cache::data().session.unwrap().status.unwrap();
                            cache::unlock();

                            if l.data.name.to_lowercase().contains("test") && !cache::dev() || data.shipping.is_none() {
                                log!("Skipping Dev Product or no Data received yet");
                                continue;
                            }

                            if l.data.skus.len() == 0 {
                                log!("Skipping Product with no Skus");
                                continue;
                            }

                            if status == "Running" && running {
                                if l.data.link.contains("kith") && data.settings.unwrap().stores.unwrap().kith_eu.unwrap() {
                                    log!("New Product Kith");
                                    l.data.kith_eu();
                                } else if l.data.link.contains("aboutyou") && cache::dev() /* TODO: add Store */ {
                                    log!("New Product AboutYou");
                                    l.data.aboutyou();
                                }
                            }
                        }
                        NEW_UPDATE => {
                            cache::lock();
                            cache::settings().update().await;
                            cache::unlock();
                        }
                        _ => {}
                    }
                }
                Message::Close(_) => {
                    error!("Disconnected - reconnecting");
                    socket = connect();
                }
                _ => {}
            }

            let logs = console::get();
            if logs.is_some() {
                socket.write_message(Message::Text(serde_json::to_string(&logs.unwrap()).unwrap())).unwrap();
            }
        } else {
            error!("Disconnected - reconnecting");
            socket = connect();
        }
    }
}

fn connect() -> tungstenite::WebSocket<stream::MaybeTlsStream<net::TcpStream>> {
    loop {
        let ws_request = api::websocket();
        let socket = tungstenite::connect(ws_request);
        if socket.is_err() {
            log!("Reconnecting");
            spin_sleep::sleep(std::time::Duration::from_secs(20));
            continue;
        }
        log!("Websocket connected");
        return socket.unwrap().0;
    }
}

const PING: i32 = 1;
const DATA_UPDATE: i32 = 2;
const SEND_DATA: i32 = 4;
const PAYMENTS: i32 = 5;
const NEW_PRODUCT: i32 = 6;
const NEW_UPDATE: i32 = 7;

#[derive(serde::Deserialize)]
#[derive(serde::Serialize)]
pub struct Basic {
    pub op: i32,
}

#[derive(serde::Deserialize)]
#[derive(serde::Serialize)]
pub struct Payment {
    pub op: i32,
    pub data: payments::Payment,
}

#[derive(serde::Deserialize)]
#[derive(serde::Serialize)]
pub struct Monitor {
    pub op: i32,
    pub data: Product,
}

#[derive(serde::Deserialize)]
#[derive(serde::Serialize)]
#[derive(Clone)]
#[derive(Default)]
pub struct Product {
    pub name: String,
    pub sku: String,
    pub skus: HashMap<String, String>,
    pub date: DateTime<Local>,
    pub link: String,
    pub image: String,
    pub price: f64,
}