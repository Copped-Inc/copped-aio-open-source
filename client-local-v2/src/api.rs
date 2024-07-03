use std::collections::HashMap;
use std::fs::File;
use std::{fs, io};
use std::ops::Add;
use std::process::{Command, exit};
use chrono::{DateTime, Local};
use reqwest::StatusCode;
use self_update::cargo_crate_version;
use tungstenite::handshake::client;
use crate::{cache, data, error, log, payments};

#[derive(serde::Serialize)]
#[derive(Debug)]
pub struct LoginReq {
    pub code: String,
}

#[derive(serde::Deserialize)]
#[derive(Debug)]
pub struct LoginRes {
    pub authorization: String,
}

pub async fn login(code: String) -> Result<LoginRes, &'static str> {
    let req = LoginReq { code };

    let api = Api::new_database(
        "instance",
        false,
        METHOD_POST,
        HashMap::new(),
        Some(reqwest::Body::from(serde_json::to_string(&req).unwrap())),
    );

    let r = api.request().await;
    if let Ok(res) = r {
        if res.status().is_success() {
            let body = res.text().await.unwrap();
            let l: LoginRes = serde_json::from_str(&body).unwrap();

            Ok(l)
        } else {
            Err("Failed to login")
        }
    } else {
        log!("{:}", r.as_ref().unwrap_err());
        Err("Failed to login")
    }
}

#[derive(serde::Deserialize)]
pub(crate) struct ProxyRes {
    pub proxies: Vec<ResProxy>,
}

#[derive(serde::Deserialize)]
pub(crate) struct ResProxy {
    pub ip: String,
    pub port: String,
    pub username: String,
    pub password: String,
}

pub(crate) async fn proxies() -> Result<ProxyRes, &'static str> {
    let api = Api::new_service(
        "proxies",
        true,
        METHOD_GET,
        HashMap::new(),
        None,
    );

    let r = api.request().await;
    if let Ok(res) = r {
        if res.status().is_success() {
            let body = res.text().await.unwrap();
            let l: ProxyRes = serde_json::from_str(&body).unwrap();

            Ok(l)
        } else {
            Err("Failed to get proxies")
        }
    } else {
        Err("Failed to get proxies")
    }
}

#[derive(serde::Deserialize)]
pub struct CaptchaRes {
    pub token: Vec<String>,
    pub expire: DateTime<Local>,
}

pub async fn gen_captcha(site: &'static str) -> Result<CaptchaRes, &'static str> {
    let api = Api::new_database(
        format!("captcha/{}", site).as_str(),
        true,
        METHOD_GET,
        HashMap::new(),
        None,
    );

    let r = api.request().await;
    if let Ok(res) = r {
        if res.status().is_success() {
            let body = res.text().await.unwrap();
            let c: CaptchaRes = serde_json::from_str(&body).unwrap();

            Ok(c)
        } else {
            Err("Failed to generate captcha")
        }
    } else {
        Err("Failed to generate captcha")
    }
}

pub async fn update() -> Result<(), &'static str> {
    let paths = fs::read_dir("./").unwrap();
    for path in paths {
        if let Ok(p) = path {
            if p.path().is_dir() && p.file_name().to_str().unwrap().contains("self_update") {
                fs::remove_dir_all(p.path()).unwrap();
            }
        }
    }

    let current_exe = std::env::current_exe();
    if current_exe.is_err() {
        return Err("Failed to get current exe");
    }

    let name = current_exe.as_ref().unwrap().file_name().unwrap().to_str().unwrap();
    log!("Checking for {}", name);

    let mut h = HashMap::new();
    h.insert("version".to_string(), cargo_crate_version!().to_string());

    let api = Api::new_database(
        format!("instance/update/{}", name).as_str(),
        true,
        METHOD_GET,
        h,
        None,
    );

    let r = api.request().await;
    if let Ok(res) = r {
        if res.status() == StatusCode::OK {
            Ok(())
        } else if res.status() == StatusCode::CREATED {
            log!("Update available, downloading...");
            let current_dir = std::env::current_dir();
            if current_dir.is_err() {
                return Err("Failed to get current dir");
            }

            let tmp_dir = tempfile::Builder::new()
                .prefix("self_update")
                .tempdir_in(current_dir.unwrap());

            if tmp_dir.is_err() {
                return Err("Failed to create temp dir");
            }

            let tmp_dir = tmp_dir.unwrap();

            let tmp_tarball_path = tmp_dir.path().join(name);
            let tmp_tarball = File::create(&tmp_tarball_path);
            if tmp_tarball.is_err() {
                return Err("Failed to open temp tarball");
            }

            io::copy(&mut res.bytes().await.unwrap().as_ref(), &mut tmp_tarball.unwrap()).unwrap();

            let tmp_file = tmp_dir.path().join("replacement_tmp");
            let dest = self_update::Move::from_source(&tmp_tarball_path)
                .replace_using_temp(&tmp_file)
                .to_dest(&current_exe.as_ref().unwrap());

            if dest.is_err() {
                return Err("Failed to move tarball");
            }

            if !cfg!(target_os = "windows") {
                log!("Setting permissions...");
                let cmd = Command::new("chmod")
                    .arg("+x")
                    .arg(&current_exe.as_ref().unwrap())
                    .spawn();

                if cmd.is_err() {
                    return Err("Failed to run chmod");
                }

                let output = cmd.unwrap().wait_with_output();
                if output.is_err() {
                    return Err("Failed to get output of chmod");
                }
            }

            if cfg!(target_os = "windows") {
                let mut cmd = Command::new(&current_exe.unwrap());
                cmd.spawn().unwrap();

                exit(0)
            }

            exit(1)
        } else {
            error!("Failed to update: {}", res.status().clone());
            Err("Failed to update")
        }
    } else {
        error!("Failed to update: {}", r.as_ref().unwrap_err().to_string());
        Err("Failed to update")
    }
}

