use std::sync::{Arc, Mutex, MutexGuard};
use tokio::runtime;
use tokio::runtime::Handle;
use crate::log;

static mut RUNTIME: Option<runtime::Runtime> = None;
static mut HANDLE: Option<&Handle> = None;

pub fn init() {
    let runtime = runtime::Builder::new_multi_thread().
        enable_all().
        max_blocking_threads(512).
        worker_threads(512).
        build().unwrap();

    unsafe {
        RUNTIME = Some(runtime);
    }

    rt().spawn(async {
        log!("Thread pool initialized");
    });

    unsafe {
        HANDLE = Some(rt().handle());
    }
}

pub fn rt() -> &'static runtime::Runtime {
    unsafe {
        RUNTIME.as_ref().unwrap()
    }
}

pub fn handle() -> &'static Handle {
    unsafe {
        HANDLE.unwrap()
    }
}

pub trait WaitLock<T> {
    fn wait_lock(&self) -> MutexGuard<T>;
}

impl<T> WaitLock<T> for Arc<Mutex<T>> {
    fn wait_lock(&self) -> MutexGuard<T> {
        loop {
            let t = self.try_lock();
            if t.is_err() {
                std::thread::sleep(std::time::Duration::from_millis(10));
                continue;
            }
            return t.unwrap();
        }
    }
}
