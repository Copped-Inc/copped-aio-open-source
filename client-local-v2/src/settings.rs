use platform_dirs::AppDirs;
use std::fs;
use std::process::exit;
use self_update::cargo_crate_version;
use crate::{api, cache, error, log};
use crate::login::Login;

pub fn delete() {
    let settings = AppDirs::new(Some("Copped AIO"), true).unwrap().
        config_dir.join(format!("settings{}.json", if cache::dev() { "-dev" } else { "" }));
    let _ = fs::remove_file(settings);

    let sessions = AppDirs::new(Some("Copped AIO"), true).unwrap().
        config_dir.join(format!("sessions{}.json", if cache::dev() { "-dev" } else { "" }));
    let _ = fs::remove_file(sessions);

    exit(0);
}

#[derive(serde::Deserialize)]
#[derive(serde::Serialize)]
#[derive(Debug)]
#[derive(Clone)]
#[derive(Default)]
pub struct Settings {
    pub authorization: String,
    pub id: String,
    pub price: f64,
    pub provider: String,
    pub task_max: i32,
    pub region: String,
}

impl Settings {
    pub async fn get() -> Self {
        let path = AppDirs::new(Some("Copped AIO"), true).unwrap().
            config_dir.join(format!("settings{}.json", if cache::dev() { "-dev" } else { "" }));

        if !path.exists() {
            fs::create_dir_all(AppDirs::new(Some("Copped AIO"), true).unwrap().config_dir).unwrap();

            log!("Login needed");
            let auth = Login::create().await;
            log!("Got authorization");

            let s = Settings::new(
                auth.authorization,
                auth.code,
                0.0,
                "Self Hosted".to_string(),
                5,
                "Unavailable".to_string()
            ).save();

            if let Ok(_) = s {
                log!("Settings file created");
            } else {
                error!("Could not create settings file");
            }
        }

        let d = fs::read_to_string(path)
            .expect("Failed to read settings file");

        let mut s: Settings = serde_json::from_str(&d).expect("Failed to parse settings file");
        if s.task_max > 5 {
            println!("Task max is too high ({}), setting to 5", s.task_max);
            s.task_max = 5;
            let _ = s.clone().save();
        }

        s
    }

    fn new(auth: String, id: String, price: f64, provider: String, task_max: i32, region: String) -> Self {
        Self {
            authorization: auth,
            id,
            price,
            provider,
            task_max,
            region,
        }
    }

    fn save(self) -> std::io::Result<Self> {
        let path = AppDirs::new(Some("Copped AIO"), true).unwrap().
            config_dir.join(format!("settings{}.json", if cache::dev() { "-dev" } else { "" }));

        let file = fs::File::create(path)?;
        serde_json::to_writer_pretty(file, &self)?;

        Ok(self)
    }

    pub async fn update(self) {
        loop {
            log!("Current version: {}", cargo_crate_version!());
            if cache::dev() {
                log!("Skipping update");
                return;
            }

            let _ = fs::remove_file("update.bat");
            let r = api::update().await;
            if r.is_ok() {
                log!("Newest version installed");
                return;
            } else {
                error!("Update failed - Retrying. {}", r.unwrap_err());
                spin_sleep::sleep(std::time::Duration::from_secs(1));
            }
        }
    }

    pub fn running(self) -> bool {
        let instances = cache::data().instances.unwrap();
        for i in instances {
            if i.id == self.id {
                return true;
            }
        }
        false
    }
}