use uuid::Uuid;
use serde::Deserialize;

#[derive(Debug, Clone, Deserialize)]
pub struct Config {
    pub id: Uuid,
    pub server_public_key: String,
    pub client_public_key: String,
    pub client_private_key: String,
    pub token: String,
    pub connector_address: String,
}

//pub secret: 6ba7885277793bca54b3c26ee9a6b72a
impl Config {
    pub fn new(
        path: &str,
    ) -> Self {
        let config_str = std::fs::read_to_string(path).expect("Failed to read config file");
        
        let config: Config = serde_json::from_str(&config_str).expect("Failed to parse config file");

        
        config
    }
}