mod console;
mod cache;
mod settings;
mod api;
mod threads;
mod request;
mod websocket;
mod data;
mod jig;
mod session;
mod modules;

#[macro_use]
extern crate lazy_static;

use crate::cache::Cache;
use crate::request::proxy;
use crate::settings::Settings;

#[tokio::main]
async fn main() {
    log!("Starting Copped AIO v3");
    Cache::init();
    Cache::set_settings(Settings::get().await);
    Cache::settings().update().await;

    threads::init();
    proxy::init().await;

    modules::build(false);
    websocket::connect().await;
}
