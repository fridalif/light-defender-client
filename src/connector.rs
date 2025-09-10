use futures_util::{SinkExt, StreamExt};
use std::sync::{Arc, Mutex};
use tokio_tungstenite::{connect_async, MaybeTlsStream, WebSocketStream};

use crate::config::AppConfig;

pub struct Connector {
    pub app_config: AppConfig
}

pub trait ConnectorTrait {
    async fn new(app_config: AppConfig) -> Self;
    async fn connect(&self) -> Result<Arc<Mutex<WebSocketStream<MaybeTlsStream<tokio::net::TcpStream>>>>, Box<dyn std::error::Error>>;
    async fn run(&self);
    async fn serve_websocket(&self, ws_stream: Arc<Mutex<WebSocketStream<MaybeTlsStream<tokio::net::TcpStream>>>>);
}

impl ConnectorTrait for Connector {
    async fn new(app_config: AppConfig) -> Self {
        Connector { app_config }
    }

    async fn run (&self) {
        loop {
            match self.connect().await {
                Ok(ws_stream) => {
                    self.serve_websocket(ws_stream).await;
                }
                Err(e) => {
                    println!("Failed to connect: {}", e);
                }
            }
            tokio::time::sleep(std::time::Duration::from_secs(60)).await;
        }
    }

    async fn connect(&self) -> Result<Arc<Mutex<WebSocketStream<MaybeTlsStream<tokio::net::TcpStream>>>>, Box<dyn std::error::Error>> {
        let (ws_stream, _) = connect_async("ws://localhost:8000").await?;
        Ok(Arc::new(Mutex::new(ws_stream)))
    }

    async fn serve_websocket(&self, ws_stream: Arc<Mutex<WebSocketStream<MaybeTlsStream<tokio::net::TcpStream>>>>) {
        let ws_stream = self.connect().await.unwrap();
        let (mut ws_stream, _) = ws_stream.lock().await.split();
    }
}

