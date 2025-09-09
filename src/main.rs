
mod config;
mod cryptography;

use std::{error::Error, sync::{Arc, Mutex}};
use config::Config;

fn main() -> Result<(), Box<dyn Error>>{
    let current_configuration = Config::new("../etc/config.bin")?;
    
}
