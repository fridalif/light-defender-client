package connector

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"light-defender-client/pkg/config"
	"light-defender-client/pkg/cryptography"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type ConnectorI interface {
	Run() error
}

type Connector struct {
	AppConfig     *config.Config
	privateKey    *rsa.PrivateKey
	connectClosed bool
}

func NewConnector(appConfig *config.Config) *Connector {
	return &Connector{AppConfig: appConfig, connectClosed: true}
}

func (c *Connector) Run() error {
	for {

		// Генерируем ключи rsa
		privateKey, publicKey, err := cryptography.GenerateKeys(2048)
		if err != nil {
			return err
		}
		c.privateKey = privateKey

		// Устанавливаем соединение
		wsConnection, _, err := websocket.DefaultDialer.Dial(c.AppConfig.PubConfig.ConnectorAddress, nil)
		if err != nil {
			log.Println("dial:", err)
			continue
		}
		c.connectClosed = false
		publicKeyB64 := base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PublicKey(publicKey))
		firstMessage := map[string]interface{}{
			"public_key": publicKeyB64,
			"token":      c.AppConfig.PubConfig.Token,
		}

		err = wsConnection.WriteJSON(firstMessage)
		if err != nil {
			wsConnection.Close()
			log.Println(err)
			continue
		}

		done := make(chan struct{})

		go c.readMessages(done, wsConnection)
		go c.writeMessages(done, wsConnection)

		<-done

		wsConnection.Close()
		time.Sleep(5 * time.Second)
		break
	}
	return nil
}

func (c *Connector) readMessages(done chan struct{}, conn *websocket.Conn) {
	defer func() {
		if !c.connectClosed {
			close(done)
			c.connectClosed = true
		}
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("Received: %s", message)
	}
}

func (c *Connector) writeMessages(done chan struct{}, conn *websocket.Conn) {
	defer func() {
		if !c.connectClosed {
			close(done)
			c.connectClosed = true
		}
	}()

	attempts := 0
	for {
		if attempts > 10 {
			return
		}
		time.Sleep(5 * time.Second)
		attempts++
		message := []byte("hello")
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("write:", err)
			return
		}
	}
}
