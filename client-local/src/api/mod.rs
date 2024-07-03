use std::fs::File;
use std::io::Write;
use std::ops::Add;
use std::process::{Command, exit};
use chrono::{DateTime, Local};
use hyper::{Body, body, Client, Method, Request};
use tokio::io::{self};
use serde::{Serialize, Deserialize};
use hyper_tls::HttpsConnector;
use crate::data::values::Product;

// pub(crate) static DATABASE_URL: &str = "http://localhost:91/";
pub(crate) static DATABASE_URL: &str = "https://database.copped-inc.com/";
// pub(crate) static SERVICE_URL: &str = "http://localhost:93/";
pub(crate) static SERVICE_URL: &str = "https://service.copped-inc.com/";
pub(crate) static mut AUTH: String = String::new();

#[derive(Deserialize)]
#[derive(Debug)]
pub(crate) struct LoginRes {
    pub authorization: String,
}

pub(crate) async fn login(code: String) -> Result<LoginRes, Box<dyn std::error::Error + Send + Sync>> {
    let https = HttpsConnector::new();
    let client = Client::builder()
        .build::<_, Body>(https);

    let req = Request::builder()
        .method(Method::POST)
        .uri(DATABASE_URL.clone().to_string().add("instance"))
        .header("code", code.clone())
        .body(Body::empty())?;

    let res = client.request(req).await;
    if let Ok(mut res) = res {
        if res.status().is_success() {
            let body = hyper::body::to_bytes(res.body_mut()).await.unwrap();
            let l: LoginRes = serde_json::from_slice(&body).unwrap();

            Ok(l)
        } else {
            Err(format!("{}", res.status()).into())
        }
    } else {
        Err(Box::new(io::Error::new(io::ErrorKind::Other, "Failed to login")))
    }
}

#[derive(Serialize)]
pub(crate) struct PerformanceReq {
    pub performance: String,
}

#[derive(Deserialize)]
pub(crate) struct PerformanceRes {
    pub performance: String,
}

pub(crate) async fn performance() -> Result<String, Box<dyn std::error::Error + Send + Sync>> {
    let https = HttpsConnector::new();
    let client = Client::builder()
        .build::<_, Body>(https);

    let performance = PerformanceReq {
        performance: "Test".to_string(),
    };

    let uri = DATABASE_URL.clone().to_string().add("instance/performance");
    let req = Request::builder()
        .method(Method::POST)
        .uri(uri)
        .body(Body::from(serde_json::to_string(&performance).unwrap()))?;

    let res = client.request(req).await;
    if let Ok(mut res) = res {
        return if res.status().is_success() {
            let body = hyper::body::to_bytes(res.body_mut()).await.unwrap();
            let l: PerformanceRes = serde_json::from_slice(&body).unwrap();

            Ok(l.performance)
        } else {
            Err(format!("{}", res.status()).into())
        }
    }
    Err(Box::new(io::Error::new(io::ErrorKind::Other, res.unwrap_err())))
}

pub(crate) async fn update(auth: String) -> Result<(), Box<dyn std::error::Error + Send + Sync>> {
    unsafe {
        AUTH = auth.clone();
    }

    let https = HttpsConnector::new();
    let client = Client::builder()
        .build::<_, Body>(https);

    let performance = PerformanceReq {
        performance: "Test".to_string(),
    };

    let uri = DATABASE_URL.clone().to_string().add("instance/update");
    let req = Request::builder()
        .method(Method::GET)
        .uri(uri)
        .header("version", "0.0.48")
        .header("cookie", "authorization=".to_string() + &auth)
        .body(Body::from(serde_json::to_string(&performance).unwrap()))?;

    let res = client.request(req).await;
    if let Ok(res) = res {
        if res.status() == hyper::StatusCode::OK {
            return Ok(());
        } else if res.status() == hyper::StatusCode::CREATED {
            let mut file = File::create("copped-aio-update.exe").unwrap();
            file.write_all(&body::to_bytes(res.into_body()).await.unwrap()).unwrap();

            let name = std::env::current_exe()
                .ok()
                .and_then(|pb| pb.file_name().map(|s| s.to_os_string()))
                .and_then(|s| s.into_string().ok()).unwrap();

            let mut file = File::create("update.bat").unwrap();
            file.write_all(format!("@echo off\ntimeout /t 2\ndel /f \"{}\"\nrename \"copped-aio-update.exe\" \"{}\"\nstart ./\"{}\"\nexit 0", name.clone(), name.clone(), name.clone()).as_bytes()).unwrap();
            let _ = Command::new("cmd.exe")
                .args(&["/C", "start", "update.bat"])
                .status()
                .expect("failed to execute process");

            exit(0);
        } else {
            return Err(format!("{}", res.status()).into());
        }
    }
    Err(Box::new(io::Error::new(io::ErrorKind::Other, res.unwrap_err())))
}

