use async_recursion::async_recursion;
use crate::{api, cache, error, kitheu, log};
use crate::cache::Captcha;

pub async fn start() {
    expired().await;
}

async fn expired() {
    cache::lock();
    let mut captchas = kitheu::captchas();

    captchas.retain(|captcha| {
        !captcha.is_expired()
    });

    kitheu::set_captchas(captchas);
    cache::unlock();
    check().await;
}

#[async_recursion]
async fn check() {
    if kitheu::captcha_required() {
        kitheu::remove_captcha_required();

        let r = api::gen_captcha("kitheu").await;
        cache::lock();
        if let Ok(c) = r {
            kitheu::append_captcha(Captcha {
                tokens: c.token,
                expire: c.expire.timestamp(),
            });

            log!("Captcha received");
        } else {
            error!("Failed to request captcha");
            spin_sleep::sleep(std::time::Duration::from_secs(1));
            kitheu::add_captcha_required();
        }
        cache::unlock();
    }

    if kitheu::captcha_required() {
        check().await;
    }
}

pub fn get() -> Option<Captcha> {
    let retry = 100;
    for _ in 0..retry {
        cache::lock();
        let captcha = kitheu::captcha();
        cache::unlock();

        if captcha.is_some() {
            return captcha;
        }
        spin_sleep::sleep(std::time::Duration::from_secs(1));
    }
    None
}