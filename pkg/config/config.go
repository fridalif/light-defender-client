package config

import (
	"encoding/base64"
	"encoding/json"
	"light-defender-client/pkg/cryptography"
	"os"
)

type PublicConfig struct {
	ServerPublicKeyB64 string `json:"server_public_key"`
	Token              string `json:"token"`
	ConnectorAddress   string `json:"connector_address"`
	ServerPublicKey    []byte
}

type PrivateConfig struct {
}

type Config struct {
	PubConfig  *PublicConfig
	PrivConfig *PrivateConfig
}

func NewConfig() (*Config, error) {
	cipheredBytes, err := os.ReadFile("./etc/config.bin")
	if err != nil {
		return nil, err
	}
	plaintext, err := cryptography.DecryptConfig(cipheredBytes, []byte("01234567890123456789012345678901")) //6ba7885277793bca54b3c26ee9a6b72a
	if err != nil {
		return nil, err
	}
	var pubConfig PublicConfig
	err = json.Unmarshal(plaintext, &pubConfig)
	if err != nil {
		return nil, err
	}
	serverPKFromB64, err := base64.StdEncoding.DecodeString(pubConfig.ServerPublicKeyB64)

	if err != nil {
		return nil, err
	}

	pubConfig.ServerPublicKey = []byte(serverPKFromB64)

	appConfig := &Config{&pubConfig, nil}
	return appConfig, nil
}
