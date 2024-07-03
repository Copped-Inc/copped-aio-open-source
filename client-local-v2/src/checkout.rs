use crate::{api, data, error, log};

impl data::Product {
    pub async fn checkout(self, url: &str) {
        if url != "" {
            open::that(url).expect("can't open 3DS");
        }
        let r = api::checkout(self).await;
        if let Ok(_) = r {
            log!("Checkout success");
        } else {
            error!("{}", r.unwrap_err().to_string().as_str());
        }
    }
}