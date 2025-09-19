package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	cryptorand "crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
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

func EncryptMessage(publicKey *rsa.PublicKey, message []byte) ([]byte, error) {
	cipherText, err := rsa.EncryptOAEP(
		sha256.New(),
		cryptorand.Reader,
		publicKey,
		message,
		nil,
	)
	if err != nil {
		return nil, EncryptMessageError(err)
	}
	return cipherText, nil
}

func DecryptMessage(privateKey *rsa.PrivateKey, cipher []byte) ([]byte, error) {
	message, err := rsa.DecryptOAEP(
		sha256.New(),
		cryptorand.Reader,
		privateKey,
		cipher,
		nil,
	)
	if err != nil {
		return nil, DecryptMessageError(err)
	}
	return message, nil
}

func GenerateKeys(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(cryptorand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, &privateKey.PublicKey, nil
}
