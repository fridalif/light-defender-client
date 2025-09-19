package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	cryptorand "crypto/rand"
	"io"
)

func EncryptConfig(configData []byte, key []byte) ([]byte, error) {
	for len(key) < 32 {
		key = append(key, '0')
	}

	if len(key) > 32 {
		key = key[:32]
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(cryptorand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, configData, nil)
	return ciphertext, nil
}

func DecryptConfig(cipherMes []byte, key []byte) ([]byte, error) {
	for len(key) < 32 {
		key = append(key, '0')
	}

	if len(key) > 32 {
		key = key[:32]
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := cipherMes[:gcm.NonceSize()]
	cipherMessage := cipherMes[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, cipherMessage, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}
