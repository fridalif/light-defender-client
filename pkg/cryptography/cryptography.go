package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	cryptorand "crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"errors"
	"fmt"
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
	// 1. Генерируем случайный AES ключ (256 бит для AES-GCM)
	aesKey := make([]byte, 32)
	if _, err := io.ReadFull(cryptorand.Reader, aesKey); err != nil {
		return nil, fmt.Errorf("ошибка генерации AES ключа: %w", err)
	}

	// 2. Шифруем сообщение с помощью AES-GCM
	encryptedData, iv, err := encryptWithAES(aesKey, message)
	if err != nil {
		return nil, fmt.Errorf("ошибка AES шифрования: %w", err)
	}

	// 3. Шифруем AES ключ с помощью RSA-OAEP
	encryptedAESKey, err := rsa.EncryptOAEP(
		sha256.New(),
		cryptorand.Reader,
		publicKey,
		aesKey,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка шифрования RSA ключа: %w", err)
	}

	// 4. Формируем итоговый результат: [RSA зашифрованный ключ] + [IV] + [AES зашифрованные данные]
	result := append(encryptedAESKey, iv...)
	result = append(result, encryptedData...)

	return result, nil
}

// DecryptMessage дешифрует сообщение, зашифрованное гибридным подходом
func DecryptMessage(privateKey *rsa.PrivateKey, cipherData []byte) ([]byte, error) {
	// 1. Извлекаем зашифрованный AES ключ (размер зависит от RSA ключа)
	keySize := privateKey.Size()
	if len(cipherData) < keySize+12 { // минимальный размер: ключ + IV
		return nil, errors.New("недостаточная длина зашифрованных данных")
	}

	encryptedAESKey := cipherData[:keySize]
	iv := cipherData[keySize : keySize+12]
	encryptedData := cipherData[keySize+12:]

	// 2. Дешифруем AES ключ с помощью RSA
	aesKey, err := rsa.DecryptOAEP(
		sha256.New(),
		cryptorand.Reader,
		privateKey,
		encryptedAESKey,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка дешифрования RSA ключа: %w", err)
	}

	// 3. Дешифруем данные с помощью AES-GCM
	message, err := decryptWithAES(aesKey, iv, encryptedData)
	if err != nil {
		return nil, fmt.Errorf("ошибка AES дешифрования: %w", err)
	}

	return message, nil
}

// encryptWithAES шифрует данные с помощью AES-GCM
func encryptWithAES(key []byte, plaintext []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	// Создаем GCM режим
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	// Генерируем nonce (IV)
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(cryptorand.Reader, nonce); err != nil {
		return nil, nil, err
	}

	// Шифруем данные
	ciphertext := aesGCM.Seal(nil, nonce, plaintext, nil)

	return ciphertext, nonce, nil
}

// decryptWithAES дешифрует данные с помощью AES-GCM
func decryptWithAES(key []byte, nonce []byte, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func GenerateKeys(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(cryptorand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, &privateKey.PublicKey, nil
}
