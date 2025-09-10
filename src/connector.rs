use futures_util::{SinkExt, StreamExt};
use std::sync::{Arc, Mutex};
use tokio_tungstenite::{connect_async, MaybeTlsStream, WebSocketStream};

use crate::config::AppConfig;

pub struct Connector {
    pub app_config: Arc<Mutex<AppConfig>>,
    pub ws_stream: Option<Arc<Mutex<WebSocketStream<MaybeTlsStream<tokio::net::TcpStream>>>>>
}

pub trait ConnectorTrait {
    async fn new(app_config: Arc<Mutex<AppConfig>>) -> Self;
    async fn connect(&self) -> Result<Arc<Mutex<WebSocketStream<MaybeTlsStream<tokio::net::TcpStream>>>>, Box<dyn std::error::Error>>;
    async fn run(&mut self);
    async fn serve_websocket(&self);
}

impl ConnectorTrait for Connector {
    async fn new(app_config: Arc<Mutex<AppConfig>>) -> Self {
        Connector { app_config, ws_stream: None }
    }

    async fn run (&mut self) {
        loop {
            match self.connect().await {
                Ok(ws_stream) => {
                    self.ws_stream = Some(ws_stream);
                }
                Err(e) => {
                    println!("Failed to connect: {}", e);
                }
            }
            tokio::time::sleep(std::time::Duration::from_secs(60)).await;
        }
    }

    async fn connect(&self) -> Result<Arc<Mutex<WebSocketStream<MaybeTlsStream<tokio::net::TcpStream>>>>, Box<dyn std::error::Error>> {
        let app_config = self.app_config.lock().map_err(|_| "Failed to lock app config")?;
        let (ws_stream, _) = connect_async(&*app_config.client_config.connector_address).await?;
        Ok(Arc::new(Mutex::new(ws_stream)))
    }

    async fn serve_websocket(&self) {
        let ws_stream = self.connect().await.unwrap();
        let (mut ws_stream, _) = ws_stream.lock().await.split();
    }
}

