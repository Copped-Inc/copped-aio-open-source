use std::thread::sleep;
use async_recursion::async_recursion;
use chrono::{Local, Timelike};
use crate::api::captcha;

pub(crate) static mut CAPTCHA: Vec<Captcha> = Vec::new();
pub(crate) static mut NEED_CAPTCHA: bool = false;
static mut RETRY_COUNTER: i32 = 0;

#[derive(Clone)]
pub(crate) struct Captcha {
    pub(crate) captcha_one: String,
    pub(crate) captcha_two: String,
    pub(crate) time: i64,
}

pub(crate) async fn captcha_loop(auth: String, task_max: i32) {
    unsafe {
        loop {

            let now = Local::now();
            let hour = now.clone().hour();
            let minute = now.clone().minute();

            CAPTCHA.retain(|x| x.time + 120 > now.timestamp());

            if (hour == 9 && minute > 57 || NEED_CAPTCHA) && CAPTCHA.len() <= task_max as usize {
                gen_captcha(auth.clone()).await;
                NEED_CAPTCHA = false;
            }

            sleep(std::time::Duration::from_millis(100));

        }
    }
}

#[async_recursion]
pub(crate) async unsafe fn gen_captcha(auth: String) {

    let res = captcha(auth.clone()).await;
    if let Ok(res) = res {
        CAPTCHA.push(Captcha {
            captcha_one: res.token[0].clone(),
            captcha_two: res.token[1].clone(),
            time: res.expire.timestamp(),
        });
    } else {
        RETRY_COUNTER += 1;
        if RETRY_COUNTER > 3 {
            return;
        }

        sleep(std::time::Duration::from_secs(5));
        gen_captcha(auth.clone()).await;
    }

}

pub(crate) fn get_captcha() -> Result<Captcha, String> {

    unsafe {
        if RETRY_COUNTER > 3 {
            RETRY_COUNTER = 0;
            return Err("captcha error".to_string());
        }
    }

    let c = unsafe { CAPTCHA.pop() };
    if let Some(c) = c {
        return Ok(c);
    }

    unsafe {
        NEED_CAPTCHA = true;
    }

    sleep(std::time::Duration::from_secs(1));
    get_captcha()

}