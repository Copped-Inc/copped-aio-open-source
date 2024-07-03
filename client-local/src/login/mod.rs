mod values;

use crate::console::{input, log};
use crate::login::values::{Login};

pub(crate) async fn create() -> Login {
    loop {
        let i = input("Enter your code: ");
        let r = Login::new().code(i).check().await;

        if let Ok(l) = r {
            return l;
        } else {
            log("Failed to login");
        }
    }
}
