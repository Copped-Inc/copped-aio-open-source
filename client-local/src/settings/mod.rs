pub(crate) mod values;

use std::fs;
use std::fs::File;
use std::process::exit;
use platform_dirs::AppDirs;
use crate::{create, log};
use crate::settings::values::Settings;

pub(crate) async fn get() -> Settings {
    let app_dirs = AppDirs::new(Some("Copped AIO"), true).unwrap();
    let config_file_path = app_dirs.config_dir.join("settings.json");
    fs::create_dir_all(&app_dirs.config_dir).unwrap();

    let _ = if config_file_path.exists() {
        File::open(&config_file_path).unwrap();
        log("Found settings file");
    } else {
        log("Login");
        let auth = create().await;
        log("Got authorization");

        log("Test performance - currently disabled");
        let max_tasks = 0/*test().await*/;
        log("Use of standard max tasks of the alpha version");

        let s = Settings::new();
        let r = s
            .auth(auth.authorization)
            .id(auth.code)
            .price(0.0)
            .provider("Self Hosted".to_string())
            .task_max(max_tasks)
            .region("Unavailable".to_string())
            .create();

        if let Ok(_) = r {
            log("Settings file created");
        } else {
            log("r.unwrap_err().to_string().as_str()");
        }
    };

    let d = fs::read_to_string(config_file_path)
        .expect("Failed to read settings file");

    let s:Settings = serde_json::from_str(&d)
        .expect("Failed to parse settings file");

    s
}

pub(crate) fn delete() {
    let app_dirs = AppDirs::new(Some("Copped AIO"), true).unwrap();
    let config_file_path = app_dirs.config_dir.join("settings.json");
    fs::remove_file(config_file_path).unwrap();
    exit(0);
}