
mod config;
mod cryptography;

use std::{error::Error, sync::{Arc, Mutex}};
use config::Config;

fn main() -> Result<(), Box<dyn Error>>{
    let current_configuration = Arc::new(Mutex::new(Config::new("../etc/config.bin")));
    if current_configuration.lock().unwrap().client_public_key.is_empty() {
        let (client_private_key, client_public_key) = cryptography::generate_rsa_keys().unwrap();
        current_configuration.lock().unwrap().client_public_key = client_public_key;   
        print!("{:?}", current_configuration);
    }
    print!("{:?}", current_configuration);
    Ok(())
}
