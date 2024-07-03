extern crate chrono;

use std::io;
use std::io::Write;
use chrono::Local;

const DATE_FORMAT_STR: &'static str = "%Y-%m-%d %H:%M:%S";

pub(crate) fn input(prompt: &str) -> String {
    print!("{}", prompt);
    io::stdout().flush().unwrap();

    let mut i = String::new();
    io::stdin().read_line(&mut i).unwrap();
    i = i.trim().to_string();

    if i == "" {
        log("Empty input");
        input(prompt);
    }
    i
}

pub(crate) fn log(s: &str) {
    let date = Local::now();
    println!("[{}] {}", date.format(DATE_FORMAT_STR), s);
}