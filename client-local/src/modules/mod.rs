use std::ops::Add;
use std::sync::{Arc, mpsc, Mutex};
use std::sync::mpsc::{Receiver, Sender};
use std::time::SystemTime;
use chrono::DateTime;
use rand::prelude::SliceRandom;
use tokio::runtime;
use tokio::runtime::Runtime;
use tokio::task::JoinHandle;
use crate::captcha::{captcha_loop};
use crate::data::values::{Billing, Data, Product, Shipping};
use crate::jig::{address, name, xyz};
use crate::log;
use crate::modules::session::Session;
use crate::request::proxy::ReqProxy;
use crate::settings::values::Settings;
use crate::websocket::Monitor;

mod kitheu;
pub(crate) mod shopify;
mod rates;
pub(crate) mod session;

#[derive(Clone)]
pub(crate) struct Task {
    pub module: i32,
    pub prod_id: String,
    pub store: String,
    pub size: String,
    pub proxy: ReqProxy,
    pub shipping: Shipping,
    pub name: String,
    pub image: String,
    pub price: f64,
    pub est_sell: f64,
    billing: Billing,
}

pub(crate) static mut RUNTIME: Option<Runtime> = None;

#[allow(dead_code)]
pub(crate) struct ThreadPool {
    pub(crate) active: bool,
    workers: Vec<Worker>,
    sender: Sender<Job>,
    session_workers: Vec<Worker>,
    pub(crate) session_sender: Sender<Data>,
}

#[derive(Clone)]
pub(crate) struct Job {
    pub(crate) data: Data,
    pub(crate) sku: String,
    pub(crate) size: String,
    pub(crate) store: String,
    pub(crate) module: i32,
    pub(crate) name: String,
    pub(crate) image: String,
    pub(crate) price: f64,
    pub(crate) est_sell: f64,
}

impl ThreadPool {
    pub(crate) fn new(s: Settings) -> ThreadPool {

        let (_, max_tokio_blocking_threads) = (num_cpus::get(), 512);
        let rt = runtime::Builder::new_multi_thread()
            .enable_all()
            .thread_stack_size(8 * 1024 * 1024)
            .worker_threads((s.task_max + 2) as usize)
            .max_blocking_threads(max_tokio_blocking_threads)
            .build();

        let rt = rt.unwrap();
        let (tx, rx): (Sender<Job>, Receiver<Job>) = mpsc::channel();
        let receiver = Arc::new(Mutex::new(rx));
        let mut workers = Vec::with_capacity(s.task_max as usize);

        for _ in 0..s.task_max {
            workers.push(Worker::new(Arc::clone(&receiver), &rt));
        }

        let (ts, rs): (Sender<Data>, Receiver<Data>) = mpsc::channel();
        let receiver = Arc::new(Mutex::new(rs));
        let mut s_workers = Vec::with_capacity(2);

        s_workers.push(Worker::new_session(Arc::clone(&receiver), &rt));
        s_workers.push(Worker::new_captcha(s.authorization, s.task_max, &rt));

        unsafe {
            RUNTIME = Some(rt);
        }

        ThreadPool {
            active: true,
            workers,
            sender: tx,
            session_workers: s_workers,
            session_sender: ts
        }
    }

    pub(crate) fn send(&mut self, d: Data, monitor: Monitor, module: i32) {
        let s = monitor.clone().skus.clone();
        let store = monitor.clone().link.split("/").nth(2).unwrap().to_string();
        for _ in 0..self.workers.len() {
            let mut rng = rand::thread_rng();
            let keys: Vec<String> = s.keys().map(|x| x.to_string()).collect();
            let vals: Vec<String> = s.values().map(|x| x.to_string()).collect();

            let key = keys.choose(&mut rng);
            if key.is_none() {
                continue;
            }
            let key = key.unwrap();

            let val = vals.choose(&mut rng);
            if val.is_none() {
                continue;
            }
            let val = val.unwrap();

            let j = Job {
                data: d.clone(),
                sku: val.to_string(),
                size: key.to_string(),
                store: store.to_string(),
                module: module.to_owned(),
                name: monitor.name.clone(),
                image: monitor.image.clone(),
                price: monitor.price.clone(),
                est_sell: monitor.est_sell.clone()
            };
            self.sender.send(j).unwrap();
        }
    }
}

#[allow(dead_code)]
struct Worker {
    thread: JoinHandle<()>,
}

impl Worker {
    fn new(receiver: Arc<Mutex<Receiver<Job>>>, rt: &Runtime) -> Worker {
        let thread = rt.spawn(async move {
            loop {
                let mut t = Task::new();
                let job = receiver.lock().unwrap().recv();
                if let Ok(job) = job {
                    log(format!("Got Product {} {}", job.size, job.sku).as_str());
                    t.parse_data(job.data.clone());
                    t.prod_id = job.sku.clone();
                    t.size = job.size.clone();
                    t.store = job.store.clone();
                    t.module = job.module.clone();
                    t.name = job.name.clone();
                    t.image = job.image.clone();
                    t.price = job.price.clone();
                    t.est_sell = job.est_sell.clone();

                    t.start().await;

                } else {
                    break;
                }
            }
        });
        Worker {
            thread,
        }
    }

    fn new_session(receiver: Arc<Mutex<Receiver<Data>>>, rt: &Runtime) -> Worker {
        let thread = rt.spawn(async move {
            let d = receiver.lock().unwrap().recv();
            if let Ok(mut d) = d {
                unsafe {
                    d.get_shipping_rates().await;
                    d.clone().gen_session(d.clone().session.tasks as usize).await;
                }
            }
        });
        Worker {
            thread,
        }
    }

    fn new_captcha(auth: String, task_max: i32, rt: &Runtime) -> Worker {
        let thread = rt.spawn(async move {
            captcha_loop(auth, task_max).await;
        });
        Worker {
            thread,
        }
    }
}

impl Task {
    pub(crate) fn new() -> Task {
        Task {
            prod_id: "".to_string(),
            store: "".to_string(),
            size: "".to_string(),
            module: 0,
            proxy: ReqProxy::get(),
            shipping: Shipping::new(),
            name: "".to_string(),
            image: "".to_string(),
            price: 0.0,
            est_sell: 0.0,
            billing: Billing::new(),
        }
    }

    pub(crate) fn parse_data(&mut self, d: Data) {
        let mut rng = rand::thread_rng();
        self.shipping = d.shipping.unwrap().choose(&mut rng).unwrap().clone();
        self.billing = d.billing.unwrap().choose(&mut rng).clone().unwrap().clone();

        self.shipping.last = name(self.shipping.last.clone());
        self.shipping.address1 = address(self.shipping.address1.clone());
        self.shipping.address2 = address(self.shipping.address2.clone());

        if !self.shipping.email.contains("@") {
            self.shipping.email = self.shipping.last.clone().add(xyz().as_str()).add("@").add(self.shipping.email.as_str());
        }
    }

    pub(crate) async fn start(self) {
        match self.module {
            0 => self.kith().await,
            1 => self.shopify().await,
            _ => return,
        }
    }

    pub(crate) async fn checked_out(self, s: Session) {
        Product {
            date: DateTime::from(SystemTime::now()),
            name: self.name.clone(),
            link: "https://".to_owned().add(self.store.clone().as_str()).add("/variants/").add(self.prod_id.as_str()).to_string(),
            image: self.image,
            store: self.store,
            size: self.size,
            price: self.price,
            est_sell: self.est_sell
        }.checkout(s.checkout.as_str()).await;
    }
}