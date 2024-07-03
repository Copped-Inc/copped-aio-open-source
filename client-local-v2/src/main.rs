extern crate core;

use std::env;
use crate::request::proxy;
use crate::settings::Settings;

mod console;
mod settings;
mod login;
mod api;
mod cache;
mod websocket;
mod data;
mod threads;
mod kitheu;
mod request;
mod session;
mod jig;
mod checkout;
mod aboutyou;
mod payments;

#[tokio::main]
async fn main() {
    log!("{}", "Starting Copped AIO");

    cache::init();
    let args: Vec<String> = env::args().collect();
    if args.len() > 1 && args[1] == "dev"{
        log!("Running in dev mode");
        cache::set_dev(true);
    }

    cache::set_settings(Settings::get().await);
    cache::settings().update().await;
    session::load_sessions();

    threads::init();
    proxy::init().await;

    cache::check_thread();
    websocket::websocket().await;
}