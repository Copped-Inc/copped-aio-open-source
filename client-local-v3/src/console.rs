use std::{io, net};
use std::io::Write;
use std::sync::{Arc, Mutex};
use chrono::{DateTime, Local};
use tungstenite::stream;

lazy_static! {
    static ref LOGS: Arc<Mutex<Vec<Log>>> = Arc::new(Mutex::new(Vec::new()));
}
static ERROR_KEYS: [&str; 3] = ["err", "fail", "panic"];

#[derive(serde::Serialize, Clone)]
pub struct Log {
    pub state: i32,
    pub date: DateTime<Local>,
    pub message: String,
}

#[derive(serde::Serialize, Clone)]
pub struct Websocket {
    op: i32,
    data: Vec<Log>,
}

#[macro_export]
macro_rules! log {
    () => {
        print!("\n")
    };
    ($($arg:tt)*) => {
        let date = chrono::Local::now();
        $crate::console::save(&format!($($arg)*), date);
        print!("[{}] ", date.format("%Y-%m-%d %H:%M:%S"));
        println!($($arg)*);
    }
}

pub fn save(arg: &str, date: DateTime<Local>) {
    let state = {
        let a = arg.to_lowercase();
        if ERROR_KEYS.iter().any(|&x| a.contains(x)) {
            1
        } else { 0 }
    };

    let mut logs = LOGS.lock().unwrap();
    logs.push(Log {
        state,
        date,
        message: arg.to_string(),
    });
}

pub fn send(socket: &mut tungstenite::WebSocket<stream::MaybeTlsStream<net::TcpStream>>) {
    let mut logs = LOGS.lock().unwrap();
    if logs.len() > 0 {
        let w = Websocket {
            op: 11,
            data: logs.clone(),
        };

        logs.clear();
        socket.write_message(tungstenite::Message::Text(serde_json::to_string(&w).unwrap())).unwrap();
    }
}

pub fn input(prompt: &str) -> String {
    print!("{} ", prompt);

    io::stdout().flush().unwrap();
    let mut s = String::new();

    io::stdin().read_line(&mut s).unwrap();
    s = s.trim().to_string();

    if s.is_empty() { return input(prompt); }

    return s
}