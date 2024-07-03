use crate::api::checkout;
use crate::data::values::Product;
use crate::log;

impl Product {
    pub(crate) async fn checkout(self, url: &str) {
        open::that(url).expect("can't open 3DS");
        unsafe {
            let r = checkout(self).await;
            if let Ok(_) = r {
                log("Checkout success");
            } else {
                log(r.unwrap_err().to_string().as_str());
            }
        }
    }
}