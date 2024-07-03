#![allow(dead_code)]
#![allow(unused_variables)]
#![allow(unused_imports)]
#![allow(unused_must_use)]

use std::fmt::format;
use std::ops::Add;
use std::sync::mpsc;
use std::sync::mpsc::{Receiver, Sender};
use std::thread;
use std::thread::sleep;
use std::time::Duration;
use systemstat::{System, Platform, saturating_sub_bytes};
use futures::executor::block_on;
use futures::task::SpawnExt;
use serde_json::from_slice;
use tokio::runtime;
use tokio::runtime::Runtime;
use tokio::time::*;
use crate::api::performance;
use crate::log;

pub(crate) async fn test() -> i32 {
    let mut i: i32 = 2000;
    let (tx, rx): (Sender<bool>, Receiver<bool>) = mpsc::channel();

    loop {
        log(format!("Test {} Tasks", i).as_str());

        let (_, max_tokio_blocking_threads) = (num_cpus::get(), 512);
        let rt = runtime::Builder::new_multi_thread()
            .enable_all()
            .thread_stack_size(8 * 1024 * 1024)
            .worker_threads(i as usize)
            .max_blocking_threads(max_tokio_blocking_threads)
            .build();

        let rt = rt.unwrap();
        let thread = tx.clone();

        thread::spawn(move || {
            let sys = System::new();
            match sys.cpu_load_aggregate() {
                Ok(cpu) => {
                    sleep(Duration::from_secs(1));
                    let cpu = cpu.done().unwrap();
                    if cpu.idle * 100.0 < 20.0 || cpu.user * 100.0 > 80.0 {
                        thread.send(true).unwrap();
                    } else {
                        thread.send(false).unwrap();
                    }
                },
                Err(x) => println!("\nCPU load: error: {}", x)
            }
        });

        for _ in 0..i {
            rt.spawn(async {
                let _ = performance().await;
            });
        }

        sleep(Duration::from_secs(10));
        let result = rx.recv();
        rt.shutdown_background();
        if let Ok(result) = result {
            if result {
                i -= 500;
                break
            } else {
                i += 500;
            }
        }
    }
    i
}