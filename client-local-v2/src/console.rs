use std::io;
use std::io::Write;
use chrono::{DateTime, Local};

pub static mut LOG: Vec<Log> = Vec::new();

#[derive(serde::Serialize)]
#[derive(Clone)]
#[derive(Default)]
pub struct Log {
    pub state: i32,
    pub date: DateTime<Local>,
    pub message: String,
}


#[derive(serde::Serialize)]
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
        #[allow(unused_unsafe)]
        unsafe {
            crate::console::LOG.push(crate::console::Log {
                state: 0,
                date,
                message: format!($($arg)*).replace("\n", "").replace("\r", ""),
            });
        }

        print!("[{}] ", date.format("%Y-%m-%d %H:%M:%S"));
        println!($($arg)*);
    };
}

#[macro_export]
macro_rules! error {
    () => {
        print!("\n")
    };
    ($($arg:tt)*) => {
        let date = chrono::Local::now();
        #[allow(unused_unsafe)]
        unsafe {
            crate::console::LOG.push(crate::console::Log {
                state: 1,
                date,
                message: format!($($arg)*).replace("\n", "").replace("\r", ""),
            });
        }

        print!("[{}] ", date.format("%Y-%m-%d %H:%M:%S"));
        println!($($arg)*);
    };
}

pub fn get() -> Option<Websocket> {
    unsafe {
        let log = LOG.clone();
        if log.len() == 0 {
            return None;
        }

        LOG = Vec::new();
        return Some(Websocket {
            op: 11,
            data: log,
        });
    }
}

pub(crate) fn input(prompt: &str) -> String {
    print!("{} ", prompt);

    io::stdout().flush().unwrap();
    let mut s = String::new();

    io::stdin().read_line(&mut s).unwrap();
    s = s.trim().to_string();

    if s == "" {
        log!("Empty input");
        return input(prompt);
    }

    return s
}