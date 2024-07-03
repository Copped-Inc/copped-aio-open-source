use std::{env, fs};
use std::path::PathBuf;
use std::process::exit;
use platform_dirs::AppDirs;
use async_recursion::async_recursion;
use self_update::cargo_crate_version;
use crate::{api, console, log};
use crate::cache::Cache;

pub fn delete() {
    let settings = path("settings");
    let _ = fs::remove_file(settings);
    delete_sessions();

    exit(0);
}

pub fn delete_sessions() {
    let sessions = path("kith_eu");
    let _ = fs::remove_file(sessions);
}

#[derive(serde::Deserialize, serde::Serialize, Clone, Default)]
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
        let p = path("settings");
        if !p.exists() {
            fs::create_dir_all(AppDirs::new(Some("Copped AIO"), true).unwrap().config_dir).unwrap();
            let auth = Login::login().await;

            Settings {
                authorization: auth.authorization,
                id: auth.code,
                price: 0.0,
                provider: {
                    if Cache::dev() {
                        "Self Hosted (Dev)"
                    } else if cfg!(target_os = "windows") {
                        "Self Hosted"
                    } else {
                        "Cloud Hosted"
                    }.to_string()
                },
                task_max: 5,
                region: "Unavailable".to_string()
            }.save();
        }

        let d = fs::read_to_string(p).expect("Failed to read settings file");
        let mut s: Settings = serde_json::from_str(&d).expect("Failed to parse settings file");
        if s.task_max > 10 && !Cache::dev() {
            log!("Error: Task max is too high ({}), setting to 5", s.task_max);
            s.task_max = 10;
            s.clone().save();
        }

        s
    }

    fn save(self) {
        let p = path("settings");
        let file = fs::File::create(p).expect("Failed to create settings file");
        serde_json::to_writer_pretty(file, &self).expect("Failed to write settings file");
    }

    pub async fn update(self) {
        log!("Current version: {}", cargo_crate_version!());
        if Cache::dev() {
            log!("Skipping update");
            return;
        }

        let paths = fs::read_dir("./").unwrap();
        for path in paths {
            if path.is_err() { continue; }
            let p = path.unwrap();
            if p.path().is_dir() && p.file_name().to_str().unwrap().contains("self_update") {
                fs::remove_dir_all(p.path()).unwrap();
            }
        }

        let r = api::update().await;
        if r.is_err() {
            spin_sleep::sleep(std::time::Duration::from_secs(10));
            return self.save();
        }

        log!("Newest version installed");
    }
}

#[derive(Default, Clone)]
struct Login {
    pub authorization: String,
    pub code: String,
}

impl Login {
    #[async_recursion]
    async fn login() -> Self {
        let (code, passed): (String, bool);
        match env::var("INSTANCE_ID") {
            Ok(c) => {
                log!("Using passed code {}", c);
                code = c; passed = true;
            }
            Err(_) => {
                code = console::input("Enter your code:"); passed = false;
            }
        }

        let l = Login::new_with_code(code).check().await;
        if let Err(e) = l {
            log!("{}", e);
            if passed { panic!("Invalid code passed") }
            Login::login().await;
        }

        l.unwrap()
    }

    fn new_with_code(code: String) -> Self {
        Self {
            code,
            ..Self::default()
        }
    }

    async fn check(self) -> Result<Self, &'static str> {
        let r = api::login(self.clone().code).await;
        if r.is_err() {
            return Err(r.unwrap_err());
        }

        Ok(Self {
            authorization: r.unwrap().authorization,
            ..self
        })
    }
}

pub fn path(file: &str) -> PathBuf {
    AppDirs::new(Some("Copped AIO"), true).unwrap().
        config_dir.join(format!("{}{}.json", file, if Cache::dev() { "-dev" } else { "" }))
}
