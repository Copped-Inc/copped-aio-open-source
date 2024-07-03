use std::net;
use tungstenite::{Message, stream};
use tungstenite::handshake::client;
use crate::{console, data, log, session};
use crate::cache::Cache;

const PING: i32 = 1;
const DATA_UPDATE: i32 = 2;
const SEND_DATA: i32 = 4;
const PAYMENTS: i32 = 5;
const PRODUCT: i32 = 6;
const UPDATE: i32 = 7;

#[derive(serde::Deserialize, serde::Serialize)]
pub struct WebsocketMessage {
    pub op: i32,
}

pub async fn connect() {
    let mut socket = get();
    loop {
        let msg = socket.read_message();
        if msg.is_err() {
            log!("Disconnected - reconnecting");
            socket = get();
            continue;
        }

        let msg = msg.unwrap();
        let res = read(msg, &mut socket).await;
        if res.is_err() {
            log!("Disconnected - reconnecting");
            socket = get();
        }

        console::send(&mut socket);
    }
}

async fn read(msg: Message, socket: &mut tungstenite::WebSocket<stream::MaybeTlsStream<net::TcpStream>>) -> Result<(), ()> {
    match msg {
        Message::Text(text) => {
            let l: WebsocketMessage = serde_json::from_str(text.as_str()).unwrap();
            match l.op {
                PING => socket.write_message(Message::Text(serde_json::to_string(&l).unwrap())).unwrap(),
                DATA_UPDATE => data::update(text),
                SEND_DATA => data::parse(text),
                PAYMENTS => {}
                PRODUCT => session::websocket(text),
                UPDATE => Cache::settings().update().await,
                _ => {}
            }
            Ok(())
        }
        Message::Close(_) => Err(()),
        _ => Ok(()),
    }
}

fn get() -> tungstenite::WebSocket<stream::MaybeTlsStream<net::TcpStream>> {
    loop {
        let ws_request = request();
        let socket = tungstenite::connect(ws_request);
        if socket.is_err() {
            log!("Reconnecting");
            spin_sleep::sleep(std::time::Duration::from_secs(20));
            continue;
        }
        log!("Websocket connected");
        Cache::set_connected();
        return socket.unwrap().0;
    }
}

fn request() -> client::Request {
    let s = Cache::settings();
    client::Request::builder()
        .uri("wss://database.copped-inc.com/websocket"/*"ws://localhost:91/websocket"*/)
        .method("GET")
        .header("Host", "database.copped-inc.com")
        .header("Connection", "Upgrade")
        .header("Upgrade", "websocket")
        .header("Sec-WebSocket-Version", "13")
        .header("Sec-WebSocket-Key", client::generate_key())
        .header("User-Agent", "Instance")
        .header("Price", format!("{:.1}", s.price))
        .header("Provider", format!("{}", s.provider))
        .header("Task-Max", format!("{}", s.task_max))
        .header("Region", format!("{}", s.region))
        .header("Id", format!("{}", s.id))
        .header("Reconnect",
            if Cache::connected() { "1" }
            else { "" }.to_string()
        )
        .header("Cookie", format!("authorization={}", s.authorization))
        .body(())
        .unwrap()
}