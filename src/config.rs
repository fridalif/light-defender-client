use uuid::Uuid;


#[derive(Debug, Clone)]
pub struct Config {
    pub id: Uuid,
    pub server_public_key: String,
    pub client_public_key: String,
    pub client_private_key: String,
    pub token: String,
    pub connector_address: String,
}