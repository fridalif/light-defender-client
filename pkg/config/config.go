package config

import (
	"fmt"
	"light-defender-client/pkg/cryptography"
	"os"
)

type PublicConfig struct {
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
	fmt.Println(string(plaintext))
	return nil, nil
}
