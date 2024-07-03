use std::collections::HashMap;
use std::net::TcpStream;
use std::ops::{Add};
use std::thread::sleep;
use async_recursion::async_recursion;
use chrono::{DateTime, Local};
use tungstenite::handshake::client::{generate_key, Request, Response};
use tungstenite::{connect, Message, WebSocket};
use serde::{Deserialize, Serialize};
use tungstenite::stream::MaybeTlsStream;
use crate::api::DATABASE_URL;
use crate::{data, log};
use crate::data::values::SESSION;
use crate::modules::ThreadPool;
use crate::request::proxy::ReqProxy;
use crate::settings::values::Settings;

static mut RECONNECT: bool = false;

fn get_websocket_request(s: Settings) -> Request {
    Request::builder()
        .uri(DATABASE_URL.clone().replace("http", "ws").add("websocket"))
        .method("GET")
        .header("Host", (DATABASE_URL.clone().split("://").nth(1).unwrap()).split("/").nth(0).unwrap())
        .header("Connection", "Upgrade")
        .header("Upgrade", "websocket")
        .header("Sec-WebSocket-Version", "13")
        .header("Sec-WebSocket-Key", generate_key())
        .header("User-Agent", "Instance")
        .header("Price", format!("{:.1}", s.price))
        .header("Provider", format!("{}", s.provider))
        .header("Task-Max", format!("{}", s.task_max))
        .header("Region", format!("{}", s.region))
        .header("Id", format!("{}", s.id))
        .header("Reconnect", format!("{}", unsafe {
            if RECONNECT {
                "1"
            } else {
                ""
            }
        }))
        .header("Cookie", format!("authorization={}", s.authorization))
        .body(())
        .unwrap()
}

fn connect_websocket(s: Settings) -> tungstenite::Result<(WebSocket<MaybeTlsStream<TcpStream>>, Response)> {
    loop {
        let ws_request = get_websocket_request(s.clone());
        let socket = connect(ws_request);
        if let Err(_) = socket {
            log("Reconnecting");
            sleep(std::time::Duration::from_secs(20));
            return connect_websocket(s.clone());
        }
        return socket;
    }
}

#[async_recursion]
#[allow(unused_assignments)]
pub(crate) async unsafe fn websocket(s: Settings) {

    let mut data = data::values::Data::new();
    let mut pool = ThreadPool::new(s.clone());
    let mut last_request = Local::now();

    let (mut socket, _) = connect_websocket(s.clone()).unwrap();
    log("Websocket connected");
    RECONNECT = true;

    loop {
        if last_request.timestamp() < Local::now().timestamp() - (60 * 15) {
            ReqProxy::get_from_api(s.authorization.clone()).await;
            last_request = Local::now();
            log("Proxies updated");
        }

        let msg = socket.read_message();
        if let Ok(msg) = msg {
            match msg {
                Message::Text(text) => {
                    let body = hyper::body::to_bytes(text).await.unwrap();
                    let l: BasicResp = serde_json::from_slice(&body).unwrap();
                    match l.op {
                        1 => /* Ping */ {
                            socket.write_message(Message::Text(serde_json::to_string(&l).unwrap())).unwrap();
                        }
                        2 => /* DataUpdate */ {
                            data.update(body, &s);
                        }
                        4 => /* SendData */ {
                            data.parse_data(body);
                        }
                        6 => /* NewProduct */ {
                            let l: MonitorResp = serde_json::from_slice(&body).unwrap();
                            if data.session.status == "Running" && s.clone().running(data.instances.as_ref().unwrap()) {
                                if l.data.link.contains("kith") && data.settings.clone().stores.kith_eu {

                                    pool.send(data.clone(), l.data, 0);
                                    log("New Product Kith");

                                } else if l.data.store == "shopify" && data.settings.clone().stores.shopify {

                                    pool.send(data.clone(), l.data, 1);
                                    log("New Product Shopify");

                                }
                            }
                        }
                        7 => /* NewUpdate */ {
                            s.clone().update().await;
                        }
                        _ => {}
                    }

                    if data.session.tasks != 0 &&
                        data.clone().billing.is_some() &&
                        data.clone().billing.unwrap().len() > 0 &&
                        SESSION.len() == 0 {
                        pool.session_sender.send(data.clone()).unwrap();
                    }
                }
                Message::Close(_) => {
                    log("Disconnected - reconnecting");
                    socket = connect_websocket(s.clone()).unwrap().0;
                }
                _ => {}
            }
        } else {
            log("Disconnected - reconnecting");
            socket = connect_websocket(s.clone()).unwrap().0;
        }
    }
}

#[derive(Deserialize)]
#[derive(Serialize)]
pub(crate) struct BasicResp {
    pub op: i32,
}


#[derive(Deserialize)]
#[derive(Serialize)]
pub(crate) struct MonitorResp {
    pub op: i32,
    pub data: Monitor,
}

#[derive(Deserialize)]
#[derive(Serialize)]
#[derive(Clone)]
pub(crate) struct Monitor {
    pub date: DateTime<Local>,
    pub name: String,
    pub link: String,
    pub stockx: String,
    pub image: String,
    pub store: String,
    pub skus: HashMap<String, String>,
    pub price: f64,
    pub est_sell: f64,
}