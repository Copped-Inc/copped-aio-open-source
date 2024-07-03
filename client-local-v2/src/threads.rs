use tokio::runtime;
use crate::log;

static mut RUNTIME: Option<runtime::Runtime> = None;

pub fn init() {

    let runtime = runtime::Builder::new_multi_thread().
        enable_all().
        max_blocking_threads(512).
        worker_threads(512).
        build().unwrap();

    unsafe {
        RUNTIME = Some(runtime);
        rt().spawn(async {
            log!("Thread pool initialized");
        });
    }

}

pub fn rt() -> &'static runtime::Runtime {
    unsafe {
        RUNTIME.as_ref().unwrap()
    }
}