#[derive(Deserialize)]
pub(crate) struct CaptchaRes {
    pub token: Vec<String>,
    pub expire: DateTime<Local>,
}

pub(crate) async fn captcha(auth: String) -> Result<CaptchaRes, Box<dyn std::error::Error + Send + Sync>> {
    let https = HttpsConnector::new();
    let client = Client::builder()
        .build::<_, Body>(https);

    let uri = DATABASE_URL.clone().to_string().add("captcha");
    let req = Request::builder()
        .method(Method::GET)
        .uri(uri)
        .header("cookie", "authorization=".to_string() + &auth)
        .body(Body::empty())?;

    let res = client.request(req).await;

    if let Ok(mut res) = res {
        if res.status() == hyper::StatusCode::OK {
            let body = hyper::body::to_bytes(res.body_mut()).await.unwrap();
            let l: CaptchaRes = serde_json::from_slice(&body).unwrap();
            return Ok(l);
        }
        Err(Box::new(io::Error::new(io::ErrorKind::Other, "Unhandled status code")))
    } else {
        Err(Box::new(io::Error::new(io::ErrorKind::Other, res.unwrap_err())))
    }
}

pub(crate) async unsafe fn checkout(p: Product) -> Result<(), Box<dyn std::error::Error + Send + Sync>> {
    let https = HttpsConnector::new();
    let client = Client::builder()
        .build::<_, Body>(https);

    let uri = DATABASE_URL.clone().to_string().add("data/checkout");
    let req = Request::builder()
        .method(Method::POST)
        .uri(uri)
        .header("cookie", "authorization=".to_string() + &AUTH.clone())
        .body(Body::from(serde_json::to_string(&p).unwrap()))?;

    let res = client.request(req).await;
    if let Ok(_) = res {
        return Ok(());
    }
    Err(Box::new(io::Error::new(io::ErrorKind::Other, res.unwrap_err())))
}

#[derive(Deserialize)]
pub(crate) struct ProxyRes {
    pub proxies: Vec<ResProxy>,
}

#[derive(Deserialize)]
pub(crate) struct ResProxy {
    pub ip: String,
    pub port: String,
    pub username: String,
    pub password: String,
}

pub(crate) async unsafe fn proxies(auth: String) -> Result<ProxyRes, Box<dyn std::error::Error + Send + Sync>> {
    let https = HttpsConnector::new();
    let client = Client::builder()
        .build::<_, Body>(https);

    let uri = SERVICE_URL.clone().to_string().add("proxies");
    let req = Request::builder()
        .method(Method::GET)
        .uri(uri)
        .header("cookie", "authorization=".to_string() + &auth)
        .body(Body::empty())?;

    let res = client.request(req).await;
    if let Ok(mut res) = res {
        if res.status() == hyper::StatusCode::OK {
            let body = hyper::body::to_bytes(res.body_mut()).await.unwrap();
            let l: ProxyRes = serde_json::from_slice(&body).unwrap();
            return Ok(l);
        }
        Err(Box::new(io::Error::new(io::ErrorKind::Other, "Unhandled status code")))
    } else {
        Err(Box::new(io::Error::new(io::ErrorKind::Other, res.unwrap_err())))
    }
}