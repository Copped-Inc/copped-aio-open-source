extern crate core;

mod login;
mod console;
mod api;
mod performance;
mod settings;
mod websocket;
mod data;
mod request;
mod modules;
mod jig;
mod captcha;

use crate::console::log;
use crate::login::create;
use crate::request::proxy::ReqProxy;
use crate::websocket::websocket;
use crate::settings::get;

#[tokio::main]
async fn main() {
    let mut s = get().await;
    s.task_max = 5;

    unsafe { ReqProxy::get_from_api(s.authorization.clone()).await; }
    s.clone().update().await;

    unsafe { websocket(s).await; }
}
