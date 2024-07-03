use std::env;
use crate::{console, error, log};
use crate::api;

#[derive(Clone)]
#[derive(Debug)]
pub struct Login {
    pub authorization: String,
    pub code: String,
}

impl Login {
    pub async fn create() -> Self {
        let name = "INSTANCE_ID";
        match env::var(name) {
            Ok(code) => {
                log!("Using passed code {}", code);
                let r = Login::code(code).check().await;
                if let Ok(l) = r {
                    return l;
                } else {
                    panic!("Failed to login");
                }
            }
            Err(_) => {
                loop {
                    let i = console::input("Enter your code:");
                    let r = Login::code(i).check().await;

                    if let Ok(l) = r {
                        return l;
                    } else {
                        error!("Failed to login");
                    }
                }
            }
        }
    }

    fn code(code: String) -> Self {
        Self {
            authorization: "".to_string(),
            code,
        }
    }

    async fn check(self) -> Result<Self, &'static str> {
        let r = api::login(self.clone().code).await;
        if let Ok(l) = r {
            Ok(Self {
                authorization: l.authorization,
                code: self.code,
            })
        } else {
            Err(r.unwrap_err())
        }
    }
}