pub async fn checkout(product: data::Product) -> Result<(), &'static str> {
    log!("Checking out {}", serde_json::to_string(&product).unwrap());
    let api = Api::new_database(
        "data/checkout",
        true,
        METHOD_POST,
        HashMap::new(),
        Some(reqwest::Body::from(serde_json::to_string(&product).unwrap())),
    );

    let r = api.request().await;
    if let Ok(res) = r {
        if res.status() == StatusCode::OK {
            Ok(())
        } else {
            Err("Failed to send checkout")
        }
    } else {
        Err("Failed to send checkout")
    }
}

pub async fn payment(p: payments::Payment) -> Result<(), &'static str> {
    log!("Sending payment {}", serde_json::to_string(&p).unwrap());
    let api = Api::new_database(
        "instance/payments",
        true,
        METHOD_POST,
        HashMap::new(),
        Some(reqwest::Body::from(serde_json::to_string(&p).unwrap())),
    );

    let r = api.request().await;
    if let Ok(res) = r {
        if res.status() == StatusCode::OK {
            Ok(())
        } else {
            Err("Failed to send payment")
        }
    } else {
        Err("Failed to send payment")
    }
}

pub fn websocket() -> client::Request {
    let s = cache::settings();
    client::Request::builder()
        .uri(API_URL.clone().replace("http", "ws").add("websocket"))
        .method("GET")
        .header("Host", (API_URL.clone().split("://").nth(1).unwrap()).split("/").nth(0).unwrap())
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
        .header("Reconnect", format!("{}", {
            if cache::connected() {
                "1"
            } else {
                ""
            }
        }))
        .header("Cookie", format!("authorization={}", s.authorization))
        .body(())
        .unwrap()
}

pub const METHOD_POST: i8 = 0;
pub const METHOD_GET: i8 = 1;
pub const METHOD_PUT: i8 = 2;
pub const METHOD_DELETE: i8 = 3;

const API_URL: &str = "https://database.copped-inc.com/";
/*const API_URL: &str = "http://localhost:91/";*/
const SERVICE_URL: &str = "https://service.copped-inc.com/";

struct Api {
    url: String,
    path: String,
    auth: bool,
    method: i8,
    headers: HashMap<String, String>,
    body: Option<reqwest::Body>,
}

impl Api {
    fn new_database(path: &str, auth: bool, method: i8, headers: HashMap<String, String>, body: Option<reqwest::Body>) -> Self {
        Self {
            url: API_URL.to_string(),
            path: path.to_string(),
            auth,
            method,
            headers,
            body,
        }
    }
    fn new_service(path: &str, auth: bool, method: i8, headers: HashMap<String, String>, body: Option<reqwest::Body>) -> Self {
        Self {
            url: SERVICE_URL.to_string(),
            path: path.to_string(),
            auth,
            method,
            headers,
            body,
        }
    }

    async fn request(self) -> Result<reqwest::Response, reqwest::Error> {
        let client = reqwest::Client::builder().build()?;

        let mut request = match self.method {
            METHOD_POST => client.post(self.url.to_string().add(self.path.as_str())),
            METHOD_GET => client.get(self.url.to_string().add(self.path.as_str())),
            _ => client.post(self.url.to_string().add(self.path.as_str())),
        };

        for (key, value) in self.headers {
            request = request.header(key, value);
        }

        if self.auth {
            request = request.header("cookie", "authorization=".to_string().add(cache::settings().authorization.as_str()));
        }

        if let Some(body) = self.body {
            request = request.body(body);
        }

        request.send().await
    }
}