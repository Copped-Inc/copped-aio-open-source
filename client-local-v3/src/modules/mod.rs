use crate::cache::Cache;
use crate::session::{Sessions, TaskType};
use crate::settings;

mod kitheu;
mod queueit;

pub fn build(f: bool) {
    if f {
        settings::delete_sessions();
        Cache::set_sessions(vec![]);
    }

    Sessions::new("kith_eu")
        .task_type(TaskType::Product)
        .fallback(2)
        .delete_values(vec![
            "captcha_cart",
            "captcha_checkout",
            "captcha",
            "location"
        ])
        .add_run(vec![
            kitheu::harvest_captcha,
            kitheu::cart,
            kitheu::harvest_captcha,
            kitheu::get_captcha,
            kitheu::cart_token,
            kitheu::address,
            kitheu::card,
            kitheu::checkout,
        ]).build();

    Sessions::new("queue_it")
        .task_type(TaskType::OneTime)
        .task_multiply(200)
        .checkout_value("redirect_url")
        .add_run(vec![
            queueit::queue_it,
            queueit::event_id,
            queueit::enqueue,
            queueit::status,
        ]).build();
}
