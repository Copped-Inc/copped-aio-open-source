use crate::{api, cache, log, session};

#[derive(Clone)]
#[derive(serde::Serialize)]
#[derive(serde::Deserialize)]
pub struct Payment {
    pub id: String,
    pub store: String,
    pub data: String,
    pub state: i8,
}

fn remove(id: String) {
    cache::lock();
    let mut payments = cache::payments();

    payments.retain(|x| x.id != id);

    cache::set_payments(payments);
    cache::unlock();
}

fn add(payment: Payment) {
    cache::lock();
    let mut payments = cache::payments();
    payments.push(payment);
    cache::set_payments(payments);
    cache::unlock();
}

fn get(id: String) -> Option<Payment> {
    cache::lock();
    let payments = cache::payments();
    cache::unlock();

    for payment in payments {
        if payment.id == id {
            return Some(payment);
        }
    }

    return None;
}

impl Payment {
    pub fn update(self) {
        cache::lock();
        let mut payments = cache::payments();

        payments.retain(|x| x.id != self.id);
        payments.push(self);

        cache::set_payments(payments);
        cache::unlock();
    }
}

impl session::Session {
    pub async fn payments(&mut self, store: String) -> Result<String, &'static str> {
        let id = uuid::Uuid::new_v4().to_string();
        let mut data = String::new();
        for (key, value) in self.add_data.iter() {
            data.push_str(key);
            data.push_str(":");
            data.push_str(value);
            data.push_str(";");
        }

        let payment = Payment {
            id: id.clone(),
            store,
            data,
            state: 0,
        };

        log!("Waiting for payments instance to be ready...");
        loop {
            cache::lock();
            let data = cache::data();
            cache::unlock();

            if data.instances.is_none() {
                spin_sleep::sleep(std::time::Duration::from_millis(1000));
                continue;
            }

            let instances = data.instances.unwrap();
            let found = instances.iter().find(|x| x.provider == "Payments" && x.status == "Running");
            if found.is_some() {
                break;
            }

            spin_sleep::sleep(std::time::Duration::from_millis(1000));
        }

        let r = api::payment(payment.clone()).await;
        if r.is_err() {
            return Err(r.unwrap_err());
        }

        add(payment);
        return Ok(id);
    }
}

pub fn check(id: String) -> Result<(), &'static str> {
    loop {
        spin_sleep::sleep(std::time::Duration::from_millis(1000));
        let payment = get(id.clone());
        if payment.is_none() {
            return Err("Payment not found");
        }
        let payment = payment.unwrap();

        if payment.state == 0 {
            remove(id.clone());
            return Err("Timeout");
        } else if payment.state == 2 {
            remove(id.clone());
            return Err("Payment failed");
        } else if payment.state == 3 {
            remove(id.clone());
            return Ok(());
        }
    }
}
