use std::collections::HashMap;
use std::{fs, thread};
use chrono::{DateTime, Local};
use rand::Rng;
use crate::{api, data, jig, log, settings};
use crate::cache::Cache;
use crate::request::cookie::Cookie;
use crate::request::proxy::Proxy;
use crate::threads::handle;

#[derive(serde::Deserialize, serde::Serialize, Clone, PartialEq)]
pub enum TaskType {
    /// Need preloaded task
    Account,
    /// Can use preloaded task
    Product,
    /// Never use preloaded task
    OneTime,
}

impl Default for TaskType {
    fn default() -> Self {
        TaskType::Product
    }
}

#[derive(serde::Deserialize, serde::Serialize, Default, Clone)]
pub struct Sessions {
    pub module: String,
    pub tasks: Vec<Task>,
    pub task_type: TaskType,
    pub fallback: i32,
    pub delete_values: Vec<String>,
    pub task_multiply: i32,
    pub checkout_value: String,

    #[serde(skip_serializing, skip_deserializing)]
    preload: Vec<fn(&mut Task) -> Result<(), &'static str>>,
    #[serde(skip_serializing, skip_deserializing)]
    run: Vec<fn(&mut Task) -> Result<(), &'static str>>,
}

#[derive(serde::Deserialize, serde::Serialize, Clone, Default)]
pub struct Task {
    pub state: i32,
    pub task_type: TaskType,
    pub module: String,
    pub checkout_value: String,
    pub proxy: Proxy,
    pub shipping: data::Shipping,
    pub billing: data::Billing,
    pub product_id: String,
    pub size: String,
    pub link: String,
    pub data: HashMap<String, String>,
    pub cookies: Vec<Cookie>,
}

#[derive(serde::Deserialize, serde::Serialize)]
pub struct Websocket {
    pub data: Product,
}

#[derive(serde::Deserialize, serde::Serialize)]
pub struct Product {
    pub name: String,
    pub sku: String,
    pub skus: HashMap<String, String>,
    pub date: DateTime<Local>,
    pub link: String,
    pub image: String,
    pub price: f64,
}

pub fn websocket(s: String) {
    let d: Websocket = serde_json::from_str(&s.as_str()).unwrap();

    if d.data.name.to_lowercase().contains("test") && !Cache::dev() {
        log!("Skipping Dev Product");
        return;
    }

    if d.data.link.contains("kith") && Cache::kith_eu() {
        log!("Kith EU new Product");
        if d.data.skus.len() == 0 {
            log!("Skipping Product with no Skus");
            return;
        }

        let sizes: Vec<String> = d.data.skus.keys().map(|x| x.to_string()).collect();
        let skus: Vec<String> = d.data.skus.values().map(|x| x.to_string()).collect();

        let mut rng = rand::thread_rng();
        let i = rng.gen_range(0..skus.len());
        Sessions::run(skus[i].clone(), sizes[i].clone(), d.data.link, "kith_eu".to_string());
    } else if Cache::queue_it() {
        log!("Queue-it new Product");
        Sessions::run("".to_string(), "".to_string(), d.data.link, "queue_it".to_string());
    }
}

#[allow(dead_code)]
impl Sessions {
    pub fn new(module: &'static str) -> Self {
        Sessions {
            module: module.to_string(),
            task_multiply: 1,
            ..Default::default()
        }
    }

    pub fn task_type(&mut self, task_type: TaskType) -> &mut Self {
        self.task_type = task_type;
        self
    }

    pub fn fallback(&mut self, fallback: i32) -> &mut Self {
        self.fallback = fallback;
        self
    }

    pub fn delete_values(&mut self, delete_values: Vec<&'static str>) -> &mut Self {
        self.delete_values = delete_values.iter().map(|x| x.to_string()).collect();
        self
    }

    pub fn task_multiply(&mut self, multiply: i32) -> &mut Self {
        self.task_multiply = multiply;
        self
    }

    pub fn checkout_value(&mut self, checkout_value: &'static str) -> &mut Self {
        self.checkout_value = checkout_value.to_string();
        self
    }

    pub fn add_preload(&mut self, f: Vec<fn(&mut Task) -> Result<(), &'static str>>) -> &mut Self {
        self.preload = f;
        self
    }

    pub fn add_run(&mut self, f: Vec<fn(&mut Task) -> Result<(), &'static str>>) -> &mut Self {
        self.run = f;
        self
    }

    pub fn build(&mut self) {
        let p = settings::path(&self.module);
        if p.exists() {
            let d = fs::read_to_string(p).expect("Failed to read settings file");
            let s = serde_json::from_str(&d);
            if s.is_ok() {
                let s: Sessions = s.unwrap();
                self.tasks = s.tasks;
            }
        }

        if self.task_type == TaskType::Account && self.tasks.len() <= (Cache::settings().task_max * self.task_multiply) as usize {
            self.preload();
        }

        self.save();
    }

