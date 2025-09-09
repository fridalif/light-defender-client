use uuid::Uuid;
use serde::{Serialize, Deserialize};
use crate::cryptography;


#[derive(Debug, Clone, Deserialize)]
pub struct Config {
    pub id: Uuid,
    pub server_public_key: String,
    pub client_public_key: String,
    pub client_private_key: String,
    pub token: String,
    pub connector_address: String,
}

impl Config {
    pub fn new(
        path: &str,
    ) -> Self {
        let config_str_encrypted = std::fs::read(path)?;
        let config_str = cryptography::decrypt_config(&config_str_encrypted, "6ba7885277793bca54b3c26ee9a6b72a")?;
        let config: Config = serde_json::from_vec(&config_str).expect("Failed to parse config file");
        config
    }
}
