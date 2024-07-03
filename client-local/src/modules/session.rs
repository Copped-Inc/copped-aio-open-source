use std::ops::Add;
use chrono::{DateTime, Local};
use crate::modules::Task;
use crate::request::cookie::ReqCookie;

#[derive(Clone)]
pub(crate) struct Session {
    pub(crate) task: Task,
    pub(crate) state: i32,
    pub(crate) cookies: Vec<ReqCookie>,
    pub(crate) checkout: String,
    pub(crate) price: String,
    pub(crate) expires: DateTime<Local>
}

pub(crate) static mut SESSIONS: Vec<Session> = vec![];

impl Session {
    pub(crate) fn new(task: Task, state: i32, cookies: Vec<ReqCookie>, checkout: String, price: String) -> Session {
        return Session {
            task,
            state,
            cookies,
            checkout,
            price,
            expires: Local::now().add(chrono::Duration::minutes(30))
        }
    }

    pub(crate) fn get_from_task(task: Task) -> Session {
        unsafe {
            for i in 0..SESSIONS.len() {
                if i > SESSIONS.len() - 1 {
                    break;
                }

                if SESSIONS[i].expires < Local::now() {

                    SESSIONS.remove(i);
                    continue;

                } else if SESSIONS[i].task.prod_id == task.prod_id {

                    return SESSIONS.remove(i);

                }
            }
        }
        Session::new(task, 0, vec![], String::new(), String::new())
    }

    pub(crate) fn save(&self) {
        unsafe {
            SESSIONS.push(self.clone());
        }
    }
}