
mod config;
mod cryptography;

use std::{error::Error, sync::{Arc, Mutex}};
use config::{Config, AppConfig};


fn main() -> Result<(), Box<dyn Error>>{
    let current_configuration = Config::new("../etc/config.bin");
    let (client_private_key, client_public_key) = cryptography::generate_rsa_keys().unwrap();

    std::fs::write("../etc/client_public_key.bin", cryptography::encrypt_config(&client_public_key.as_bytes(), b"01234567890123456789012345678901").unwrap()).expect("Failed to write msg");
    std::fs::write("../etc/client_private_key.bin", cryptography::encrypt_config(&client_private_key.as_bytes(), b"01234567890123456789012345678901").unwrap()).expect("Failed to write msg");
    

    let real_config = Arc::new(Mutex::new(AppConfig::new(current_configuration, client_private_key, client_public_key)));


    print!("{:?}", real_config);
    
    Ok(())
}
