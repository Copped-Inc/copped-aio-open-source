use crate::api::login;

pub struct Login {
    pub authorization: String,
    pub code: String,
}

impl Login {
    pub fn new() -> Self {
        Self {
            authorization: String::from(""),
            code: "".to_string()
        }
    }

    pub fn code(self, code: String) -> Self {
        Self {
            authorization: self.authorization,
            code,
        }
    }

    pub async fn check(self) -> Result<Self, Box<dyn std::error::Error + Send + Sync>> {
        let code = self.code;
        let r = login(code.clone()).await;

        if let Ok(l) = r {
            Ok(Self {
                authorization: l.authorization,
                code,
            })
        } else {
            Err(r.unwrap_err())
        }
    }
}