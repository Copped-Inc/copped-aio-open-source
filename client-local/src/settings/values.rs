use std::fs;
use std::fs::File;
use std::io::Write;
use std::thread::sleep;
use platform_dirs::AppDirs;
use async_recursion::async_recursion;
use serde::{Serialize, Deserialize};
use crate::api::update;
use crate::data::values::Instance;
use crate::log;

#[derive(Deserialize)]
#[derive(Serialize)]
#[derive(Debug)]
#[derive(Clone)]
pub struct Settings {
    pub authorization: String,
    pub id: String,
    pub price: f64,
    pub provider: String,
    pub task_max: i32,
    pub region: String,
}

impl Settings {
    pub(crate) fn new() -> Self {
        Self {
            authorization: String::from(""),
            id: "".to_string(),
            price: 0.0,
            provider: String::from(""),
            task_max: 0,
            region: String::from(""),
        }
    }

    pub(crate) fn auth(self, auth: String) -> Self {
        Self {
            authorization: auth,
            ..self
        }
    }

    pub(crate) fn id(self, id: String) -> Self {
        Self {
            id,
            ..self
        }
    }

    pub(crate) fn price(self, price: f64) -> Self {
        Self {
            price,
            ..self
        }
    }

    pub(crate) fn provider(self, provider: String) -> Self {
        Self {
            provider,
            ..self
        }
    }

    pub(crate) fn task_max(self, task_max: i32) -> Self {
        Self {
            task_max,
            ..self
        }
    }

    pub(crate) fn region(self, region: String) -> Self {
        Self {
            region,
            ..self
        }
    }

    pub(crate) fn create(self) -> std::io::Result<Self> {
        let app_dirs = AppDirs::new(Some("Copped AIO"), true).unwrap();
        let config_file_path = app_dirs.config_dir.join("settings.json");
        fs::create_dir_all(&app_dirs.config_dir).unwrap();

        let mut file = if config_file_path.exists() {
            File::open(config_file_path).unwrap()
        } else {
            File::create(config_file_path).unwrap()
        };

        file.write_all(serde_json::to_string(&self).unwrap().as_bytes())?;
        Ok(self)
    }

    #[async_recursion]
    pub(crate) async fn update(self) {
        let _ = fs::remove_file("update.bat");
        let s = update(self.clone().authorization).await;
        if let Ok(_) = s {
            log("Newest version installed");
        } else {
            log("Update failed - Retrying.");
            sleep(std::time::Duration::from_secs(1));
            self.update().await;
        }
    }

    pub(crate) fn running(self, instances: &Vec<Instance>) -> bool {
        for i in instances {
            if i.id == self.id {
                if i.status == "Running" {
                    return true;
                }
            }
        }
        false
    }
}