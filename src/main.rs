
mod config;
mod cryptography;
mod connector;

use std::{error::Error, sync::{Arc, Mutex}};
use config::{Config, AppConfig};

#[tokio::main]
async fn main() -> Result<(), Box<dyn Error>>{
    let current_configuration = Config::new("../etc/config.bin");
    let (client_private_key, client_public_key) = cryptography::generate_rsa_keys().unwrap();

    let real_config = Arc::new(Mutex::new(AppConfig::new(current_configuration, client_private_key, client_public_key)));
    print!("{:?}", real_config);
    

    Ok(())
}
