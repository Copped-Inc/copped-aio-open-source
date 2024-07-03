use std::ops::Add;
use std::time;
use chrono::{DateTime, Utc};
use regex::Regex;
use reqwest::StatusCode;
use serde::{Deserialize, Serialize};
use crate::request::expected::Expected;
use crate::{log, session};
use crate::session::Helper;

pub fn queue_it(t: &mut session::Task) -> Result<(), &'static str> {
    let req = t.request()
        .url(t.link.clone())
        .tls()
        .expected(Expected::status(StatusCode::FOUND));

    let res = req.get();
    if res.is_err() {
        return Err(res.err().unwrap());
    }

    if let Some(l) = res.unwrap().header.get("Location") {
        t.data.add_string("queueit", l.to_str().unwrap().to_string());
        return Ok(());
    }

    Err("Failed to get Location header")
}

pub fn event_id(t: &mut session::Task) -> Result<(), &'static str> {
    let req = t.request()
        .url(t.data.get_string("queueit"))
        .tls()
        .expected(Expected::multiple_status(vec![StatusCode::OK, StatusCode::FOUND]));

    let res = req.get();
    if res.is_err() {
        return Err(res.err().unwrap());
    }

    let res = res.unwrap();
    if res.status == StatusCode::FOUND {
        return Err("No queue live");
    }

    let re_event_id = Regex::new(r#"eventId:( |)'(.*?)',"#).unwrap();
    let re_layout = Regex::new(r#"layout:( |)'(.*?)',"#).unwrap();
    let re_target_url = Regex::new(r#"targetUrl:( |)decodeURIComponent\('(.*?)'\),"#).unwrap();
    let re_costumer_id = Regex::new(r#"customerId:( |)'(.*?)',"#).unwrap();
    let re_layout_version = Regex::new(r#"layoutVersion:( |)(.*?),"#).unwrap();

    let body = res.body.clone();
    let event_id = re_event_id.captures(body.as_str());
    let layout = re_layout.captures(body.as_str());
    let target_url = re_target_url.captures(body.as_str());
    let costumer_id = re_costumer_id.captures(body.as_str());
    let layout_version = re_layout_version.captures(body.as_str());
    if event_id.is_none() || layout.is_none() || target_url.is_none() || costumer_id.is_none() || layout_version.is_none() {
        return Err("Failed to get event_id || layout || target_url || costumer_id || layout_version");
    }

    t.data.add_string("event_id", event_id.unwrap().get(2).unwrap().as_str().to_string());
    t.data.add_string("layout", layout.unwrap().get(2).unwrap().as_str().to_string());
    t.data.add_string("target_url", urlencoding::decode(target_url.unwrap().get(2).unwrap().as_str()).expect("UTF-8").to_string());
    t.data.add_string("costumer_id", costumer_id.unwrap().get(2).unwrap().as_str().to_string());
    t.data.add_string("layout_version", layout_version.unwrap().get(2).unwrap().as_str().to_string());
    t.data.add_string("url", "https://".to_string().add(&t.data.get_string("costumer_id")).add(".queue-it.net/spa-api/queue/").add(&t.data.get_string("costumer_id")).add("/").add(&t.data.get_string("event_id")).add("/"));
    t.cookies = res.cookies;
    Ok(())
}

#[derive(Serialize)]
#[serde(rename_all = "camelCase")]
struct EnqueueReq {
    pub layout_name: String,
    pub custom_url_params: String,
    pub target_url: String,
    #[serde(rename = "Referrer")]
    pub referrer: String,
}

#[derive(Deserialize)]
#[serde(rename_all = "camelCase")]
struct EnqueueRes {
    pub queue_id: String,
}

pub fn enqueue(t: &mut session::Task) -> Result<(), &'static str> {
    let body = EnqueueReq {
        layout_name: t.data.get_string("layout"),
        custom_url_params: "".to_string(),
        target_url: t.data.get_string("target_url"),
        referrer: "https://".to_string().add(&t.data.get_string("costumer_id")).add(".queue-it.net/"),
    };

    let url = t.data.get_string("url").add("enqueue?cid=de-DE");

    let req = t.request()
        .url(url)
        .body_json(&body)
        .tls()
        .expected(Expected::status(StatusCode::OK));

    let res = req.post();
    if res.is_err() {
        return Err(res.err().unwrap());
    }

    let res = res.unwrap();
    let l: EnqueueRes = res.json().unwrap();
    t.data.add_string("queue_id", l.queue_id);
    t.cookies = res.cookies;
    Ok(())
}

#[derive(Serialize)]
#[serde(rename_all = "camelCase")]
pub struct StatusReq {
    pub target_url: String,
    pub custom_url_params: String,
    pub layout_version: i64,
    pub layout_name: String,
    pub is_client_reday_to_redirect: bool,
    pub is_before_or_idle: bool,
}

#[derive(Deserialize, Clone)]
#[serde(rename_all = "camelCase")]
pub struct StatusRes {
    pub redirect_url: Option<String>,
    pub is_before_or_idle: Option<bool>,
    pub ticket: Option<Ticket>
}

#[derive(Deserialize, Clone)]
#[serde(rename_all = "camelCase")]
pub struct Ticket {
    #[serde(rename = "eventStartTimeUTC")]
    pub event_start_time_utc: DateTime<Utc>,
    #[serde(rename = "expectedServiceTimeUTC")]
    pub expected_service_time_utc: Option<DateTime<Utc>>,
}

pub fn status(t: &mut session::Task) -> Result<(), &'static str> {
    if t.data.get_string("timeout").is_empty() {
        t.data.add_string("timeout", Utc::now().to_string());
    }

    if Utc::now().signed_duration_since(t.data.get_string("timeout").parse::<DateTime<Utc>>().unwrap()).num_minutes() > 60 {
        return Err("Timeout exceeded");
    }

    let body = StatusReq {
        target_url: t.data.get_string("target_url"),
        custom_url_params: "".to_string(),
        layout_version: t.data.get_string("layout_version").parse::<i64>().unwrap(),
        layout_name: t.data.get_string("layout"),
        is_client_reday_to_redirect: true,
        is_before_or_idle: false,
    };

    let url = t.data.get_string("url").add(&t.data.get_string("queue_id")).add("/status");

    let req = t.request()
        .url(url)
        .body_json(&body)
        .tls()
        .expected(Expected::status(StatusCode::OK));

    let res = req.post();
    if res.is_err() {
        return Err(res.err().unwrap());
    }

    let res = res.unwrap();
    let l: StatusRes = res.json().unwrap();
    if let Some(u) = l.redirect_url {
        log!("Redirecting to {}", u);
        t.data.add_string("redirect_url", u.to_string());
        return Ok(());
    }

    let before = l.is_before_or_idle.unwrap();
    let ticket = l.ticket.unwrap();
    if before {
        log!("Waiting for event start");
        let start_time = ticket.event_start_time_utc;
        let now = Utc::now();
        let diff = start_time.signed_duration_since(now);
        if diff.num_milliseconds() < 1000 {
            spin_sleep::sleep(time::Duration::from_millis(1000));
        } else {
            spin_sleep::sleep(diff.to_std().unwrap());
        }
        return status(t);
    }

    if ticket.expected_service_time_utc.is_none() {
        spin_sleep::sleep(time::Duration::from_millis(5000));
        return status(t);
    }

    let expected = ticket.expected_service_time_utc.unwrap();
    let now = Utc::now();
    let diff = expected.signed_duration_since(now);
    if diff.num_milliseconds() < 60000 {
        log!("Waiting before passing");
        spin_sleep::sleep(time::Duration::from_millis(60000));
    } else if diff.num_minutes() > 60 {
        log!("Stop task because timeout exceeded");
        return Err("Timeout exceeded");
    } else {
        log!("Waiting half of the time expected to pass {}", expected);
        spin_sleep::sleep(time::Duration::from_millis((diff.to_std().unwrap().as_millis() / 2) as u64));
    }

    t.cookies = res.cookies;
    return status(t);
}