    pub fn save(&self) {
        if self.task_type != TaskType::OneTime {
            let p = settings::path(&self.module);
            let file = fs::File::create(p).expect("Failed to create settings file");
            serde_json::to_writer_pretty(file, self).expect("Failed to write settings file");
        }

        Cache::update_session(self.clone());
    }

    fn preload(&mut self) {
        for _ in self.tasks.len()..Cache::settings().task_max as usize {
            let task = Task::get(self.module.clone(), self.task_type.clone(), self.checkout_value.clone());
            task.run_thread(self.preload.clone());
        }
    }

    pub fn run(pid: String, size: String, link: String, module: String) {
        let mut sessions = Cache::sessions();
        for s in sessions.iter_mut() {
            if s.module == module {
                s.start(pid.clone(), size.clone(), link.clone());
            }
        }
    }

    fn start(&mut self, pid: String, size: String, link: String) {
        let mut max = Cache::settings().task_max * self.task_multiply;
        log!("Running {} Tasks", max);
        match self.task_type {
            TaskType::Account => {
                for task in self.tasks.iter_mut() {
                    task.product_id = pid.clone();
                    task.size = size.clone();
                    task.link = link.clone();
                    task.run_thread(self.run.clone());
                }
                self.tasks.clear();
            },
            TaskType::Product => {
                let tasks = self.tasks.clone();
                for i in (0..tasks.len()).rev() {
                    if tasks[i].product_id == pid && tasks[i].size == size && max > 0 {
                        tasks[i].run_thread(self.run.clone());
                        self.tasks.remove(i);
                        max -= 1;
                    }
                }

                for _ in 0..max {
                    let mut task = Task::get(self.module.clone(), self.task_type.clone(), self.checkout_value.clone());
                    task.product_id = pid.clone();
                    task.size = size.clone();
                    task.link = link.clone();
                    task.run_thread(self.run.clone());
                }
            },
            TaskType::OneTime => {
                for _ in 0..max {
                    let mut task = Task::get(self.module.clone(), self.task_type.clone(), self.checkout_value.clone());
                    task.link = link.clone();
                    task.run_thread(self.run.clone());
                }
            },
        }
        self.save();
    }
}

#[allow(dead_code)]
impl Task {
    fn get(module: String, task_type: TaskType, checkout_value: String) -> Self {
        let (shipping, billing) = jig::get();
        Task {
            proxy: Proxy::get(),
            task_type,
            module,
            checkout_value,
            shipping,
            billing,
            ..Default::default()
        }
    }

    fn run_thread(&self, fs: Vec<fn(&mut Task) -> Result<(), &'static str>>) {
        let mut task = self.clone();
        let functions = fs.clone();
        thread::spawn(move || {
            log!("Starting: {} {} {}", task.link, task.product_id, task.size);
            if task.state >= functions.len() as i32 {
                log!("{}: {}", "Invalid state", task.product_id);
                return;
            }

            for i in task.state as usize..functions.len() {
                let r = functions[i](&mut task);
                if r.is_ok() { task.state += 1; }
                else {
                    log!("{}: {} {} {}", r.unwrap_err(), task.link, task.product_id, task.size);
                    break;
                }
            }

            if task.state == functions.len() as i32 {
                handle().block_on(task.checkout());
                log!("Success: {} {}", task.link, task.product_id);
                if task.task_type == TaskType::Account {
                    task.state = 0;
                    task.save();
                }
                return;
            }

            if task.state > 0 { task.save(); }
        });
    }

    fn save(&mut self) {
        if self.task_type == TaskType::OneTime {
            return;
        }

        Cache::add_task(self.clone());
        let s = Cache::sessions();
        for session in s.iter() {
            if session.module == self.module {
                session.save();
                break;
            }
        }
    }

    async fn checkout(&self) {
        let mut checkout = String::new();
        if self.checkout_value != "" {
            checkout = self.data.get(&self.checkout_value).unwrap().to_owned();
        }

        let _ = api::checkout(api::CheckoutReq {
            link: self.link.clone(),
            site: self.module.clone(),
            size: self.size.clone(),
            checkout,
        }).await;
    }
}

pub trait Helper {
    fn get_string(&self, f: &'static str) -> String;
    fn add_string(&mut self, f: &'static str, v: String);
    fn add_str(&mut self, f: &'static str, v: &'static str);
}

impl Helper for HashMap<String, String> {
    fn get_string(&self, f: &'static str) -> String {
        match self.get(f) {
            Some(v) => v.to_string(),
            None => String::new(),
        }
    }

    fn add_string(&mut self, f: &'static str, v: String) {
        self.insert(f.to_string(), v);
    }

    fn add_str(&mut self, f: &'static str, v: &'static str) {
        self.insert(f.to_string(), v.to_string());
    }
}
