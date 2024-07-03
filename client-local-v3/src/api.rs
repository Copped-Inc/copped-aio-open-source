use std::collections::HashMap;
use std::fs::File;
use std::io;
use std::ops::Add;
use std::process::{Command, exit};
use chrono::{DateTime, Local};
use reqwest::StatusCode;
use self_update;
use self_update::cargo_crate_version;
use crate::cache::Cache;
use crate::log;
use crate::request::proxy;

#[derive(serde::Serialize, Debug)]
pub struct LoginReq {
    pub code: String,
}

#[derive(serde::Deserialize, Debug)]
pub struct LoginRes {
    pub authorization: String,
}

pub async fn login(code: String) -> Result<LoginRes, &'static str> {
    let r = Api::new_database(
        "instance",
        false,
        Method::Post,
        HashMap::new(),
        body(&LoginReq {
            code,
        }),
    ).await;

    if r.is_err() {
        log!("Error: {}", r.as_ref().unwrap_err());
        return Err("Failed to login");
    }

    let l: LoginRes = r.unwrap().json().await.unwrap();
    Ok(l)
}

#[derive(serde::Deserialize)]
pub struct ProxyRes {
    pub proxies: Vec<proxy::Proxy>,
}

pub async fn proxies() -> Result<ProxyRes, &'static str> {
    let r = Api::new_serivce(
        "proxies",
        true,
        Method::Get,
        HashMap::new(),
        None,
    ).await;

    if r.is_err() {
        log!("Error: {}", r.as_ref().unwrap_err());
        return Err("Failed to get proxies");
    }

    let p: ProxyRes = r.unwrap().json().await.unwrap();
    Ok(p)
}

pub async fn update() -> Result<(), &'static str> {
    let current_exe = std::env::current_exe();
    if current_exe.is_err() {
        return Err("Failed to get current exe");
    }

    let name = current_exe.as_ref().unwrap().file_name().unwrap().to_str().unwrap();
    log!("Checking for {}", name);

    let mut h = HashMap::new();
    h.insert("version".to_string(), cargo_crate_version!().to_string());

    let r = Api::new_database(
        format!("instance/update/{}", name).as_str(),
        true,
        Method::Get,
        h,
        None,
    ).await;

    if let Ok(res) = r {
        if res.status() == StatusCode::OK {
            return Ok(());
        } else if res.status() != StatusCode::CREATED {
            log!("Failed to update: {}", res.status().clone());
            return Err("Failed to update");
        }

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

        if cfg!(target_os = "windows") {
            let mut cmd = Command::new(&current_exe.unwrap());
            cmd.spawn().unwrap();

            exit(0)
        } else {
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

            exit(1)
        }
    }

    log!("Failed to update: {}", r.as_ref().unwrap_err().to_string());
    Err("Failed to update")
}

#[derive(serde::Deserialize)]
pub struct CaptchaRes {
    pub token: Vec<String>,
    pub expire: DateTime<Local>,
}

pub async fn gen_captcha(site: &'static str) -> Result<CaptchaRes, &'static str> {
    let r = Api::new_database(
        format!("captcha/{}", site).as_str(),
        true,
        Method::Get,
        HashMap::new(),
        None,
    ).await;

    if r.is_err() {
        log!("Error: {}", r.as_ref().unwrap_err());
        return Err("Failed to get captcha");
    }

    let p = r.unwrap().json().await;
    if p.is_err() {
        return Err("Failed to get captcha");
    }
    Ok(p.unwrap())
}

#[derive(serde::Serialize)]
pub struct CheckoutReq {
    pub link: String,
    pub site: String,
    pub size: String,
    pub checkout: String,
}

pub async fn checkout(req_body: CheckoutReq) -> Result<(), &'static str> {
    let r = Api::new_database(
        "data/checkout",
        true,
        Method::Post,
        HashMap::new(),
        body(&req_body),
    ).await;

    if r.is_err() {
        log!("Error: {}", r.as_ref().unwrap_err());
        return Err("Failed send to checkout");
    }

    Ok(())
}

#[allow(dead_code)]
#[derive(PartialEq)]
pub enum Method {
    Post,
    Get,
    Put,
    Delete,
}

struct Api {
    url: String,
    path: String,
    auth: bool,
    method: Method,
    headers: HashMap<String, String>,
    body: Option<reqwest::Body>,
}

fn body<T>(value: &T) -> Option<reqwest::Body>
where T: serde::Serialize {
    Some(reqwest::Body::from(serde_json::to_string(value).unwrap()))
}

impl Api {
    async fn new(url: String, path: &str, auth: bool, method: Method, headers: HashMap<String, String>, body: Option<reqwest::Body>) -> Result<reqwest::Response, reqwest::Error> {
        Api {
            url: url.to_string(),
            path: path.to_string(),
            auth,
            method,
            headers,
            body,
        }.request().await
    }

    async fn new_database(path: &str, auth: bool, method: Method, headers: HashMap<String, String>, body: Option<reqwest::Body>) -> Result<reqwest::Response, reqwest::Error> {
        Api::new("https://database.copped-inc.com/".to_string(), path, auth, method, headers, body).await
    }

    async fn new_serivce(path: &str, auth: bool, method: Method, headers: HashMap<String, String>, body: Option<reqwest::Body>) -> Result<reqwest::Response, reqwest::Error> {
        Api::new("https://service.copped-inc.com/".to_string(), path, auth, method, headers, body).await
    }

    async fn request(self) -> Result<reqwest::Response, reqwest::Error> {
        let client = reqwest::Client::new();
        let mut request = match self.method {
            Method::Get => client.get(self.url + &self.path),
            Method::Post => client.post(self.url + &self.path),
            Method::Put => client.put(self.url + &self.path),
            Method::Delete => client.delete(self.url + &self.path),
        };

        for (key, value) in self.headers {
            request = request.header(key, value);
        }

        if self.auth {
            request = request.header("cookie", "authorization=".to_string().add(Cache::settings().authorization.as_str()));
        }

        if let Some(body) = self.body {
            request = request.body(body);
        }

        request.send().await
    }
}
